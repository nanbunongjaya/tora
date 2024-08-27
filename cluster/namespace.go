package cluster

import (
	"context"
	"fmt"
	"log"

	"tora/config"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cs *ClientSet) UpsertNamespace() error {
	ns := cs.newNamespace()

	// Check if the Service exists
	existingNamespace, err := cs.getNamespace()
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to check namespace existence: %w", err)
	}

	// Namespace exists
	if existingNamespace != nil {
		log.Println("Namespace exist")
		return nil
	}

	// Namespace does not exist, create it
	return cs.createNamespace(ns)
}

func (cs *ClientSet) newNamespace() *corev1.Namespace {
	/*
		YAML:
		┌─────────────────────┐
		│ apiVersion: v1      │
		│ kind: Namespace     │
		│ metadata:           │
		│   name: tora-slaves │
		└─────────────────────┘
	*/
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.TORA_SLAVE_NAMESPACE,
		},
	}
}

func (cs *ClientSet) getNamespace() (*corev1.Namespace, error) {
	return cs.clientset.CoreV1().
		Namespaces().
		Get(context.Background(), config.TORA_SLAVE_NAMESPACE, metav1.GetOptions{})
}

func (cs *ClientSet) createNamespace(ns *corev1.Namespace) error {
	_, err := cs.clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	log.Println("Namespace created successfully")

	return nil
}
