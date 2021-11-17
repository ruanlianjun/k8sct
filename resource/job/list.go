package job

import (
	"context"
	"github.com/ruanlianjun/k8sct/common"
	v1 "k8s.io/api/batch/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type JobListChannel struct {
	List  chan *v1.JobList
	Error chan error
}

func GetJobList(client k8sClient.Interface, nsQuery common.NamespaceQuery) (*v1.JobList, error) {
	list, err := client.BatchV1().Jobs(nsQuery.ToRequestParam()).List(context.TODO(), common.ListEverything)

	var filterItems []v1.Job

	for _, item := range list.Items {
		if nsQuery.Matches(item.ObjectMeta.Namespace) {
			filterItems = append(filterItems, item)
		}
	}

	list.Items = filterItems
	return list, err
}

func GetJobListChannel(client k8sClient.Interface, nsQuery common.NamespaceQuery, readNums int) JobListChannel {
	channel := JobListChannel{
		List:  make(chan *v1.JobList, readNums),
		Error: make(chan error, readNums),
	}

	go func() {
		list, err := GetJobList(client, nsQuery)
		for i := 0; i < readNums; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()
	return channel
}
