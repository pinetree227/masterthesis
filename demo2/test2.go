package main

import (
	"flag"
	"fmt"
	"time"
//	"os"
//	"path/filepath"
	"strconv"
//	"strings"
	"context"
//	"sync"
//	"math"
//	"regexp"
	"k8s.io/client-go/tools/clientcmd"
        "k8s.io/client-go/util/homedir"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
//	"image/color"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"

	"gonum.org/v1/plot"
        "gonum.org/v1/plot/plotter"
        "gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
//	"k8s.io/client-go/tools/clientcmd"
)

type Car struct {
X string
Y string
Apptype string
Node string
Name string
}

func SeparateData(df dataframe.DataFrame, specie string) dataframe.DataFrame {
	// Filter
	DF := df.Filter(
		dataframe.F{
			Colname:    "Node",
			Comparator: series.Eq,
			Comparando: specie,  // "setosa"
		},
	)

	return DF
}

func makeScatterPlot(df dataframe.DataFrame)(*plotter.Scatter,*plotter.Labels) {
	records := df.Records()
	n := len(records)
	pts := make(plotter.XYs, n-1)
	names := make([]string,n-1)
	for i, r := range records {
        if i == 0 {
			// Skip colname
			continue
		}else if i == n {
			// len(records) == n but len(pts) == n-1
			break
		}
 
		// str to float64
		pts[i-1].X, _ = strconv.ParseFloat(r[0], 64)
		pts[i-1].Y, _ = strconv.ParseFloat(r[1], 64)
		names[i-1] = r[4]

    }
	// fmt.Println(pts)
	textPoints := plotter.XYLabels{
		XYs: pts,
		Labels: names,
	}
	lab, err := plotter.NewLabels(textPoints)
	s, err := plotter.NewScatter(pts)
	if err != nil {
panic(err.Error())
	}
 
	return s,lab
}

func SaveScatterPlot(df1, df2, df3, df4 dataframe.DataFrame, species []string,num int) {
	// Create a new plot
	p := plot.New()

	// Set its title and axis labels
	p.Title.Text = "demo"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"
	p.Add(plotter.NewGrid())

	// Make a scatter plotter
	sp1,l1 := makeScatterPlot(df1)
	sp2,l2 := makeScatterPlot(df2)
	sp3,l3 := makeScatterPlot(df3)
	sp4,l4 := makeScatterPlot(df4)

	// Set color with "gonum.org/v1/plot/plotutil"
	sp1.GlyphStyle.Color = plotutil.Color(0)
	sp2.GlyphStyle.Color = plotutil.Color(1)
	sp3.GlyphStyle.Color = plotutil.Color(2)
	sp4.GlyphStyle.Color = plotutil.Color(3)

	// Set color with "image/color"
	// sp1.GlyphStyle.Color = color.RGBA{R: 128, G: 255, B: 255, A: 128}
	// sp2.GlyphStyle.Color = color.RGBA{R: 255, G: 128, B: 255, A: 128}
	// sp2.GlyphStyle.Color = color.RGBA{R: 255, G: 255, B: 128, A: 128}

	// Set shape
	sp1.Shape = &draw.CircleGlyph{}
	sp2.Shape = &draw.PyramidGlyph{}
	sp3.Shape = &draw.BoxGlyph{}
	sp4.Shape = &draw.RingGlyph{}

	// Add the plotters to the plot, with a legend
	p.Add(sp1)
	p.Add(sp2)
    p.Add(sp3)
    p.Add(sp4)
    p.Add(l1)
    p.Add(l2)
    p.Add(l3)
    p.Add(l4)
    p.Legend.Add(species[0], sp1)
    p.Legend.Add(species[1], sp2)
    p.Legend.Add(species[2], sp3)
    p.Legend.Add(species[3], sp4)
	// Set the range of the axis
    p.X.Min = 0
    p.X.Max = 100
    p.Y.Min = 0
    p.Y.Max = 100

	// Save the plot to a PNG file
if err := p.Save(4*vg.Inch, 4*vg.Inch, "/home/vboxuser/Desktop/images/ScatterPlot"+strconv.Itoa(num)+".png"); err != nil {
	panic(err.Error())
	}
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
	clientset, err := kubernetes.NewForConfig(config)
        if err != nil {
                panic(err.Error())
        }
	num := 0
	for {
	
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	
	if len(pods.Items) == 0{
		break
	}
	var p []Car
	for _, pod := range pods.Items {
		for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodScheduled {
			newCar := Car{pod.Labels["podx"],pod.Labels["pody"],pod.Labels["apptype"],pod.Spec.NodeName,pod.Name[len(pod.Name)-2:]}
			p = append(p,newCar)
		}
	}
}
DF := dataframe.LoadStructs(p)
fmt.Println(DF)
	// Read file
	/*
	f, err := os.Open("./files/iris.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

	// File to Dataframe
	df := dataframe.ReadCSV(f)
	fmt.Println(df)
	fmt.Println(df.Describe())

	// Select {"sepal_length", "sepal_width", "species"}
	DF := Preprocess(df)
	fmt.Println(DF)
*/

	// Separate data by species
	species := []string{"worker2525", "worker2575", "worker7525","worker7575"}
	DF1 := SeparateData(DF, species[0])
	DF2 := SeparateData(DF, species[1])
	DF3 := SeparateData(DF, species[2])
	DF4 := SeparateData(DF, species[3])
	// fmt.Println(seDF)
	// fmt.Println(veDF)
	//fmt.Println(DF1)
	// Save the plot to a PNG file
	SaveScatterPlot(DF1, DF2, DF3,DF4, species,num)
	time.Sleep(20*time.Second)
	num += 1
}
}
