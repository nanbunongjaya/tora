package cluster

import (
	"context"
	"fmt"
	"log"
	"tora/config"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cs *ClientSet) CreateRole() error {
	/*
		YAML:
		┌──────────────────────────────────────────────┐
		│ apiVersion: rbac.authorization.k8s.io/v1     │
		│ kind: Role                                   │
		│ metadata:                                    │
		│   name: tora-controller-role                 │
		│   namespace: tora-slave                      │
		│ rules:                                       │
		│ - apiGroups: [""]                            │
		│   resources: ["pods", "services"]            │
		│   verbs: ["create", "list", "get", "watch"]  │
		└──────────────────────────────────────────────┘
	*/

	// Set Role
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.TORA_CONTROLLER_ROLE_NAME,
			Namespace: config.TORA_SLAVE_NAMESPACE,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "services"},
				Verbs:     []string{"create", "list", "get", "watch"},
			},
		},
	}

	// Create Role
	_, err := cs.clientset.RbacV1().
		Roles(config.TORA_SLAVE_NAMESPACE).
		Create(context.Background(), role, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create role: %v", err)
	}

	log.Println("Role created successfully")
	return nil
}

func (cs *ClientSet) CreateRoleBinding(targetNamespace, sourceNamespace string) error {
	/*
		YAML:
		┌──────────────────────────────────────────┐
		│ apiVersion: rbac.authorization.k8s.io/v1 │
		│ kind: RoleBinding                        │
		│ metadata:                                │
		│   name: tora-controller-role-binding     │
		│   namespace: <target-namespace>          │
		│ subjects:                                │
		│ - kind: ServiceAccount                   │
		│   name: tora-service-account             │
		│   namespace: <source-namespace>          │
		│ roleRef:                                 │
		│   kind: Role                             │
		│   name: tora-controller-role             │
		│   apiGroup: rbac.authorization.k8s.io    │
		└──────────────────────────────────────────┘
	*/

	// Set RoleBinding
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.TORA_CONTROLLER_ROLE_BINDING_NAME,
			Namespace: targetNamespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      config.TORA_SERVICE_ACCOUNT_NAME,
				Namespace: sourceNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     config.TORA_CONTROLLER_ROLE_NAME,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	// Create RoleBinding
	_, err := cs.clientset.RbacV1().
		RoleBindings(targetNamespace).
		Create(context.Background(), roleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create rolebinding: %v", err)
	}

	log.Println("RoleBinding created successfully")
	return nil
}
