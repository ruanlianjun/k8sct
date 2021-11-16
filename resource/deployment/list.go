package deployment

import (
	"context"
	"github.com/ruanlianjun/k8s-operate/common"
	"k8s.io/api/apps/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type ListChannel struct {
	List  chan *v1.DeploymentList
	Error chan error
}

func GetDeploymentList(client k8sClient.Interface, nsQuery common.NamespaceQuery) (*v1.DeploymentList, error) {
	list, err := client.AppsV1().Deployments(nsQuery.ToRequestParam()).List(context.TODO(), common.ListEverything)

	var filterItem []v1.Deployment
	for _, item := range list.Items {
		if nsQuery.Matches(item.ObjectMeta.Namespace) {
			filterItem = append(filterItem, item)
		}
	}

	list.Items = filterItem

	return list, err
}

func GetDeploymentListChannel(client k8sClient.Interface, nsQuery common.NamespaceQuery, readNums int) ListChannel {
	channel := ListChannel{
		List:  make(chan *v1.DeploymentList, readNums),
		Error: make(chan error, readNums),
	}

	go func() {
		list, err := GetDeploymentList(client, nsQuery)
		for i := 0; i < readNums; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}
