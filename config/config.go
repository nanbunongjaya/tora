package config

const (
	// Config of cluster
	NAMESPACE_FILE = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

	// Config of slaves
	TORA_SLAVE_NAMESPACE                      = "tora-slaves"
	TORA_SLAVE_BASE_IMAGE_NAME                = "tora-slave-base-image"
	TORA_CONFIGMAP_NAME                       = "tora-slave-services-plugin"
	TORA_CONTROLLER_CLUSTER_ROLE_NAME         = "tora-controller-cluster-role"
	TORA_CONTROLLER_CLUSTER_ROLE_BINDING_NAME = "tora-controller-cluster-role-binding"
	TORA_SERVICE_ACCOUNT_NAME                 = "tora-service-account"

	// Config of base image
	TORA_SLAVE_BASE_IMAGE     = REPO_NAME + "/" + TORA_SLAVE_BASE_IMAGE_NAME + ":" + TORA_SLAVE_BASE_IMAGE_TAG
	TORA_SLAVE_BASE_IMAGE_TAG = "latest"
	REPO_NAME                 = ""

	// Config of plugins
	SO_FILE = "tora_slave_services_plugin.so"
	GO_FILE = "tora_slave_services/tora_slave_services.go"
)
