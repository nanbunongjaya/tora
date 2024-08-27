package cluster

import (
	"context"
	"os"

	"tora/config"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cs *ClientSet) CreateConfigMap(sofile string) error {
	content, err := os.ReadFile(sofile)
	if err != nil {
		return err
	}

	/*
		YAML:
		┌─────────────────────────────────────┐
		│ apiVersion: v1                      │
		│ kind: ConfigMap                     │
		│ metadata:                           │
		│   name: tora-slave-services-plugin  │
		│   namespace: tora_slaves            │
		│ data:                               │
		│   tora_slave_services_plugin.so: |  │
		│     <file-contents>                 │
		└─────────────────────────────────────┘
	*/

	// Set ConfigMap
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.TORA_CONFIGMAP_NAME,
			Namespace: config.TORA_SLAVE_NAMESPACE,
		},
		Data: map[string]string{
			config.SO_FILE: string(content),
		},
	}

	// Create ConfigMap
	_, err = cs.clientset.CoreV1().
		ConfigMaps(config.TORA_SLAVE_NAMESPACE).
		Create(context.Background(), configMap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
