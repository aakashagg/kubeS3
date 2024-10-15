package controller

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	v1 "kubeS3/api/v1"
	"kubeS3/internal/aws"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	Session *session.Session
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=pods/log,verbs=get

func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// TODO maybe we need hook for pod resource in k8s, need to check

	logger := log.FromContext(ctx)

	logger.Info("Reconciling Pod Controller")

	var pod corev1.Pod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Printf("Reconciling Pod: %s/%s\n", pod.Name, pod.Namespace)

	s3Bucket, err := r.getS3data(pod) //getting the S3Data name form Annotation
	if err != nil {

		logger.Info("S3DataName not found in the pod annotation", "PodName", pod.Name)
		logger.Info("maybe this pod is not supposed to upload data to S3")

		return ctrl.Result{}, err
	}

	logger.Info("S3DataName for the pod:", "PodName", pod.Name, "S3DataName", s3Bucket.Name)

	pathInDataResource := s3Bucket.Spec.PathOfPod

	s3Client := aws.S3Client(r.Session) //creating a new S3 client from the session

	// main magic happens here, uploading data to S3 bucket
	err = aws.UploadDirToS3(s3Client, pathInDataResource, s3Bucket.Spec.S3BucketName)
	if err != nil {
		logger.Error(err, "Error uploading data to S3, requeing the request after 500 ms")

		if err.Error() == "directory "+pathInDataResource+" does not exist" {
			logger.Info("Directory does not exist, no need to requeue")
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true, RequeueAfter: 500}, err
	}

	logger.Info("Data uploaded to S3 successfully")

	return ctrl.Result{}, nil
}

func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}

func (r *PodReconciler) getS3data(pod corev1.Pod) (v1.S3Data, error) {

	S3DataName := pod.Annotations["S3DataName"]

	var S3Data v1.S3Data
	err := r.Get(context.TODO(), client.ObjectKey{
		Namespace: pod.Namespace,
		Name:      S3DataName,
	}, &S3Data)
	if err != nil {
		return v1.S3Data{}, err
	}
	return S3Data, nil
}
