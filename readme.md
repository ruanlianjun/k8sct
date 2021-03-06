###### k8s-common-operate
> watch informer操作封装

##### 使用
> import "github.com/ruanlianjun/k8sct"
```go
package test

import (
	"fmt"
	"github.com/ruanlianjun/k8sct/common"
	"github.com/ruanlianjun/k8sct/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"testing"
)

func TestWatch(t *testing.T) {
	stopChan:=make(chan struct{})
	//defaultResyncPeriod update可能sync信息 maxRetries处理失败后最大重试次数
	dyInformer:=watch.NewInformerCli(DyK8sCli(),0,10,"default")
	err:=dyInformer.AddEventHandler(common.ResourceKindConfigMap,watch.WarpHandleFunc(func(rawData interface{}) error {
		fmt.Printf("----------->addFunc:%+v\n",rawData)
		return nil
	}),watch.WarpHandleFunc(func(rawData interface{}) error {
		if data,ok:=rawData.([]interface{}) ;ok{
			fmt.Printf("----------->updateFunc:\n"+
				"------------->old=%+v \n " +
				"------------->new:%+v\n",data[0],data[1])
		}
		fmt.Printf("update error\n")
		return nil
	}),watch.WarpHandleFunc(func(rawData interface{}) error {
		fmt.Printf("----------->deleteFunc:%+v\n",rawData)
		return nil
	})).Run(10,stopChan)//最大10个协程处理数据
	log.Println("informer:",err)
}

func DyK8sCli() dynamic.Interface {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/ruanlianjun/.kube/config")

	if err != nil {
		panic(err)
	}
	k, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return k
}
```