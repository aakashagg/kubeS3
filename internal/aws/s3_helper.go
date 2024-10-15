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

// UploadDirToS3 uploads all files in a directory to an S3 bucket.

func UploadDirToS3(s3Client *s3.S3, dirPath, bucketName string) error {
	// Check if the directory exists
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", dirPath)
	}
	if err != nil {
		return err
	}

	// Check if the directory is empty
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("directory %s is empty", dirPath)
	}

	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error accessing path", path)
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

// create a function that empties the s3 bucket
func EmptyBucket(s3Client *s3.S3, bucketName string) error {
	// List all objects in the bucket
	resp, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	// Delete all objects in the bucket
	for _, obj := range resp.Contents {
		_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    obj.Key,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// create a function to get the size of the bucket
func GetBucketSize(s3Client *s3.S3, bucketName string) (int64, error) {
	resp, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return 0, err
	}
	var size int64
	for _, obj := range resp.Contents {
		size += *obj.Size
	}
	return size, nil
}
