package common

type ResourceStatus struct {
	Running     string `json:"running"`
	Pending     string `json:"pending"`
	Failed      string `json:"failed"`
	Succeeded   string `json:"succeeded"`
	Unknown     string `json:"unknown"`
	Terminating string `json:"terminating"`
}

type ResourceType string

func (r ResourceType) String() string {
	return string(r)
}

const (
	ResourceKindConfigMap                ResourceType =              "configmaps.v1."
	ResourceKindDaemonSet                ResourceType =   "daemonsets.v1.apps"
	ResourceKindDeployment               ResourceType = "deployments.v1.apps"
	ResourceKindEvent                    ResourceType = "events.v1.events.k8s.io"
	ResourceKindHorizontalPodAutoScaler  ResourceType = "horizontalpodautoscalers.v1.autoscaling"
	ResourceKindIngress                  ResourceType = "ingresses.v1.networking.k8s.io"
	ResourceKindJob                      ResourceType = "jobs.v1.batch"
	ResourceKindCronJob                  ResourceType = "cronjobs.v1.batch"
	ResourceKindLimitRange               ResourceType = "limitrange.v1."
	ResourceKindNamespace                ResourceType = "namespaces.v1."
	ResourceKindNode                     ResourceType = "nodes.v1"
	ResourceKindPersistentVolumeClaim    ResourceType = "persistentvolumeclaims.v1."
	ResourceKindPersistentVolume         ResourceType = "persistentvolumes.v1."
	ResourceKindCustomResourceDefinition ResourceType = "customresourcedefinitions.v1.apiextensions.k8s.io"
	ResourceKindPod                      ResourceType = "pods.v1."
	ResourceKindReplicaSet               ResourceType ="replicasets.v1.apps"
	ResourceKindReplicationController    ResourceType ="replicationcontrollers.v1."
	ResourceKindResourceQuota            ResourceType ="resourcequotas.v1."
	ResourceKindSecret                   ResourceType ="secrets.v1."
	ResourceKindService                  ResourceType = "services.v1."
	ResourceKindServiceAccount           ResourceType ="serviceaccounts.v1."
	ResourceKindStatefulSet              ResourceType ="statefulsets.v1.apps"
	ResourceKindStorageClass             ResourceType ="storageclasses.v1.storage.k8s.io"
	ResourceKindEndpoint                 ResourceType = "endpoints.v1."
	ResourceKindNetworkPolicy            ResourceType = "networkpolicies.v1.networking.k8s.io"
	ResourceKindClusterRole              ResourceType = "clusterroles.v1.rbac.authorization.k8s.io"
	ResourceKindClusterRoleBinding       ResourceType = "clusterrolebindings.v1.rbac.authorization.k8s.io"
	ResourceKindRole                     ResourceType = "roles.v1.rbac.authorization.k8s.io"
	ResourceKindRoleBinding              ResourceType = "rolebindings.v1.rbac.authorization.k8s.io"
)