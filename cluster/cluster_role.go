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

func (cs *ClientSet) UpsertClusterRole() error {
	clusterRole := cs.newClusterRole()

	// Check if the ClusterRole exists
	existingClusterRole, err := cs.getClusterRole()
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to check cluster role existence: %w", err)
	}

	// ClusterRole exists, update it
	if existingClusterRole != nil {
		return cs.updateClusterRole(clusterRole, existingClusterRole)
	}

	// ClusterRole does not exist, create it
	return cs.createClusterRole(clusterRole)
}

func (cs *ClientSet) newClusterRole() *rbacv1.ClusterRole {
	/*
	   YAML:
	   ┌─────────────────────────────────────────────────────────────────┐
	   │ apiVersion: rbac.authorization.k8s.io/v1                        │
	   │ kind: ClusterRole                                               │
	   │ metadata:                                                       │
	   │   name: tora-controller-cluster-role                            │
	   │ rules:                                                          │
	   │ - apiGroups: [""]                                               │
	   │   resources: ["pods", "deployments", "services"]                │
	   │   verbs: ["create", "list", "get", "watch", "update", "delete"] │
	   └─────────────────────────────────────────────────────────────────┘
	*/
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.TORA_CONTROLLER_CLUSTER_ROLE_NAME,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "deployments", "services"},
				Verbs:     []string{"create", "list", "get", "watch", "update", "delete"},
			},
		},
	}
}

func (cs *ClientSet) getClusterRole() (*rbacv1.ClusterRole, error) {
	return cs.clientset.RbacV1().
		ClusterRoles().
		Get(context.Background(), config.TORA_CONTROLLER_CLUSTER_ROLE_NAME, metav1.GetOptions{})
}

func (cs *ClientSet) createClusterRole(clusterRole *rbacv1.ClusterRole) error {
	_, err := cs.clientset.RbacV1().
		ClusterRoles().
		Create(context.Background(), clusterRole, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create cluster role: %v", err)
	}

	log.Println("ClusterRole created successfully")

	return nil
}

func (cs *ClientSet) updateClusterRole(clusterRole, existingClusterRole *rbacv1.ClusterRole) error {
	// Apply resource version
	clusterRole.ResourceVersion = existingClusterRole.ResourceVersion

	_, err := cs.clientset.RbacV1().
		ClusterRoles().
		Update(context.Background(), clusterRole, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update cluster role: %w", err)
	}

	log.Println("ClusterRole updated successfully")

	return nil
}
