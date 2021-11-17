package namespace

import (
	"context"
	"github.com/ruanlianjun/k8sct/common"
	v1 "k8s.io/api/core/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type ListChannel struct {
	List  chan *v1.NamespaceList
	Error chan error
}

func GetNamespaceList(client k8sClient.Interface) (*v1.NamespaceList, error) {
	list, err := client.CoreV1().Namespaces().List(context.TODO(), common.ListEverything)
	return list, err
}

func GetNamespaceListChannel(client k8sClient.Interface, readNums int) ListChannel {
	channel := ListChannel{
		List:  make(chan *v1.NamespaceList, readNums),
		Error: make(chan error, readNums),
	}

	go func() {
		list, err := GetNamespaceList(client)
		for i := 0; i < readNums; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()
	return channel
}
