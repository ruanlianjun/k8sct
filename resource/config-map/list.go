package config_map

import (
	"context"
	"github.com/ruanlianjun/k8sct/common"
	v1 "k8s.io/api/core/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type ConfigMapListChannel struct {
	List  chan *v1.ConfigMapList
	Error chan error
}

func GetConfigMapList(client k8sClient.Interface, nsQuery common.NamespaceQuery) (*v1.ConfigMapList, error) {
	list, err := client.CoreV1().ConfigMaps(nsQuery.ToRequestParam()).List(context.TODO(), common.ListEverything)

	var filterItems []v1.ConfigMap

	for _, item := range list.Items {
		if nsQuery.Matches(item.ObjectMeta.Namespace) {
			filterItems = append(filterItems, item)
		}
	}

	list.Items = filterItems

	return list, err
}

func GetConfigMapListChannel(client k8sClient.Interface, nsQuery common.NamespaceQuery, readNums int) ConfigMapListChannel {
	channel := ConfigMapListChannel{
		List:  make(chan *v1.ConfigMapList, readNums),
		Error: make(chan error, readNums),
	}

	go func() {
		list, err := GetConfigMapList(client, nsQuery)

		for i := 0; i < readNums; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}
