package test

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/ruanlianjun/k8sct/common"
	"github.com/ruanlianjun/k8sct/watch"
	"go.uber.org/goleak"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	defer goleak.VerifyNone(t)
	go Informer("ai-devel")
	go Informer("ning-test")
	go Informer("cci-namespace-53968331")
	go Informer("cci-namespace-hf-1")
	time.Sleep(time.Second*20)
}

func DyK8sCli() dynamic.Interface {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/ruanlianjun/Desktop/hw_config")

	if err != nil {
		panic(err)
	}
	k, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return k
}

func Informer(namespace string)  {
	stopChan := make(chan struct{})

	go func() {
		time.Sleep(time.Second*10)
		close(stopChan)
	}()

	//defaultResyncPeriod update可能sync信息 maxRetries处理失败后最大重试次数
	dyInformer := watch.NewInformerCli(DyK8sCli(), 0, 10, namespace)
	err := dyInformer.AddEventHandler(common.ResourceKindPod, watch.WarpHandleFunc(func(rawData interface{}) error {
		var pod v1.Pod

		marshal, err := jsoniter.Marshal(rawData)
		if err!=nil {
			panic("raw data to marshal err")
		}
		err = jsoniter.Unmarshal(marshal, &pod)
		if err != nil {
			panic("unmarshal err")
		}

		fmt.Printf("==================>%#v\n",pod)



		fmt.Printf("----------->namespace:%s addFunc:%+v\n", namespace,common.ResourceKindPod)
		return nil
	}), watch.WarpHandleFunc(func(rawData interface{}) error {
		if _, ok := rawData.([]interface{}); ok {
			fmt.Printf("----------->namespace:%s updateFunc: %#v\n",namespace,common.ResourceKindPod)
		}
		return nil
	}), watch.WarpHandleFunc(func(rawData interface{}) error {
		fmt.Printf("----------->namespace:%s deleteFunc:%+v\n", namespace,common.ResourceKindPod)
		return nil
	})).AddEventHandler(common.ResourceKindEvent,watch.WarpHandleFunc(func(rawData interface{}) error {
		fmt.Printf("----------->namespace:%s,  %s addFunc:%+v\n", namespace,common.ResourceKindEvent,rawData)
		return nil
	}),watch.WarpHandleFunc(func(rawData interface{}) error {
		if update,ok:=rawData.([]interface{});ok {
			fmt.Printf("----------->namespace:%s %s updateFunc: old:%#v new:%#v\n",namespace,common.ResourceKindEvent,update[0],update[1])
		}

		return nil
	}),watch.WarpHandleFunc(func(rawData interface{}) error {
		fmt.Printf("----------->namespace:%s %s deleteFunc:%+v\n",namespace, common.ResourceKindEvent,rawData)
		return nil
	})).AddEventHandler(common.ResourceKindJob,watch.WarpHandleFunc(func(rawData interface{}) error {
		fmt.Printf("----------->namespace:%s  %s addFunc:%+v\n",namespace, common.ResourceKindJob,rawData)
		return nil
	}),watch.WarpHandleFunc(func(rawData interface{}) error {
		if update,ok:=rawData.([]interface{});ok {
			fmt.Printf("----------->namespace:%s %s updateFunc: old:%#v new:%#v\n",namespace,common.ResourceKindJob,update[0],update[1])
		}

		return nil
	}),watch.WarpHandleFunc(func(rawData interface{}) error {
		fmt.Printf("----------->namespace:%s %s deleteFunc:%+v\n", namespace,common.ResourceKindJob,rawData)
		return nil
	})).Run(10, stopChan) //最大10个协程处理数据
	if err!=nil {
		panic(err)
	}
}