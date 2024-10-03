/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	storagev1 "kubeS3/api/v1"
)

const (
	defaultRegion     = "us-east-1"
	finalizerS3Bucket = "s3bucket.finalizers.kubes3.io"
)

// S3BucketReconciler reconciles a S3Bucket object
type S3BucketReconciler struct {
	Session *session.Session
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=storage.awsresources.com,resources=s3buckets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=storage.awsresources.com,resources=s3buckets/status,verbs=get;update;create;patch
// +kubebuilder:rbac:groups=storage.awsresources.com,resources=s3buckets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the S3Bucket object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *S3BucketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling S3Bucket...")

	// Fetch the S3Bucket resource
	s3Bucket := &storagev1.S3Bucket{}
	if err := r.Get(ctx, req.NamespacedName, s3Bucket); err != nil {
		logger.Error(err, "Failed to fetch S3Bucket resource")
		return ctrl.Result{}, err
	}

	sess := r.Session
	bucketName := s3Bucket.Spec.BucketName

	logger.Info("S3Bucket data", "BucketName", bucketName, "Region", s3Bucket.Spec.Region, "State", s3Bucket.Status.State, "Size", s3Bucket.Status.Size, "ARN", s3Bucket.Status.ARN)

	// check if bucket is being deleted if so handle deletion
	if err := r.handleBucketDeletion(ctx, sess, s3Bucket, bucketName); err != nil {
		logger.Error(err, "Failed to handle bucket deletion")
		return ctrl.Result{}, err
	}

	logger.Info("Bucket is not being deleted, checking for other operations...")

	// Check if the S3 bucket already exists
	bucketExists, err := IfBucketExistsOnS3(sess, bucketName)
	if err != nil {
		logger.Error(err, "Failed to check if S3 bucket exists")
		return ctrl.Result{}, err
	}

	if bucketExists {
		// Update the S3 bucket configuration

		logger.Info("Bucket already exists, updating the bucket", "BucketName", bucketName)

		if err := UpdateBucket(sess, bucketName); err != nil {
			logger.Error(err, "Failed to update S3 bucket")
			// update logic is work in progress
			return ctrl.Result{}, err
		}
	} else {
		// Create a new S3 bucket
		logger.Info("Bucket does not exists, creating a new S3 bucket", "BucketName", bucketName)

		if err := CreateBucket(sess, bucketName); err != nil {
			logger.Error(err, "Failed to create S3 bucket")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *S3BucketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&storagev1.S3Bucket{}).
		Complete(r)
}
