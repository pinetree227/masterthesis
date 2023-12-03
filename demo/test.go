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

func createCustomResource(clientset *versioned.Clientset, name string) error {

    x:=rand.Intn(100)
    y:=rand.Intn(100)

    cr := &pinev1.LocationCtl{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
        },
        Spec: pinev1.LocationCtlSpec{
            // Set your custom resource spec fields here
  PodX: strconv.Itoa(x),
  PodY: strconv.Itoa(y),
  Replicas: 1,
        },
    }
    _, err := clientset.CtlV1().LocationCtls("default").Create(context.TODO(), cr, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("failed to create custom resource: %v", err)
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
	n := 160

	name := "example-custom-resource"
	crdname := ""
	for i := 1; i <= n; i++ {
        crdname = name + strconv.Itoa(i)
	err = createCustomResource(clientset, crdname)
	if err != nil {
	panic(err.Error())
	}
}
}
        pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
        if err != nil {
                panic(err.Error())
        }

        var wg sync.WaitGroup
        var mu sync.Mutex
        var sum, sumSquared float64

        for _, pod := range pods.Items {
                wg.Add(1)
                go func(pod v1.Pod) {
                        defer wg.Done()
                        nodeName := pod.Spec.NodeName
                        value1 := pod.Labels["podx"]
                        if value1 != "" {
                                if _, err := strconv.Atoi(value1); err != nil {
                                        return
                                }
                        }
                        value2 := pod.Labels["pody"]
                        if value2 != "" {
                                if _, err := strconv.Atoi(value2); err != nil {
                                 return
                         }
                 }
                        podx, _ := strconv.Atoi(value1)
                        pody, _ := strconv.Atoi(value2)
                        y,x,err := extractNumbers(nodeName)
                        if err != nil {
                                        fmt.Printf("Error getting last digit from node name %s: %v\n", nodeName, err)
                                        return
                                }

                                mu.Lock()
                                defer mu.Unlock()
                                dest := (x - podx) * (x - podx) + (y - pody) * (y - pody)
                                sum += math.Sqrt(float64(dest))
                                sumSquared += float64(dest)
                        }(pod)
        }

        wg.Wait()
	time.Sleep(10*time.Second)
