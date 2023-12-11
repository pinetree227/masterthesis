/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"
	"time"
	"strconv"
       "sync"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"math/rand"
	pinev1 "github.com/pinetree227/location-ctl/api/ctl/v1"
        "github.com/pinetree227/location-ctl/generated/ctl/clientset/versioned"
)



func createCustomResource(clientset *versioned.Clientset, name string) error {
    //options1 := []string{"A","B","C","D","E","F","G","H"}
    options2 := []string{"realtime","no-realtime"}
    apptype := options2[rand.Intn(len(options2))]
    x:=rand.Float64()*100
    y:=rand.Float64()*100

//    newResource := &MyCustomResource{name,x,y,0,apptype}
    cr := &pinev1.LocationCtl{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
        },
        Spec: pinev1.LocationCtlSpec{
            // Set your custom resource spec fields here
  PodX: strconv.FormatFloat(x,'f',2,64),
  PodY: strconv.FormatFloat(y,'f',2,64),
  Update: 0,
  Apptype: apptype,
  Replicas: 1,
        },
    }
    _, err := clientset.CtlV1().LocationCtls("default").Create(context.TODO(), cr, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("failed to create custom resource: %v", err)
    }

    return nil
}

func updateCustomResource(clientset *versioned.Clientset, mdView *pinev1.LocationCtl)  error {
	mdView.Spec.PodX = strconv.FormatFloat(rand.Float64()*100,'f',2,64)
	mdView.Spec.PodY = strconv.FormatFloat(rand.Float64()*100,'f',2,64)
	if   rand.Float64()  < 0.05{
		mdView.Spec.Update = 1
	}
	_, err := clientset.CtlV1().LocationCtls("default").Update(context.TODO(), mdView, metav1.UpdateOptions{})
        if err != nil {
                return err
        }
return nil

}


func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", "/etc/kubernetes/admin.conf", "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	n := 100
	//var p []*MyCustomResource
	name := "example-custom-resource0"
	crdname := ""
	for i := 1; i <= n; i++ {
        crdname = name + strconv.Itoa(i)
	err := createCustomResource(clientset, crdname)
	if err != nil {
	panic(err.Error())
	}
	}

	var wg sync.WaitGroup
//	var mu sync.Mutex
	for {
		        var mdViewList *pinev1.LocationCtlList
        mdViewList, err = clientset.CtlV1().LocationCtls("default").List(context.TODO(),metav1.ListOptions{})
                if err != nil {
                panic(err.Error())
        }

	var timesum time.Duration
	start := time.Now()
//		mu.Lock()
//		defer mu.Unlock()
for _, cr := range mdViewList.Items {
	       wg.Add(1)

	        go func(clientset *versioned.Clientset, cr pinev1.LocationCtl){
		defer  wg.Done()

  //            mu.Lock()
    //          defer mu.Unlock()
    start2 := time.Now()
        err = updateCustomResource(clientset, &cr)
        if err != nil {
        panic(err.Error())}
	elapsed2 := time.Since(start2)
	timesum += elapsed2
}(clientset, cr)
}
        elapsed := time.Since(start)
        sleepDuration := time.Second - elapsed
        if sleepDuration > 0 {
                time.Sleep(sleepDuration)
        }
	average := timesum / time.Duration(n)
	fmt.Println(average)
wg.Wait()
}
}
