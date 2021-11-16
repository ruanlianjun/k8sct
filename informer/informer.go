package informer

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"time"
)

type Cli struct {
	informer dynamicinformer.DynamicSharedInformerFactory
	Controller *QueueController
}

func NewInformerCli(client dynamic.Interface, defaultResyncPeriod time.Duration,maxRetries int, namespace string) *Cli {

	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		client, defaultResyncPeriod, namespace, nil)

	queue:=LimitQueue()
	queueController := NewController(informerFactory, queue, maxRetries)

	return &Cli{informer: informerFactory,Controller: queueController}
}

func (c *Cli) Run()  {
	
}