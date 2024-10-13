package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/fs"
	"os"
	"path/filepath"
)

// S3Client creates an S3 service client from the given session.
func S3Client(sess *session.Session) *s3.S3 {
	return s3.New(sess)
}

// uploadDirToS3 uploads all files in a directory to an S3 bucket.
func uploadDirToS3(s3Client *s3.S3, dirPath, bucketName string) error {
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			fmt.Println("Found another directory in the path", path)
			fmt.Println("taking that as well")
			return nil
		}

		fmt.Printf("Uploading %s\n", path)
		return uploadFileToS3(s3Client, path, bucketName)
	})
	return err
}

// uploadFileToS3 uploads a file to an S3 bucket.
func uploadFileToS3(s3Client *s3.S3, filename, bucketName string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Upload the file
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(filepath.Base(filename)),
		Body:          file,
		ContentLength: aws.Int64(fileInfo.Size()),
	})
	return err
}
