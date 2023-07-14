package sknlinechart_test

import (
	"fyne.io/fyne/v2/theme"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/skoona/sknlinechart"
	"math/rand"
	"time"
)

var _ = Describe("Maps and slices utilities", func() {

	var dataPoints []*sknlinechart.LineChartDatapoint

	BeforeEach(func() {
		rand.NewSource(1000.0)
		for x := 1; x < 11; x++ {
			val := rand.Float32() * 75.0
			if val > 75.0 {
				val = 75.0
			} else if val < 30.0 {
				val = 30.0
			}
			point := sknlinechart.NewLineChartDatapoint(val, theme.ColorBlue, time.Now().Format(time.RFC3339))
			dataPoints = append(dataPoints, &point)
		}
	})

	Describe("shift data points through slice limited to 10 objects", func() {
		var first, last, newOne sknlinechart.LineChartDatapoint
		var orignalCount int

		BeforeEach(func() {
			newOne = sknlinechart.NewLineChartDatapoint(960.13, "TEST", time.Now().Format(time.RFC3339))
			first = *dataPoints[0]
			last = *dataPoints[(len(dataPoints) - 1)]
			orignalCount = len(dataPoints)
			dataPoints = sknlinechart.ShiftSlice(&newOne, dataPoints)
		})

		It("previous last should not equal current last", func() {
			Expect((*dataPoints[(len(dataPoints) - 1)])).NotTo(Equal(last))
		})
		It("first should be removed", func() {
			Expect(dataPoints[0]).NotTo(Equal(first))
		})
		It("first should be removed", func() {
			Expect(dataPoints[0]).NotTo(Equal(first))
		})
		It("last should be equal newOne", func() {
			Expect((*dataPoints[(len(dataPoints) - 1)])).To(Equal(newOne))
		})
		It("slice size should not change", func() {
			Expect(len(dataPoints)).To(Equal(orignalCount))
		})
	})

	Describe("Remove datapoints when given too many on addition to series", func() {
		var orignalCount int

		It("should remove one point from source", func() {
			orignalCount = len(dataPoints)
			dataPoints = sknlinechart.RemoveIndexFromSlice(0, dataPoints)
			Expect(len(dataPoints)).To(Equal(orignalCount - 1))
		})
		It("should detect empty slice and return it empty", func() {
			var a []*sknlinechart.LineChartDatapoint
			var b []*sknlinechart.LineChartDatapoint
			b = sknlinechart.RemoveIndexFromSlice(0, a)
			Expect(len(a)).To(Equal(len(b)))
			Expect(a).To(Equal(b))
		})

	})

})
