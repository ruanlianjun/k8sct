package storage_class

import (
	"context"
	"github.com/ruanlianjun/k8sct/common"
	"k8s.io/api/storage/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type StorageClassListChannel struct {
	List  chan *v1.StorageClassList
	Error chan error
}

func GetStorageClassList(client k8sClient.Interface) (*v1.StorageClassList, error) {
	list, err := client.StorageV1().StorageClasses().List(context.TODO(), common.ListEverything)
	return list, err
}

func GetStorageClassListChannel(client k8sClient.Interface, readNums int) StorageClassListChannel {
	channel := StorageClassListChannel{
		List:  make(chan *v1.StorageClassList, readNums),
		Error: make(chan error, readNums),
	}

	go func() {
		list, err := GetStorageClassList(client)
		for i := 0; i < readNums; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}
