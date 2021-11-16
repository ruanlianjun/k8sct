package replication_controller

import (
	"context"
	"github.com/ruanlianjun/k8s-operate/common"
	v1 "k8s.io/api/core/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type ReplicationControllerChannel struct {
	List  chan *v1.ReplicationControllerList
	Error chan error
}

func GetReplicationController(client k8sClient.Interface, nsQuery common.NamespaceQuery) (*v1.ReplicationControllerList, error) {
	list, err := client.CoreV1().ReplicationControllers(nsQuery.ToRequestParam()).List(context.TODO(), common.ListEverything)
	var filterItem []v1.ReplicationController
	for _, item := range list.Items {
		if nsQuery.Matches(item.ObjectMeta.Namespace) {
			filterItem = append(filterItem, item)
		}
	}
	list.Items = filterItem
	return list, err
}

func GetReplicationControllerChannel(client k8sClient.Interface, nsQuery common.NamespaceQuery, readNum int) ReplicationControllerChannel {
	channel := ReplicationControllerChannel{
		List:  make(chan *v1.ReplicationControllerList, readNum),
		Error: make(chan error, readNum),
	}

	go func() {
		list, err := GetReplicationController(client, nsQuery)
		for i := 0; i < readNum; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}
