package main
import (
    "context"
//    "strings"
"fmt"
"regexp"
"math"
"strconv"
    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"
    framework "k8s.io/kubernetes/pkg/scheduler/framework"
)

type SamplePlugin2 struct{}

var _ framework.FilterPlugin = &SamplePlugin2{}

func (pl *SamplePlugin2) Name() string {
    return "SamplePlugin2"
}

func (pl *SamplePlugin2) Filter(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	if nodeInfo.Node() == nil {
		return framework.NewStatus(framework.Error, "node not found")
	}
	nodeName := nodeInfo.Node().Name
	y,x,err := extractNumbers(nodeName)
	if err != nil {
		return framework.NewStatus(framework.Error, "node name err")
	}
	value1 := pod.Labels["podx"]
        if value1 != "" {
                if _, err := strconv.Atoi(value1); err != nil {
                        return framework.NewStatus(framework.Error, "pod err") 
		}
	}
        value2 := pod.Labels["pody"]
        if value2 != "" {
                if _, err := strconv.Atoi(value2); err != nil {
                        return framework.NewStatus(framework.Error, "pod err")
                }
	}

	podx, _ := strconv.ParseFloat(value1,64)
	pody, _ := strconv.ParseFloat(value2,64)
	dest := (x - podx) * (x - podx) + (y - pody) * (y - pody)
        if dest > 0 {
                dest2 := math.Sqrt(float64(dest))
               // klog.V(4).InfoS("score:",nodeName,dest4)
	               if dest2 > 35{
			       return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("%g : %v",dest2,nodeName))
}
    }
    return nil
//      return framework.NewStatus(framework.Error, value1,value2,nodeName,strconv.Itoa(x),strconv.Itoa(y))

}
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

	return firstTwo*8.66, lastTwo*5, nil
}


func NewSamplePlugin(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
    return &SamplePlugin2{}, nil
}
