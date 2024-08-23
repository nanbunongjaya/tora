package cluster

import (
	"context"
	"os"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	clientset *kubernetes.Clientset
)

func Init() {
	// Create Kubernetes API client
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

func CreateConfigMap(sofile string) {
	soFileData, err := os.ReadFile(sofile)
	if err != nil {
		panic(err)
	}

	// Build ConfigMap
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tora-controller",
			Namespace: "tora",
		},
		Data: map[string]string{
			"controller.so": string(soFileData),
		},
	}

	// Create ConfigMap
	_, err = clientset.CoreV1().ConfigMaps("tora").Create(context.Background(), configMap, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}
