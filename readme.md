###### k8s-common-operate
> 一些常用的k8s操作封装
```go
package informer

import (
	"fmt"
	"github.com/ruanlianjun/k8s-operate/common"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"testing"
)

func TestInformer(t *testing.T) {
	//NewController()
	stopChan:=make(chan struct{})
	dyInformer:=NewInformerCli(DyK8sCli(),0,10,"default")
	err:=dyInformer.Controller.AddEventHandler(common.ResourceKindDeployment,WarpHandleFunc(func(rawData interface{}) error {
		fmt.Printf("----------->addFunc:%+v\n",rawData)
		return nil
	}),WarpHandleFunc(func(rawData interface{}) error {
		if data,ok:=rawData.([]interface{}) ;ok{
			fmt.Printf("----------->updateFunc:\n"+
				"------------->old=%+v \n " +
				"------------->new:%+v\n",data[0],data[1])
		}
		fmt.Printf("update error\n")
		return nil
	}),WarpHandleFunc(func(rawData interface{}) error {
		fmt.Printf("----------->deleteFunc:%+v\n",rawData)
		return nil
	})).Run(10,stopChan)//最大处理数据10的协程
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