package cluster

import (
	"context"
	"fmt"
	"log"
	"os"

	"tora/config"

	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cs *ClientSet) UpsertConfigMap(sofile string) error {
	configMap, err := cs.newConfigMap(sofile)
	if err != nil {
		return err
	}

	// Check if the ConfigMap exists
	existingConfigMap, err := cs.getConfigMap()
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to check deployment existence: %v", err)
	}

	// ConfigMap exists, update it
	if existingConfigMap != nil {
		return cs.updateConfigMap(configMap, existingConfigMap)
	}

	// ConfigMap does not exist, create it
	return cs.createConfigMap(configMap)
}

func (cs *ClientSet) newConfigMap(sofile string) (*v1.ConfigMap, error) {
	content, err := os.ReadFile(sofile)
	if err != nil {
		return nil, err
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
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.TORA_CONFIGMAP_NAME,
			Namespace: config.TORA_SLAVE_NAMESPACE,
		},
		Data: map[string]string{
			config.SO_FILE: string(content),
		},
	}, nil
}

func (cs *ClientSet) getConfigMap() (*v1.ConfigMap, error) {
	return cs.clientset.CoreV1().
		ConfigMaps(config.TORA_SLAVE_NAMESPACE).
		Get(context.Background(), config.TORA_CONFIGMAP_NAME, metav1.GetOptions{})

}

func (cs *ClientSet) createConfigMap(configMap *v1.ConfigMap) error {
	_, err := cs.clientset.CoreV1().
		ConfigMaps(config.TORA_SLAVE_NAMESPACE).
		Create(context.Background(), configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create config map: %v", err)
	}

	log.Println("ConfigMap created successfully")

	return nil
}

func (cs *ClientSet) updateConfigMap(configMap, existingConfigMap *v1.ConfigMap) error {
	// Apply resource version
	configMap.ResourceVersion = existingConfigMap.ResourceVersion

	_, err := cs.clientset.CoreV1().
		ConfigMaps(config.TORA_SLAVE_NAMESPACE).
		Update(context.Background(), configMap, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update config map: %v", err)
	}

	log.Println("ConfigMap updated successfully")

	return nil
}
