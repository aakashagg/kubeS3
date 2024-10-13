package controller

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"kubeS3/internal/aws"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	storagev1 "kubeS3/api/v1" // Adjust this import path to match your project
)

const finalizerS3Data = "s3data.finalizers.kubes3.io"

// S3DataReconciler reconciles a S3Data object
type S3DataReconciler struct {
	Session *session.Session
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=storage.awsresources.com,resources=s3datas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=storage.awsresources.com,resources=s3datas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=storage.awsresources.com,resources=s3datas/finalizers,verbs=update

func (r *S3DataReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("s3data", req.NamespacedName)

	logger.Info("Reconciling S3Data...")

	// if the S3Data resource is being deleted, handle deletion
	if err := r.handleS3DataDeletion(ctx, nil, req); err != nil {
		logger.Error(err, "Failed to handle S3Data deletion")
		return ctrl.Result{}, err
	}

	// Fetch the S3Data resource
	var s3Data storagev1.S3Data
	if err := r.Get(ctx, req.NamespacedName, &s3Data); err != nil {
		logger.Error(err, "Failed to fetch S3Data resource")
		return ctrl.Result{}, err
	}

	s3client := aws.S3Client(r.Session)
	s3Size, err := aws.GetBucketSize(s3client, s3Data.Spec.S3BucketName)
	if err != nil {
		logger.Error(err, "Failed to get S3 bucket size")
		return ctrl.Result{}, err
	}
	s3Data.Status.Size = s3Size

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *S3DataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&storagev1.S3Data{}).
		Complete(r)
}

// handleS3DataDeletion handles the deletion of S3Data resources
func (r *S3DataReconciler) handleS3DataDeletion(ctx context.Context, session *session.Session, req ctrl.Request) error {
	logger := log.FromContext(ctx)
	// Fetch the S3Data instance
	var s3Data storagev1.S3Data
	if err := r.Get(ctx, req.NamespacedName, &s3Data); err != nil {
		return client.IgnoreNotFound(err)
	}

	// Check if the S3Data resource is marked for deletion
	if s3Data.ObjectMeta.DeletionTimestamp.IsZero() {
		return nil
	}

	// Perform deletion logic here
	// For example, delete data from the S3 bucket if DeletionPolicy is true
	if s3Data.Spec.DeletionPolicy {

		logger.Info("Deleting all data from S3 bucket", "BucketName", s3Data.Spec.S3BucketName)

		s3Client := aws.S3Client(session)
		// Delete data from the S3 bucket
		err := aws.EmptyBucket(s3Client, s3Data.Spec.S3BucketName)
		if err != nil {
			return err
		}
		logger.Info("Successfully deleted all data from S3 bucket", "BucketName", s3Data.Spec.S3BucketName)

	} else {
		logger.Info("DeletionPolicy is false, skipping deletion of data from S3 bucket", "BucketName", s3Data.Spec.S3BucketName)
	}

	// Remove finalizer to allow deletion of the S3Data resource
	s3Data.ObjectMeta.Finalizers = removeString(s3Data.ObjectMeta.Finalizers, finalizerS3Data)
	if err := r.Update(ctx, &s3Data); err != nil {
		return err
	}

	return nil
}
