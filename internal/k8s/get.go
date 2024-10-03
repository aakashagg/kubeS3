package k8s

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func CreateK8sClientSet() (*kubernetes.Clientset, error) {
	var (
		config *rest.Config
		err    error
	)

	// Attempt to use in-cluster configuration
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to default kubeconfig
		var kubeconfig *api.Config

		kubeconfig, err = clientcmd.NewDefaultClientConfigLoadingRules().Load()
		if err != nil {
			return nil, err
		}

		config, err = clientcmd.NewDefaultClientConfig(*kubeconfig, &clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func FetchAndStoreSecret(clientset *kubernetes.Clientset, namespace, secretName string) (map[string]string, error) {

	// Get the secret from the specified namespace
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %v", err)
	}

	// Convert secret data to a string map
	secretData := make(map[string]string)
	for k, v := range secret.Data {
		secretData[k] = string(v)
	}

	return secretData, nil
}
