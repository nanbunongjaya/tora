package cluster

import (
	"os"

	"tora/config"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type (
	ClientSet struct {
		clientset *kubernetes.Clientset // kubernetes api client set
		namespace string                // namespace of master pod
	}
)

func Init() (*ClientSet, error) {
	// Get incluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// Create Kubernetes API client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Get master pod namespace
	namespace, err := getMasterPodNamespace()
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

func getMasterPodNamespace() (string, error) {
	namespace, err := os.ReadFile(config.NAMESPACE_FILE)
	if err != nil {
		return "", err
	}
	return string(namespace), nil
}
