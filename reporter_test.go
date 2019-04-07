package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
)

type measurementSeries map[string]*types.SpecMeasurement

type myReporter struct {
	ginkgo.Reporter

	results map[string]measurementSeries
}

var _ ginkgo.Reporter = &myReporter{}

func (r *myReporter) SpecSuiteWillBegin(config.GinkgoConfigType, *types.SuiteSummary) {}
func (r *myReporter) BeforeSuiteDidRun(*types.SetupSummary)                           {}
func (r *myReporter) SpecWillRun(*types.SpecSummary)                                  {}
func (r *myReporter) AfterSuiteDidRun(*types.SetupSummary)                            {}
func (r *myReporter) SpecSuiteDidEnd(*types.SuiteSummary)                             {}

func (r *myReporter) SpecDidComplete(specSummary *types.SpecSummary) {
	if !specSummary.IsMeasurement {
		return
	}
	method := specSummary.ComponentTexts[2]
	measurement, ok := specSummary.Measurements["runtime"]
	if !ok {
		return
	}
	r.addResult(method, "transfer", measurement)
}

func (r *myReporter) addResult(cond, ver string, measurement *types.SpecMeasurement) {
	if r.results == nil {
		r.results = make(map[string]measurementSeries)
	}
	if _, ok := r.results[cond]; !ok {
		r.results[cond] = make(measurementSeries)
	}
	r.results[cond][ver] = measurement
}

func (r *myReporter) printResult() {
	fmt.Printf("\nBenchmark results:\n")
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"", "TCP", "QUIC"}
	table.SetHeader(header)
	table.SetCaption(true, fmt.Sprintf("Based on %d samples (%d MB).", samples, sizeMB))
	table.SetAutoFormatHeaders(false)
	colAlignments := []int{tablewriter.ALIGN_LEFT}
	for i := 1; i <= len(header); i++ {
		colAlignments = append(colAlignments, tablewriter.ALIGN_RIGHT)
	}
	table.SetColumnAlignment(colAlignments)

	for _, cond := range []string{"transfer"} {
		data := make([]string, len(header))
		data[0] = cond

		for i := 1; i < len(header); i++ {
			measurement := r.results[header[i]][cond]
			var out string
			if measurement == nil {
				out = "-"
			} else {
				if len(measurement.Results) <= 1 {
					out = fmt.Sprintf(" %.2f", measurement.Average)
				} else {
					out = fmt.Sprintf("%.2f ± %.2f", measurement.Average, measurement.StdDeviation)
				}
			}
			data[i] = out
		}
		table.Append(data)
	}
	table.Render()
}
