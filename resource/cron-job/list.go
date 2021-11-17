package cron_job

import (
	"context"
	"github.com/ruanlianjun/k8sct/common"
	"k8s.io/api/batch/v1beta1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type CronJobListChannel struct {
	List  chan *v1beta1.CronJobList
	Error chan error
}

func GetCronJobList(client k8sClient.Interface, nsQuery common.NamespaceQuery) (*v1beta1.CronJobList, error) {
	list, err := client.BatchV1beta1().CronJobs(nsQuery.ToRequestParam()).List(context.TODO(), common.ListEverything)

	var filterItem []v1beta1.CronJob
	for _, item := range list.Items {
		if nsQuery.Matches(item.ObjectMeta.Namespace) {
			filterItem = append(filterItem, item)
		}

	}

	list.Items = filterItem

	return list, err
}

func GetCronJobListChannel(client k8sClient.Interface, nsQuery common.NamespaceQuery, readNums int) CronJobListChannel {
	channel := CronJobListChannel{
		List:  make(chan *v1beta1.CronJobList, readNums),
		Error: make(chan error, readNums),
	}

	go func() {
		list, err := GetCronJobList(client, nsQuery)
		for i := 0; i < readNums; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}
