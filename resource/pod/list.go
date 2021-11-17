package pod

import (
	"context"
	"github.com/ruanlianjun/k8sct/common"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type ListChannel struct {
	List  chan *v1.PodList
	Error chan error
}

//List 根据namespace获取所有k8s的pod
func List(client k8sClient.Interface, nsQuery common.NamespaceQuery, options metaV1.ListOptions) (*v1.PodList, error) {
	list, err := client.CoreV1().Pods(nsQuery.ToRequestParam()).List(context.TODO(), options)
	var filteredItems []v1.Pod
	for _, item := range list.Items {
		if nsQuery.Matches(item.ObjectMeta.Namespace) {
			filteredItems = append(filteredItems, item)
		}
	}
	list.Items = filteredItems
	return list, err
}

func GetPodListChannelWithOptions(client k8sClient.Interface, nsQuery common.NamespaceQuery, options metaV1.ListOptions, readNum int) ListChannel {
	channel := ListChannel{
		List:  make(chan *v1.PodList, readNum),
		Error: make(chan error, readNum),
	}
	go func() {
		list, err := List(client, nsQuery, options)

		for i := 0; i < readNum; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}
