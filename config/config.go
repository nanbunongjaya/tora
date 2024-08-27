package config

const (
	// Config of cluster
	NAMESPACE_FILE = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

	// Config of slaves
	TORA_SLAVE_NAMESPACE              = "tora_slaves"
	TORA_CONFIGMAP_NAME               = "tora-slave-services-plugin"
	TORA_CONTROLLER_ROLE_NAME         = "tora-controller-role"
	TORA_CONTROLLER_ROLE_BINDING_NAME = "tora-controller-role-binding"
	TORA_SERVICE_ACCOUNT_NAME         = "tora-service-account"

	// Config of plugins
	SO_FILE = "tora_slave_services_plugin.so"
	GO_FILE = "tora_slave_services/tora_slave_services.go"
)
