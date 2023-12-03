package main

import (
        "os"
//	"flag"
	"fmt"
        //"k8s.io/component-base/cli"
        _ "k8s.io/component-base/logs/json/register"
        _ "k8s.io/component-base/metrics/prometheus/clientgo"
        _ "k8s.io/component-base/metrics/prometheus/version" // for version metric registration
        "k8s.io/kubernetes/cmd/kube-scheduler/app"
//	"k8s.io/klog/v2"
)

func main() {
	fmt.Printf("ahaha")
//	klog.InitFlags(nil)
//	flag.Set("alsologtostderr","true")
        command := app.NewSchedulerCommand(
                app.WithPlugin("SamplePlugin", NewSamplePlugin),
       )
//	klog.Info("ahaha")

        if err := command.Execute(); err != nil {
                os.Exit(1)
        }
}

