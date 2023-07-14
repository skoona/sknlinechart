package sknlinechart_test

import (
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/skoona/sknlinechart"
	"math/rand"
	"reflect"
	"time"
)

var _ = Describe("verify line chart initial state", func() {

	It("should accept minimum number of points on create", func() {
		lc, err := makeUI("Testing", "Through Widget", 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(reflect.TypeOf(lc).String()).To(Equal("*sknlinechart.LineChartSkn"))
	})
	It("should accept maximum number of points on create", func() {
		lc, err := makeUI("Testing", "Through Widget", 120)
		Expect(err).NotTo(HaveOccurred())
		Expect(reflect.TypeOf(lc).String()).To(Equal("*sknlinechart.LineChartSkn"))
	})
	It("should accept overflow number of points on create", func() {
		lc, err := makeUI("Testing", "Through Widget", 130)
		Expect(err).To(HaveOccurred())
		Expect(reflect.TypeOf(lc).String()).To(Equal("*sknlinechart.LineChartSkn"))
	})
	It("should cache chart elements that don't change value or position", func() {
		lc, _ := makeUI("Testing", "Through Widget", 2)
		lc.Refresh()
		actual := lc.ObjectCount()
		Expect(actual).To(Equal(56))
	})
	It("should support a usable minimum size", func() {
		lc, _ := makeUI("Testing", "Through Widget", 2)
		actual := lc.MinSize()
		Expect(actual.Width).To(BeNumerically("==", float32(436.0)))
	})

	It("chart border labels can be changed", func() {
		lc, _ := makeUI("Testing", "Through Widget", 2)

		By("Setting a new value for top left label")
		oldValue := lc.GetTopLeftLabel()
		lc.SetTopLeftLabel("TESTING")
		Expect(lc.GetTopLeftLabel()).NotTo(Equal(oldValue))

		By("Setting a new value for title")
		oldValue = lc.GetTitle()
		lc.SetTitle("TESTING")
		Expect(lc.GetTitle()).NotTo(Equal(oldValue))

		By("Setting a new value for top right label")
		oldValue = lc.GetTopRightLabel()
		lc.SetTopRightLabel("TESTING")
		Expect(lc.GetTopRightLabel()).NotTo(Equal(oldValue))

		By("Setting a new value for middle left label")
		oldValue = lc.GetMiddleLeftLabel()
		lc.SetMiddleLeftLabel("TESTING")
		Expect(lc.GetMiddleLeftLabel()).NotTo(Equal(oldValue))

		By("Setting a new value for middle right label")
		oldValue = lc.GetMiddleRightLabel()
		lc.SetMiddleRightLabel("TESTING")
		Expect(lc.GetMiddleRightLabel()).NotTo(Equal(oldValue))

		By("Setting a new value for bottom left label")
		oldValue = lc.GetBottomLeftLabel()
		lc.SetBottomLeftLabel("TESTING")
		Expect(lc.GetBottomLeftLabel()).NotTo(Equal(oldValue))

		By("Setting a new value for bottom center label")
		oldValue = lc.GetBottomCenteredLabel()
		lc.SetBottomCenteredLabel("TESTING")
		Expect(lc.GetBottomCenteredLabel()).NotTo(Equal(oldValue))

		By("Setting a new value for bottom right label")
		oldValue = lc.GetBottomRightLabel()
		lc.SetBottomRightLabel("TESTING")
		Expect(lc.GetBottomRightLabel()).NotTo(Equal(oldValue))
	})

})

func makeUI(title, footer string, points int) (sknlinechart.LineChart, error) {
	var dataPoints = map[string][]*sknlinechart.ChartDatapoint{} // legend, points
	if points != 0 {
		rand.NewSource(1000.0)
		for x := 1; x < points+1; x++ {
			val := rand.Float32() * 75.0
			if val > 75.0 {
				val = 75.0
			} else if val < 30.0 {
				val = 30.0
			}
			point := sknlinechart.NewChartDatapoint(val, theme.ColorBlue, time.Now().Format(time.RFC3339))
			dataPoints["Testing"] = append(dataPoints["Testing"], &point)
		}
	}
	lineChart, err := sknlinechart.NewLineChart(title, footer, &dataPoints)
	lineChart.EnableDebugLogging(false)

	return lineChart, err
}
