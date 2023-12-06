package main
import (
    "context"
    "fmt"
    "regexp"
//    "strings"
    "flag"
    "math"
    "strconv"
    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"
    framework "k8s.io/kubernetes/pkg/scheduler/framework"
    "k8s.io/klog/v2"
)

type SamplePlugin struct{}
/*
var _ framework.FilterPlugin = &SamplePlugin{}

func (pl *SamplePlugin) Name() string {
    return "SamplePlugin"
}

func (pl *SamplePlugin) Filter(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
    if !strings.Contains(pod.Name, "sample") {
        return framework.NewStatus(framework.Error, "Pod name does not contain 'sample'")
    }

    return nil
}

func NewSamplePlugin(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
    return &SamplePlugin{}, nil
}
*/

var _ framework.ScorePlugin = &SamplePlugin{}
var _ framework.PreScorePlugin = &SamplePlugin{}

func (pl *SamplePlugin) Name() string {
    return "SamplePlugin"
}
/*
// preScoreState computed at PreScore and used at Score.
type preScoreState1 struct {
	nodeLatitude int
        nodeLongitude int
}
type preScoreState2 struct {
	podLatitude int
        podLongitude int
}

// Clone implements the mandatory Clone interface. We don't really copy the data since
// there is no need for that.
func (s *preScoreState1) Clone() framework.StateData {
	return s
}
func (s *preScoreState2) Clone() framework.StateData {
        return s
}

func (pl *SamplePlugin) PreScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodes []*v1.Node) *framework.Status {
	for _, node := range nodes {
		value1 := node.Labels["nodeLatitude"]
		if value1 != "" {
			num1, err := strconv.Atoi(value1)
			if err != nil {
				return nil
			}
		value2 := node.Labels["nodeLongitude"]
		if value2 != "" {
			num2, err := strconv.Atoi(value2)
			if err != nil {
				return nil
			}
		s1 := &preScoreState1{
		nodeLatitude: num1,
		nodeLongitude: num2,
		}
		state.Write(string(node).Name, s1)
		}
	}
	}
	value1 := pod.Labels["podLatitude"]
	if value1 != "" {
		num1, err := strconv.Atoi(value1)
		if err != nil {
			return nil
	        }
	value2 := pod.Labels["podLongitude"]
	if value2 != "" {
		num2, err := strconv.Atoi(value2)
		if err != nil {
			return nil
		}
	s2 := &preScoreState2{
		podLatitude: num1,
		podLongitude: num2,
	}
	state.Write(string(pod.Name), s2)
	}
	}
	return nil
}



func (pl *SamplePlugin) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	data1, err := state.Read(string(nodeName))
	if err != nil {
		return 0, framework.AsStatus(err)
	}
	s1, ok := data1.(*preScoreState1)
	if !ok {
		return 0, framework.AsStatus(err)
	}
	data2, err := state.Read(string(pod.Name))
	if err != nil {
		return 0, framework.AsStatus(err)
	}
	s2, ok := data2.(*preScoreState2)
	if !ok {
		return 0, framework.AsStatus(err)
	}
	dest := (s1.nodeLatitude - s2.podLatitude) * (s1.nodeLatitude - s2.podLatitude) + (s1.nodeLongitude - s2.podLongitude) * (s1.nodeLongitude - s2.podLongitude)
	if dest != nil{
		float64(dest) = math.Sqrt(dest)
		return int64(100 - dest), nil
	}
    return 0, nil
}
*/
const Name = "NodeNumber"
const preScoreStateKey = "PreScore" + Name

type preScoreState struct {
	x float64
	y float64
}
func init() {
	// initで呼び出しているが、以下は最初のログを書き出す前であればどこで実行してもよい
    flagset := flag.NewFlagSet("", flag.ContinueOnError)
    klog.InitFlags(flagset)
	// "-v=2" とすると "Hey from 地球" のログエントリーも表示される。
    flagset.Parse([]string{"-v=1"})
}


// Clone implements the mandatory Clone interface. We don't really copy the data since
// there is no need for that.

func (s *preScoreState) Clone() framework.StateData {
	return s
}

func (pl *SamplePlugin) PreScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodes []*v1.Node) *framework.Status {
	value1 := pod.Labels["podx"]
	if value1 != "" {
		num1, err := strconv.ParseFloat(value1,64)
		if err != nil {
			return nil
		}
	value2 := pod.Labels["pody"]
	if value2 != "" {
		num2, err := strconv.ParseFloat(value2,64)
		if err != nil {
			return nil
		}
	s := &preScoreState{
		x: num1,
		y: num2,
	}
	state.Write(preScoreStateKey, s)
	}
}
	return nil
}

func (pl *SamplePlugin) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
/*	la:=nodeName[len(nodeName)-2:]
	if _,err := strconv.Atoi(la); err != nil {
                return 0, nil
        }
	num1,err := strconv.Atoi(la)

	lo:=nodeName[len(nodeName)-4:len(nodeName)-3]
        if _,err := strconv.Atoi(lo); err != nil {
                return 0, nil
        }
	num2,err := strconv.Atoi(lo)
	*/
//	flagset := flag.NewFlagSet("",flag.ContinueOnError)
//	klog.InitFlags(flagset)
 //       logger := klig.FiomContext(ctx)
	defer klog.Flush()
	y, x, err := extractNumbers(nodeName)

	if err != nil {
		fmt.Println("エラー:", err)
		return 0, framework.AsStatus(err)
	}

	if _, err := state.Read(preScoreStateKey); err != nil {
		return 0, framework.AsStatus(err)
	}

	data, err := state.Read(preScoreStateKey)
	if _, ok := data.(*preScoreState); !ok {
		return 0, framework.AsStatus(err)
	}
	s, _ :=data.(*preScoreState);

	dest := (x - s.x) * (x - s.x) + (y - s.y) * (y - s.y)
	if dest > 0 {
		dest2 := math.Sqrt(float64(dest))
		dest3 := 110-int64(dest2)
		dest4 := int64(float64(dest3) * float64(100) / float64(110))
		klog.V(4).InfoS("score:",nodeName,dest4)
		return dest4, nil
	}
	return 0, nil
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

	return firstTwo, lastTwo, nil
}


func (pl *SamplePlugin) ScoreExtensions() framework.ScoreExtensions {
	return nil
}





func NewSamplePlugin(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
    return &SamplePlugin{}, nil
}



