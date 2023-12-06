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
	//"path/filepath"
	"time"
	"strconv"

	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"math/rand"
//	"strings"
	pinev1 "github.com/pinetree227/location-ctl/api/ctl/v1"
        "github.com/pinetree227/location-ctl/generated/ctl/clientset/versioned"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type MyCustomResource struct {
name string
x float64
y float64
update int
apptype string
loadtype string
deleted int
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

func updateCustomResource(clientset *versioned.Clientset, myResource *MyCustomResource) (int, error) {
	customResource, err := clientset.CustomResource("ctl", "v1", "default", myResource.Name)
	if err != nil {
		return 0,err
	}

	// 取得したカスタムリソースを変更
	customResource.Spec.Update = 1
	myResource.Update += 1
	switch myResource.loadtype {
	case "A":
		myResource.Y += 0.32
            customResource.Spec.PodY = myResource.Y
	case "B":
		myResource.Y -= 0.32
		customResource.Spec.Y = myResource.Y
        case "C":
		myResource.Y += 0.32
		customResource.Spec.Y = myResource.Y
	case "D":
		myResource.Y -= 0.32
		customResource.Spec.Y = myResource.Y
        case "E":
		myResource.X += 0.32
            customResource.Spec.X = myResource.X
        case "F":
		myResource.X -= 0.32
            customResource.Spec.X = myResource.X
	case "G":
		myResource.X += 0.32
            customResource.Spec.X = myResource.X
        case "H":
		myResource.X -= 0.32
            customResource.Spec.X = myResource.X
    }
    if 0 =< myResource.X =< 100 && 0 =< myResource.Y =< 100 {
	    _, err := clientset.UpdateCustomResource("ctl", "v1", "default", myResource.Name, customResource)
        if err != nil {
                return 0,err
        }
    }else{
	err := clientset.DeleteCustomResource("ctl", "v1", "default", myResource.Name, customResource,metav1.DeleteOptions{})
	if err != nil {
		return 0,err
	}
	myResource.deleted = 1
}
    return 0,nil


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
	n := 25
	var p []*MyCustomResource 
	name := "example-custom-resource"
	crdname := ""
	for i := 1; i <= n; i++ {
        crdname = name + strconv.Itoa(i)
	newResource,err := createCustomResource(clientset, crdname)
	p = append(p,newResource)
	if err != nil {
	panic(err.Error())
	}
	}
        time.Sleep(10*time.Second)
	j := 0
	for {
	start := time.Now()
	j = 0
	fmt.Println(time.Now())
	var delindex []int
	for index, v in range p {
        crdname = name + strconv.Itoa(i)
	if v.deleted != 1{
		j += 1
        del,err = updateCustomResource(clientset, v)
        if err != nil {
        panic(err.Error())
        }}}
	if j == 0 {
		break
	}
	elapsed := time.Since(start)
	sleepDurarion := time.Second - elapsed
	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}
	fmt.Println(elapsed)
}

}
