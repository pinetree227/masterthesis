package main

import (
	"flag"
	"fmt"
//	"os"
//	"path/filepath"
	"strconv"
//	"strings"
	"context"
//	"sync"
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

//	var wg sync.WaitGroup
//	var mu sync.Mutex
	var sum, sum2, count, count2,count3  float64
	count=0
	count2=0
	count3=0
	sum=0
	sum2=0
	for _, pod := range pods.Items {
if pod.Labels["apptype"] == "realtime"{
	count3 += 1.0
}
	}
	for _, pod := range pods.Items {
	//	wg.Add(1)
	//	go func(pod v1.Pod) {
//			defer wg.Done()
	for _, condition := range pod.Status.Conditions {
                if condition.Type == v1.PodReady {
//			count2 += 1
			nodeName := pod.Spec.NodeName
			        value1 := pod.Labels["podx"]
				value2 := pod.Labels["pody"]
        if value1 != "" {
                _, err := strconv.ParseFloat(value1,64)
                if err != nil {
                        return
                }
        if value2 != "" {
                _, err := strconv.ParseFloat(value2,64)
                if err != nil {
                        return
		}}
	 }
			podx, _ := strconv.ParseFloat(value1,64)
		        pody, _ := strconv.ParseFloat(value2,64)
			y,x,err := extractNumbers(nodeName)
			if err != nil {
					fmt.Printf("Error getting last digit from node name %s: %v\n", nodeName, err)
					return
				}

//				mu.Lock()
//				defer mu.Unlock()
				dest := (x - podx) * (x - podx) + (y - pody) * (y - pody)
				sums := math.Sqrt(float64(dest))
				if pod.Labels["apptype"] == "realtime"{
			//	if sums <= 35{
					sum += sums
					sum2 += sums
					count += 1.0
					count2 += 1.0
			//	}
				}else{
			//		if sums <= 105{
					sum += sums
					count2 += 1.0
			//	}
				}
			}}
			//	sumSquared += float64(dest)
//			}(pod)
	}

//	wg.Wait()
//	fmt.Println(count2)
	average := sum / count2
	average2 := sum2 / count
	success := count2/float64(len(pods.Items))
	success2 := count / count3
//	variance := sumSquared/float64(count2) - average*average

//fmt.Printf("Varianceiall: %f\n", variance)
fmt.Printf("%f\n", average)
//fmt.Println(count)
fmt.Printf("%f\n", average2)
fmt.Printf("%f\n", success)
fmt.Printf("%f\n",success2)
}

/*
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
*/
func extractNumbers(inputString string) (float64, float64, error) {
        // 正規表現のパターン
        pattern := `(\d{2})(\d{2})$`

        // 正規表現に一致するか確認
        re := regexp.MustCompile(pattern)
        matches := re.FindStringSubmatch(inputString)

        if len(matches) != 3 {
                return 0, 0, fmt.Errorf("数値に変換できません")
        }

        // 文字列から数値に変換
        firstTwo, err := strconv.ParseFloat(matches[1],64)
        if err != nil {
                return 0, 0, fmt.Errorf("数値に変換できません: %v", err)
        }

        lastTwo, err := strconv.ParseFloat(matches[2],64)
        if err != nil {
                return 0, 0, fmt.Errorf("数値に変換できません: %v", err)
        }

        return firstTwo, lastTwo, nil
}

