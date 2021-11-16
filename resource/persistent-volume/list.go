package persistent_volume

import (
	"context"
	"github.com/ruanlianjun/k8s-operate/common"
	v1 "k8s.io/api/core/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type PersistentVolumeListChannel struct {
	List  chan *v1.PersistentVolumeList
	Error chan error
}

func GetPersistentVolumeList(client k8sClient.Interface) (*v1.PersistentVolumeList, error) {
	list, err := client.CoreV1().PersistentVolumes().List(context.TODO(), common.ListEverything)
	return list, err
}

func GetPersistentVolumeListChannel(client k8sClient.Interface, readNums int) PersistentVolumeListChannel {
	channel := PersistentVolumeListChannel{
		List:  make(chan *v1.PersistentVolumeList, readNums),
		Error: make(chan error, readNums),
	}

	go func() {
		list, err := GetPersistentVolumeList(client)
		for i := 0; i < readNums; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()
	return channel
}
