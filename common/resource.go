package common

type ResourceStatus struct {
	Running     string `json:"running"`
	Pending     string `json:"pending"`
	Failed      string `json:"failed"`
	Succeeded   string `json:"succeeded"`
	Unknown     string `json:"unknown"`
	Terminating string `json:"terminating"`
}

type resourceType string

const (
	ResourceKindConfigMap  =              "configmaps.v1."
	ResourceKindDaemonSet                =   "daemonsets.v1.apps"
	ResourceKindDeployment               = "deployments.v1.apps"
	ResourceKindEvent                    = "events.v1.events.k8s.io"
	ResourceKindHorizontalPodAutoScaler  = "horizontalpodautoscalers.v1.autoscaling"
	ResourceKindIngress                  = "ingresses.v1.networking.k8s.io"
	ResourceKindJob                      = "jobs.v1.batch"
	ResourceKindCronJob                  = "cronjobs.v1.batch"
	ResourceKindLimitRange               = "limitrange.v1."
	ResourceKindNamespace                = "namespaces.v1."
	ResourceKindNode                     = "nodes.v1"
	ResourceKindPersistentVolumeClaim    = "persistentvolumeclaims.v1."
	ResourceKindPersistentVolume         = "persistentvolumes.v1."
	ResourceKindCustomResourceDefinition = "customresourcedefinitions.v1.apiextensions.k8s.io"
	ResourceKindPod                      = "pods.v1."
	ResourceKindReplicaSet               ="replicasets.v1.apps"
	ResourceKindReplicationController    ="replicationcontrollers.v1."
	ResourceKindResourceQuota            ="resourcequotas.v1."
	ResourceKindSecret                   ="secrets.v1."
	ResourceKindService                = "services.v1."
	ResourceKindServiceAccount           ="serviceaccounts.v1."
	ResourceKindStatefulSet              ="statefulsets.v1.apps"
	ResourceKindStorageClass             ="storageclasses.v1.storage.k8s.io"
	ResourceKindEndpoint                = "endpoints.v1."
	ResourceKindNetworkPolicy           = "networkpolicies.v1.networking.k8s.io"
	ResourceKindClusterRole             = "clusterroles.v1.rbac.authorization.k8s.io"
	ResourceKindClusterRoleBinding      = "clusterrolebindings.v1.rbac.authorization.k8s.io"
	ResourceKindRole                    = "roles.v1.rbac.authorization.k8s.io"
	ResourceKindRoleBinding             = "rolebindings.v1.rbac.authorization.k8s.io"
)