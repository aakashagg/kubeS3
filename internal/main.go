package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io/fs"
	"os"
	"path/filepath"
)

// Function to create an S3 bucket.
func createBucket(s3Client *s3.Client, bucketName string) error {
	_, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}

// Function to upload a file to an S3 bucket.
func uploadFileToS3(s3Client *s3.Client, filename, bucketName string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file info to set the size of the content
	_, err = file.Stat()
	if err != nil {
		return err
	}

	// Upload the file
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})
	return err
}

// Function to upload all files in a directory to an S3 bucket.
func uploadDirToS3(s3Client *s3.Client, dirPath, bucketName string) error {

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

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		return
	}

	// Create an Amazon S3 service client
	s3Client := s3.NewFromConfig(cfg)

	// Replace with your bucket name and directory path
	bucketName := "your-bucket-name"
	dirPath := "/var/log"

	// Create the S3 bucket
	err = createBucket(s3Client, bucketName)
	if err != nil {
		// Print the error and continue, the bucket might already exist
		fmt.Println("Error creating bucket:", err)
	}

	// Upload all files from dirPath to the bucket
	err = uploadDirToS3(s3Client, dirPath, bucketName)
	if err != nil {
		fmt.Println("Error uploading files:", err)
		return
	}

	fmt.Println("All files uploaded successfully.")
}
