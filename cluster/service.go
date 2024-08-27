package cluster

import (
	"context"
	"fmt"
	"log"

	"tora/config"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (cs *ClientSet) UpsertService(serviceName, deploymentName string) error {
	service := newService(serviceName, deploymentName)

	// Check if the Service exists
	existingService, err := cs.getService(serviceName)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to check service existence: %w", err)
	}

	// Service exists, update it
	if existingService != nil {
		return cs.updateService(service, existingService)
	}

	// Service does not exist, create it
	return cs.createService(service)
}

func newService(serviceName, deploymentName string) *corev1.Service {
	/*
		YAML:
		┌─────────────────────────────────────────────┐
		│ apiVersion: v1                              │
		│ kind: Service                               │
		│ metadata:                                   │
		│   name: <service-name>                      │
		│   namespace: tora-slaves                    │
		│ spec:                                       │
		│   selector:                                 │
		│     app: <deployment-name>                  │
		│   ports:                                    │
		│   - protocol: TCP                           │
		│     port: 50051                             │
		│     targetPort: 50051                       │
		└─────────────────────────────────────────────┘
	*/

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: config.TORA_SLAVE_NAMESPACE,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": deploymentName,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       50051,
					TargetPort: intstr.FromInt(50051),
				},
			},
		},
	}
}

func (cs *ClientSet) getService(serviceName string) (*corev1.Service, error) {
	return cs.clientset.CoreV1().
		Services(config.TORA_SLAVE_NAMESPACE).
		Get(context.Background(), serviceName, metav1.GetOptions{})
}

func (cs *ClientSet) createService(service *corev1.Service) error {
	_, err := cs.clientset.CoreV1().
		Services(service.Namespace).
		Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	log.Println("Service created successfully")

	return nil
}

func (cs *ClientSet) updateService(service, existingService *corev1.Service) error {
	// Apply resource version
	service.ResourceVersion = existingService.ResourceVersion

	_, err := cs.clientset.CoreV1().
		Services(service.Namespace).
		Update(context.Background(), service, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	log.Println("Service updated successfully")

	return nil
}
