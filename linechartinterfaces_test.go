package sknlinechart_test

import (
	"fyne.io/fyne/v2/theme"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/skoona/sknlinechart"
	"reflect"
	"time"
)

var _ = Describe("Linechartinterfaces", func() {

	It("Ensure Interfaces are implemented", func() {
		point := sknlinechart.NewLineChartDatapoint(62.4, theme.ColorYellow, time.Now().Format(time.RFC3339))
		dIntf := reflect.TypeOf((*sknlinechart.LineChartDatapoint)(nil)).Elem()

		var dataPoints map[string][]*sknlinechart.LineChartDatapoint
		chart, _ := sknlinechart.NewLineChart("Title", "Footer", &dataPoints)
		cIntf := reflect.TypeOf((*sknlinechart.SknLineChart)(nil)).Elem()

		By("LineChartDatapoint interface should be implemented")
		Expect(reflect.TypeOf(point).Implements(dIntf)).To(BeTrue())

		By("SknLineChart interface should be implemented")
		Expect(reflect.TypeOf(chart).Implements(cIntf)).To(BeTrue())
	})
})
