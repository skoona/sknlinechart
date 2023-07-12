package main

import (
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
	lc.SetMinSize(fyne.NewSize(500, 500))
	cnt := len(lc.ObjectsCache)
	if cnt == 0 {
		t.Errorf("Wrong number of base objects %v", cnt)
	}
}
