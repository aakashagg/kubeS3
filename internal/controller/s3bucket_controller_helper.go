package controller

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	storagev1 "kubeS3/api/v1"
)

// createSession creates a new AWS session
func createSession(region string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func ifBucketExistsOnS3(sess *session.Session, bucketName string) (bool, error) {

	svc := s3.New(sess)

	_, err := svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		var aerr awserr.Error
		if errors.As(err, &aerr) && aerr.Code() == s3.ErrCodeNoSuchBucket {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func updateBucket(sess *session.Session, bucketName string) error {

	svc := s3.New(sess)
	_, err := svc.PutBucketLogging(&s3.PutBucketLoggingInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}
	return nil
}

func createBucket(sess *session.Session, bucketName string) error {

	svc := s3.New(sess)
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}
	return nil
}

// deleteBucket deletes the S3 bucket
func deleteBucket(sess *session.Session, bucketName string) error {

	svc := s3.New(sess)

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}
	return nil
}

// removeFinalizer removes the finalizer from the S3 bucket
func (r *S3BucketReconciler) removeFinalizer(ctx context.Context, s3Bucket *storagev1.S3Bucket) error {
	// remove the finalizer
	finalizerName := "s3bucket.finalizers.kubes3.io"
	s3Bucket.Finalizers = removeString(s3Bucket.Finalizers, finalizerName)
	if err := r.Update(ctx, s3Bucket); err != nil {
		return err
	}
	return nil
}

// removeString removes the given string from the slice
func removeString(slice []string, s string) []string {
	var result []string
	for _, str := range slice {
		if str != s {
			result = append(result, str)
		}
	}
	return result
}

// handleBucketDeletion handles the deletion of the S3 bucket and removes the finalizer
func (r *S3BucketReconciler) handleBucketDeletion(ctx context.Context, sess *session.Session, s3Bucket *storagev1.S3Bucket, bucketName string) error {
	if s3Bucket.ObjectMeta.DeletionTimestamp != nil {
		// delete the bucket
		if err := deleteBucket(sess, bucketName); err != nil {
			return err
		}

		// remove the finalizer so that k8s can delete the resource
		if err := r.removeFinalizer(ctx, s3Bucket); err != nil {
			return err
		}
	}
	return nil
}
