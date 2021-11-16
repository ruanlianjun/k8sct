package persistent_volume_claim

import (
	"context"
	"github.com/ruanlianjun/k8s-operate/common"
	v1 "k8s.io/api/core/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type PersistentVolumeChainListChannel struct {
	List  chan *v1.PersistentVolumeClaimList
	Error chan error
}

func GetPersistentVolumeChainList(client k8sClient.Interface, nsQuery common.NamespaceQuery) (*v1.PersistentVolumeClaimList, error) {
	list, err := client.CoreV1().PersistentVolumeClaims(nsQuery.ToRequestParam()).List(context.TODO(), common.ListEverything)

	var filterItem []v1.PersistentVolumeClaim
	for _, item := range list.Items {
		if nsQuery.Matches(item.ObjectMeta.Namespace) {
			filterItem = append(filterItem, item)
		}
	}

	list.Items = filterItem
	return list, err
}

func GetPersistentVolumeChainListChannel(client k8sClient.Interface, nsQuery common.NamespaceQuery, readNums int) PersistentVolumeChainListChannel {
	channel := PersistentVolumeChainListChannel{
		List:  make(chan *v1.PersistentVolumeClaimList, readNums),
		Error: make(chan error, readNums),
	}

	go func() {
		list, err := GetPersistentVolumeChainList(client, nsQuery)
		for i := 0; i < readNums; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}
