package node

import (
	"context"
	"github.com/ruanlianjun/k8sct/common"
	v1 "k8s.io/api/core/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type ListChannel struct {
	List  chan *v1.NodeList
	Error chan error
}

func GetNodeList(client k8sClient.Interface) (*v1.NodeList, error) {
	list, err := client.CoreV1().Nodes().List(context.TODO(), common.ListEverything)
	return list, err
}

func GetNodeListChannel(client k8sClient.Interface, readNums int) ListChannel {
	channel := ListChannel{
		List:  make(chan *v1.NodeList, readNums),
		Error: make(chan error, readNums),
	}
	go func() {
		list, err := GetNodeList(client)
		for i := 0; i < readNums; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}
