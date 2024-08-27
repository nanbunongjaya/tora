package cluster

import (
	"context"
	"fmt"
	"log"
	"os"

	"tora/config"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNamespace() (string, error) {
	namespace, err := os.ReadFile(config.NAMESPACE_FILE)
	if err != nil {
		return "", err
	}
	return string(namespace), nil
}

func (cs *ClientSet) CreateNamespace(namespace string) error {
	/*
		YAML:
		┌─────────────────────┐
		│ apiVersion: v1      │
		│ kind: Namespace     │
		│ metadata:           │
		│   name: <namespace> │
		└─────────────────────┘
	*/

	// Set Namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	// Create Namespace
	_, err := cs.clientset.CoreV1().
		Namespaces().
		Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace: %v", err)
	}

	log.Println("Namespace created successfully")
	return nil
}
