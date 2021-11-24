package watch

import (
	"errors"
	"fmt"
	"github.com/ruanlianjun/k8sct/common"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"sync"
	"time"
)

type QueueController struct {
	factory    dynamicinformer.DynamicSharedInformerFactory
	queue      workqueue.RateLimitingInterface
	wg         *sync.WaitGroup
	liters     map[string]cache.GenericLister
	informers  map[string]cache.SharedIndexInformer
	maxRetries int //最大尝试次数
}

type HandleFunc struct {
	rawData interface{} //queue待处理的数据
	handle  func(rawData interface{}) error
}

func WarpHandleFunc(rawFunc func(rawData interface{}) error) *HandleFunc {
	return &HandleFunc{
		handle: rawFunc,
	}
}

func NewController(factory dynamicinformer.DynamicSharedInformerFactory, queue workqueue.RateLimitingInterface, maxRetries int) *QueueController {

	return &QueueController{
		factory:    factory,
		queue:      queue,
		wg:         &sync.WaitGroup{},
		maxRetries: maxRetries,
		liters:     make(map[string]cache.GenericLister, 0),
		informers:  make(map[string]cache.SharedIndexInformer, 0),
	}

}

// AddEventHandler 资源，比如 "configmaps.v1.", "deployments.v1.apps", "rabbits.v1.stable.wbsnail.com"
func (c *QueueController) AddEventHandler(resourceType common.ResourceType, addFunc, updateFunc, DeleteFunc *HandleFunc) *QueueController {

	resource, groupResource := schema.ParseResourceArg(resourceType.String())
	log.Printf("resourceType:%s groupResource:%s\n", resource, groupResource)
	if resource == nil {
		panic("parse resource type err")
	}

	c.liters[resourceType.String()] = c.factory.ForResource(*resource).Lister()

	informer := c.factory.ForResource(*resource).Informer()
	c.informers[resourceType.String()] = informer

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			addFunc.rawData = obj
			c.queue.Add(addFunc)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			updateFunc.rawData = []interface{}{oldObj, newObj}
			c.queue.Add(updateFunc)
		},
		DeleteFunc: func(obj interface{}) {
			DeleteFunc.rawData = obj
			c.queue.Add(DeleteFunc)
		},
	})

	return c
}

func (c *QueueController) Run(workerNum int, stopChan chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShuttingDown()

	log.Println("start consume queue")
	go c.factory.Start(stopChan)

	//等待首次同步全部完成
	for _, ok := range c.factory.WaitForCacheSync(stopChan) {
		if !ok {
			log.Printf("timeout wait for cache to sync\n")
			return errors.New("timeout wait for cache to sync")
		}
	}

	for i := 0; i < workerNum; i++ {
		go wait.Until(c.runWorker, time.Second, stopChan)
	}
	<-stopChan

	log.Println("waiting for processing items to finish...")
	c.wg.Wait()

	return nil
}

func (c *QueueController) runWorker() {
	for c.processNextItem() {

	}
}

func (c *QueueController) processNextItem() bool {
	obj, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	c.wg.Add(1)
	defer c.queue.Done(obj)
	result := c.processItem(obj)
	c.handleError(obj, result)
	c.wg.Done()

	return true
}

//用于同步处理某一个数据 TODO 同步处理数据
func (c *QueueController) processItem(obj interface{}) error {
	if deal, ok := obj.(*HandleFunc); ok {
		if err := deal.handle(deal.rawData); err != nil {
			log.Printf("run handle func err:%#v\n", err)
			return err
		}
	}
	return nil
}

func (c QueueController) handleError(obj interface{}, result error) {
	if result == nil {
		//执行成功后清除重试的记录
		c.queue.Forget(obj)
		return
	}

	for c.queue.NumRequeues(obj) < c.maxRetries {
		//执行失败，重试
		c.queue.AddRateLimited(obj)
		log.Printf("处理失败，重新添加队列测试:%#v\n",obj)
		return
	}
	c.queue.Forget(obj)
	runtime.HandleError(fmt.Errorf("max retries exceeded dropping item %+v out of the queue: %v\n", obj, result))
}
