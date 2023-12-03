package main

import (
	"flag"
	"fmt"
//	"os"
//	"path/filepath"
	"strconv"
//	"strings"
	"context"
	"sync"
	"math"
	"regexp"
	"k8s.io/client-go/tools/clientcmd"
        "k8s.io/client-go/util/homedir"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/tools/clientcmd"
)
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
	clientset, err := kubernetes.NewForConfig(config)
        if err != nil {
                panic(err.Error())
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

	average := sum / float64(len(pods.Items))
	variance := sumSquared/float64(len(pods.Items)) - average*average

	fmt.Printf("Sum: %f\n", sum)
	fmt.Printf("Average: %f\n", average)
	fmt.Printf("Variance: %f\n", variance)
}


func extractNumbers(inputString string) (int, int, error) {
        // 正規表現のパターン
        pattern := `(\d{2})(\d{2})$`

        // 正規表現に一致するか確認
        re := regexp.MustCompile(pattern)
        matches := re.FindStringSubmatch(inputString)

        if len(matches) != 3 {
                return 0, 0, fmt.Errorf("数値に変換できません")

        }

        // 文字列から数値に変換
        firstTwo, err := strconv.Atoi(matches[1])
        if err != nil {
                return 0, 0, fmt.Errorf("数値に変換できません: %v", err)
        }

        lastTwo, err := strconv.Atoi(matches[2])
        if err != nil {
                return 0, 0, fmt.Errorf("数値に変換できません: %v", err)
        }

        return firstTwo, lastTwo, nil
}


