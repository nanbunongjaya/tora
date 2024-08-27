package cluster

import (
	"context"
	"fmt"
	"log"

	"github.com/nanbunongjaya/tora/config"

	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cs *ClientSet) UpsertClusterRoleBinding() error {
	clusterRoleBinding := cs.newClusterRoleBinding()

	// Check if the ClusterRoleBinding exists
	existingClusterRoleBinding, err := cs.getClusterRoleBinding()
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to check cluster role binding existence: %w", err)
	}

	// ClusterRoleBinding exists, update it
	if existingClusterRoleBinding != nil {
		return cs.updateClusterRoleBinding(clusterRoleBinding, existingClusterRoleBinding)
	}

	// ClusterRoleBinding does not exist, create it
	return cs.createClusterRoleBinding(clusterRoleBinding)
}

func (cs *ClientSet) newClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	/*
		YAML:
		┌──────────────────────────────────────────┐
		│ apiVersion: rbac.authorization.k8s.io/v1 │
		│ kind: ClusterRoleBinding                 │
		│ metadata:                                │
		│   name: tora-controller-binding          │
		│ subjects:                                │
		│ - kind: ServiceAccount                   │
		│   name: tora-service-account             │
		│   namespace: <source-namespace>          │
		│ roleRef:                                 │
		│   kind: ClusterRole                      │
		│   name: tora-controller-cluster-role     │
		│   apiGroup: rbac.authorization.k8s.io    │
		└──────────────────────────────────────────┘
	*/
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.TORA_CONTROLLER_CLUSTER_ROLE_BINDING_NAME,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      config.TORA_SERVICE_ACCOUNT_NAME,
				Namespace: cs.namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     config.TORA_CONTROLLER_CLUSTER_ROLE_NAME,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

func (cs *ClientSet) getClusterRoleBinding() (*rbacv1.ClusterRoleBinding, error) {
	return cs.clientset.RbacV1().
		ClusterRoleBindings().
		Get(context.Background(), config.TORA_CONTROLLER_CLUSTER_ROLE_BINDING_NAME, metav1.GetOptions{})
}

func (cs *ClientSet) createClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	_, err := cs.clientset.RbacV1().
		ClusterRoleBindings().
		Create(context.Background(), clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create cluster role binding: %w", err)
	}

	log.Println("ClusterRoleBinding created successfully")

	return nil
}

func (cs *ClientSet) updateClusterRoleBinding(clusterRoleBinding, existingClusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	// Apply resource version
	clusterRoleBinding.ResourceVersion = existingClusterRoleBinding.ResourceVersion

	_, err := cs.clientset.RbacV1().
		ClusterRoleBindings().
		Update(context.Background(), clusterRoleBinding, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update cluster role binding: %w", err)
	}

	log.Println("ClusterRoleBinding updated successfully")

	return nil
}
