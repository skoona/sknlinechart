package sknlinechart_test

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/skoona/sknlinechart"
	"reflect"
	"time"
)

var _ = Describe("Datapoint Operations", func() {
	It("should return a valid datapoint object", func() {
		point := sknlinechart.NewChartDatapoint(62.3, theme.ColorYellow, time.Now().Format(time.RFC1123))
		Expect(reflect.TypeOf(point).String()).To(Equal("*sknlinechart.chartDatapoint"))
	})
	It("properties should respond as expected", func() {
		point := sknlinechart.NewChartDatapoint(62.4, theme.ColorYellow, time.Now().Format(time.RFC1123))

		By("should initialize the marker positions to empty")
		a, b := point.MarkerPosition()
		Expect(*a).To(Equal(fyne.Position{}))
		Expect(*b).To(Equal(fyne.Position{}))

		By("should have string values set")
		Expect(point.ColorName()).NotTo(BeEmpty())
		Expect(point.Timestamp()).NotTo(BeEmpty())
		Expect(point.ExternalID()).NotTo(BeEmpty())

		By("should have float32 value set from new")
		Expect(point.Value()).To(BeNumerically("==", float32(62.4)))

		By("should return a valid copy of current datapoint")
		Expect(point.Copy()).To(Equal(point))

		By("should be able to change string values")
		point.SetColorName(theme.ColorBlue)
		value := time.Now().Format(time.RFC1123)
		point.SetTimestamp(value)
		Expect(point.ColorName()).To(Equal(theme.ColorBlue))
		Expect(point.Timestamp()).To(Equal(value))

		By("should be able to change the float32 value")
		point.SetValue(77.12)
		Expect(point.Value()).To(BeNumerically("==", float32(77.12)))

		By("should be able set marker positions")
		c := fyne.NewPos(12, 12)
		d := fyne.NewPos(20, 20)
		point.SetMarkerPosition(&c, &d)
		a, b = point.MarkerPosition()
		Expect(*a).To(Equal(c))
		Expect(*b).To(Equal(d))
	})

})
