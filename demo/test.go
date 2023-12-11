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
  //      "sync"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"math/rand"
	pinev1 "github.com/pinetree227/location-ctl/api/ctl/v1"
        "github.com/pinetree227/location-ctl/generated/ctl/clientset/versioned"
)

type MyCustomResource struct {
Name string
X float64
Y float64
Update int
Apptype string
Loadtype string
Deleted int
}

func createCustomResource(clientset *versioned.Clientset, name string) (*MyCustomResource, error) {
    options1 := []string{"A","B","C","D","E","F","G","H"}
    options2 := []string{"realtime","no-realtime"}
    loadtype := options1[rand.Intn(len(options1))]
    apptype := options2[rand.Intn(len(options2))]
    x:=rand.Float64()*100
    y:=rand.Float64()*100

    switch loadtype {
    case "A","B":
	    x = 75.00
    case "C","D":
	    x = 25.00
    case "E","F":
	    y = 75.00
    case "G","H":
	    y = 25.00
    }
    newResource := &MyCustomResource{name,x,y,0,apptype,loadtype,0}
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
        return newResource,fmt.Errorf("failed to create custom resource: %v", err)
    }

    return newResource,nil
}

func updateCustomResource(clientset *versioned.Clientset, myResource *MyCustomResource, mdView *pinev1.LocationCtl)  error {
	myResource.Update += 1
	switch myResource.Loadtype {
	case "A","C":
		myResource.Y += 0.32*10
	case "B","D":
		myResource.Y -= 0.32*10
        case "E","G":
		myResource.X += 0.32*10
        case "F","H":
		myResource.X -= 0.32*10
    }
    if myResource.X < 100 && myResource.Y < 100 && myResource.X > 0 && myResource.Y > 0 {

	    if (myResource.X + 0.32*10 > 50.0 && myResource.X < 50.0) || (myResource.Y + 0.32*10 > 50.0 && myResource.Y < 50.0) || (myResource.X - 0.32*10 < 50.0 && myResource.X > 50.0) || (myResource.Y < 50.0 && myResource.Y > 50.0){
		mdView.Spec.Update = 1
	}
	mdView.Spec.PodX = strconv.FormatFloat(myResource.X,'f',2,64)
	mdView.Spec.PodY = strconv.FormatFloat(myResource.Y,'f',2,64)
	_, err := clientset.CtlV1().LocationCtls("default").Update(context.TODO(), mdView, metav1.UpdateOptions{})
        if err != nil {
                return err
        }
    }else{
	err := clientset.CtlV1().LocationCtls("default").Delete(context.TODO(), myResource.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	myResource.Deleted = 1
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
	var p []*MyCustomResource
	name := "example-custom-resource0"
	crdname := ""
	for i := 1; i <= n; i++ {
        crdname = name + strconv.Itoa(i)
	newResource,err := createCustomResource(clientset, crdname)
	p = append(p,newResource)
	if err != nil {
	panic(err.Error())
	}
	}

//	var wg sync.WaitGroup
//	var mu sync.Mutex
	for {
	start := time.Now()
	j := 0
	fmt.Println(time.Now())
	var mdViewList *pinev1.LocationCtlList
        mdViewList, err := clientset.CtlV1().LocationCtls("default").List(context.TODO(),metav1.ListOptions{})
                if err != nil {
                panic(err.Error())
        }

	for _, v := range p {
//		wg.Add(1)
	if v.Deleted != 1{
		j += 1
		fmt.Println(v)
//	go func(clientset *versioned.Clientset, v *MyCustomResource, mdViewList *pinev1.LocationCtlList){

//		defer  wg.Done()
//		mu.Lock()
//		defer mu.Unlock()
var temp pinev1.LocationCtl
for _, cr := range mdViewList.Items {
	if cr.Name == v.Name {
		temp=cr
		break
	}}
  //            mu.Lock()
    //          defer mu.Unlock()

        err = updateCustomResource(clientset, v,&temp)
        if err != nil {
        panic(err.Error())}
//}(clientset, v,mdViewList)
}}
//wg.Wait()
	if j == 0 {
		break
	}
	elapsed := time.Since(start)
	sleepDuration := time.Second*10 - elapsed
	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}
	fmt.Println(elapsed)
}
}
