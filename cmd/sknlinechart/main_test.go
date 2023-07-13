package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/skoona/sknlinechart"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestLimitExceededOnNew(t *testing.T) {
	_, err := makeUI("Testing", "Through Widget", 130)
	if err.Error() != "\n::NewLineChart() dataPoint contents exceeds the point count limit[Action: truncated leading]. Series: Testing, points: 130, Limit: 120" {
		t.Errorf(err.Error())
	}
}

func TestIsLineChartType(t *testing.T) {
	lc, _ := makeUI("Testing", "Through Widget", 10)
	if reflect.TypeOf(lc).String() != "*sknlinechart.LineChartSkn" {
		t.Errorf("not the expected type %v", lc)
	}
}

func TestMinimumSetOfObjects(t *testing.T) {
	lc, _ := makeUI("Testing", "Through Widget", 10)
	_ = lc.MinSize()
	cnt := lc.ObjectCount()
	fmt.Println(cnt)
	if cnt != 56 {
		t.Errorf("Wrong number of base objects %v", cnt)
	}
}

func TestMinSize(t *testing.T) {
	lc, _ := makeUI("Testing", "Through Widget", 10)
	size := lc.MinSize()
	if size != fyne.NewSize(436, 331) {
		t.Errorf("minimum size changed %v", size)
	}
}

func TestTopLeftLabel(t *testing.T) {
	lc, _ := makeUI("Testing", "Through Widget", 10)
	newText := "TESTING"
	lc.SetTopLeftLabel(newText)
	if lc.GetTopLeftLabel() != newText {
		t.Errorf("top left text failed to set. expected:%s, actual: %s", newText, lc.GetTopLeftLabel())
	}
}
func TestTitle(t *testing.T) {
	lc, _ := makeUI("Testing", "Through Widget", 10)
	newText := "Testing"
	if lc.GetTitle() != newText {
		t.Errorf("top centered text failed to set. expected:%s, actual: %s", newText, lc.GetTitle())
	}
}
func TestTopRightLabel(t *testing.T) {
	lc, _ := makeUI("Testing", "Through Widget", 10)
	newText := "TESTING"
	lc.SetTopRightLabel(newText)
	if lc.GetTopRightLabel() != newText {
		t.Errorf("top right text failed to set. expected:%s, actual: %s", newText, lc.GetTopRightLabel())
	}
}

func TestBottomLeftLabel(t *testing.T) {
	lc, _ := makeUI("Testing", "Through Widget", 10)
	newText := "TESTING"
	lc.SetBottomLeftLabel(newText)
	if lc.GetBottomLeftLabel() != newText {
		t.Errorf("bottom left text failed to set. expected:%s, actual: %s", newText, lc.GetBottomLeftLabel())
	}
}
func TestBottomCenteredLabel(t *testing.T) {
	lc, _ := makeUI("Testing", "Through Widget", 10)
	newText := "Through Widget"
	if lc.GetBottomCenteredLabel() != newText {
		t.Errorf("bottom centered text failed to set. expected:%s, actual: %s", newText, lc.GetBottomCenteredLabel())
	}
}
func TestBottomRightLabel(t *testing.T) {
	lc, _ := makeUI("Testing", "Through Widget", 10)
	newText := "TESTING"
	lc.SetBottomRightLabel(newText)
	if lc.GetBottomRightLabel() != newText {
		t.Errorf("bottom right text failed to set. expected:%s, actual: %s", newText, lc.GetBottomRightLabel())
	}
}

func makeUI(title, footer string, points int) (sknlinechart.SknLineChart, error) {
	var dataPoints = map[string][]*sknlinechart.LineChartDatapoint{} // legend, points
	rand.NewSource(1000.0)
	for x := 1; x < points+1; x++ {
		val := rand.Float32() * 75.0
		if val > 75.0 {
			val = 75.0
		} else if val < 30.0 {
			val = 30.0
		}
		point := sknlinechart.NewLineChartDatapoint(val, theme.ColorBlue, time.Now().Format(time.RFC3339))
		dataPoints["Testing"] = append(dataPoints["Testing"], &point)
	}
	lineChart, err := sknlinechart.NewLineChart(title, footer, &dataPoints)
	lineChart.EnableDebugLogging(false)

	return lineChart, err
}
