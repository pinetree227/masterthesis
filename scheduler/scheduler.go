package main

import (
        "os"

        //"k8s.io/component-base/cli"
        _ "k8s.io/component-base/logs/json/register"
        _ "k8s.io/component-base/metrics/prometheus/clientgo"
        _ "k8s.io/component-base/metrics/prometheus/version" // for version metric registration
        "k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {

        command := app.NewSchedulerCommand(
                app.WithPlugin("SamplePlugin", NewSamplePlugin),
        )

        if err := command.Execute(); err != nil {
                os.Exit(1)
        }
}

