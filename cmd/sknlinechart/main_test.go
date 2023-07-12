package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/test"
	"reflect"
	"testing"
)

func TestLimitExceededOnNew(t *testing.T) {
	_, err := makeChart("Testing", "Through Widget")
	if err.Error() != "\n::NewLineChart() dataPoint contents exceeds the point count limit[Action: truncated leading]. Series: first, points: 124, Limit: 120" {
		t.Errorf(err.Error())
	}
}

func TestIsLineChartType(t *testing.T) {
	lc, _ := makeChart("Testing", "Through Widget")
	if reflect.TypeOf(lc).String() != "*sknlinechart.LineChartSkn" {
		t.Errorf("not the expected type %v", lc)
	}
}

func TestMinimumSetOfObjects(t *testing.T) {
	lc, _ := makeChart("Testing", "Through Widget")
	_ = lc.MinSize()
	cnt := len(lc.ObjectsCache)
	fmt.Println(cnt)
	if cnt != 56 {
		t.Errorf("Wrong number of base objects %v", cnt)
	}
}

func TestMinSize(t *testing.T) {
	lc, _ := makeChart("Testing", "Through Widget")
	size := lc.MinSize()
	if size == fyne.NewSize(436, 331) {
		t.Errorf("minimum size changed %v", size)
	}
}

func TestTopLeftLabel(t *testing.T) {
	lc, _ := makeChart("Testing", "Through Widget")
	newText := "TESTING"
	lc.SetTopLeftLabel(newText)
	if lc.GetTopLeftLabel() != newText {
		t.Errorf("top left text failed to set. expected:%s, actual: %s", newText, lc.GetTopLeftLabel())
	}
}
func TestTitle(t *testing.T) {
	lc, _ := makeChart("Testing", "Through Widget")
	newText := "Testing"
	if lc.GetTitle() != newText {
		t.Errorf("top centered text failed to set. expected:%s, actual: %s", newText, lc.GetTitle())
	}
}
func TestTopRightLabel(t *testing.T) {
	lc, _ := makeChart("Testing", "Through Widget")
	newText := "TESTING"
	lc.SetTopRightLabel(newText)
	if lc.GetTopRightLabel() != newText {
		t.Errorf("top right text failed to set. expected:%s, actual: %s", newText, lc.GetTopRightLabel())
	}
}

func TestBottomLeftLabel(t *testing.T) {
	lc, _ := makeChart("Testing", "Through Widget")
	newText := "TESTING"
	lc.SetBottomLeftLabel(newText)
	if lc.GetBottomLeftLabel() != newText {
		t.Errorf("bottom left text failed to set. expected:%s, actual: %s", newText, lc.GetBottomLeftLabel())
	}
}
func TestBottomCenteredLabel(t *testing.T) {
	lc, _ := makeChart("Testing", "Through Widget")
	newText := "Through Widget"
	if lc.GetBottomCenteredLabel() != newText {
		t.Errorf("bottom centered text failed to set. expected:%s, actual: %s", newText, lc.GetBottomCenteredLabel())
	}
}
func TestBottomRightLabel(t *testing.T) {
	lc, _ := makeChart("Testing", "Through Widget")
	newText := "TESTING"
	lc.SetBottomRightLabel(newText)
	if lc.GetBottomRightLabel() != newText {
		t.Errorf("bottom right text failed to set. expected:%s, actual: %s", newText, lc.GetBottomRightLabel())
	}
}
