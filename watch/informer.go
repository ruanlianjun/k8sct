package watch

import (
	"github.com/ruanlianjun/k8sct/common"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"time"
)

type Cli struct {
	informer dynamicinformer.DynamicSharedInformerFactory
	controller *QueueController
}

func NewInformerCli(client dynamic.Interface, defaultResyncPeriod time.Duration,maxRetries int, namespace string) *Cli {

	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		client, defaultResyncPeriod, namespace, nil)
	
	queue:=LimitQueue()
	queueController := NewController(informerFactory, queue, maxRetries)

	return &Cli{informer: informerFactory,controller: queueController}
}

func (c *Cli) AddEventHandler(resourceType common.ResourceType, addFunc, updateFunc, DeleteFunc *HandleFunc) *Cli {
	c.controller = c.controller.AddEventHandler(resourceType,addFunc,updateFunc,DeleteFunc)
	return c
}

func (c *Cli) Run(workerNum int, stopCh chan struct{}) error {
	return c.controller.Run(workerNum,stopCh)
}