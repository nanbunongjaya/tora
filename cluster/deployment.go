package cluster

import (
	"context"
	"fmt"
	"log"

	"github.com/nanbunongjaya/tora/config"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cs *ClientSet) UpsertDeployment(appName string) error {
	deployment := cs.newDeployment(appName)

	// Check if the Deployment exists
	existingDeployment, err := cs.getDeployment(appName, config.TORA_SLAVE_NAMESPACE)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to check deployment existence: %v", err)
	}

	// Deployment exists, update it
	if existingDeployment != nil {
		return cs.updateDeployment(deployment, existingDeployment)
	}

	// Deployment does not exist, create it
	return cs.createDeployment(deployment)
}

func (cs *ClientSet) newDeployment(appName string) *appsv1.Deployment {
	/*
		YAML:
		┌───────────────────────────────────────────────────┐
		│ apiVersion: apps/v1                               │
		│ kind: Deployment                                  │
		│ metadata:                                         │
		│   name: <app-name>-deployment                     │
		│   namespace: tora-slaves                          │
		│ spec:                                             │
		│   replicas: 3                                     │
		│   selector:                                       │
		│     matchLabels:                                  │
		│       app: <app-name>-app                         │
		│   template:                                       │
		│     metadata:                                     │
		│       labels:                                     │
		│         app: <app-name>-app                       │
		│     spec:                                         |
		|       serviceAccountName: tora-service-account    │
		│       containers:                                 │
		│       - name: <app-name>-container                │
		│         image: <repo>/tora-slave-base-image:<tag> │
		│         ports:                                    │
		│         - containerPort: 50051                    │
		└───────────────────────────────────────────────────┘
	*/
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      withDeploymentSuffix(appName),
			Namespace: config.TORA_SLAVE_NAMESPACE,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": withAppSuffix(appName),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": withAppSuffix(appName),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  withContainerSuffix(appName),
							Image: config.TORA_SLAVE_BASE_IMAGE,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 50051,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (cs *ClientSet) getDeployment(deploymentName, namespace string) (*appsv1.Deployment, error) {
	return cs.clientset.AppsV1().
		Deployments(namespace).
		Get(context.Background(), deploymentName, metav1.GetOptions{})
}

func (cs *ClientSet) createDeployment(deployment *appsv1.Deployment) error {
	_, err := cs.clientset.AppsV1().
		Deployments(deployment.Namespace).
		Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
	}

	log.Println("Deployment created successfully")

	return nil
}

func (cs *ClientSet) updateDeployment(deployment, existingDeployment *appsv1.Deployment) error {
	// Apply resource version
	deployment.ResourceVersion = existingDeployment.ResourceVersion

	_, err := cs.clientset.AppsV1().
		Deployments(deployment.Namespace).
		Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %v", err)
	}

	log.Println("Deployment updated successfully")

	return nil
}

func withDeploymentSuffix(s string) string {
	return s + "-deployment"
}

func withAppSuffix(s string) string {
	return s + "-app"
}

func withContainerSuffix(s string) string {
	return s + "-container"
}

func int32Ptr(i int32) *int32 {
	return &i
}
