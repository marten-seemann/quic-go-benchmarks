package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQuicGoBenchmark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Benchmark Suite", []Reporter{reporter})
}

var (
	samples  int
	sizeMB   int
	reporter *myReporter
)

func init() {
	flag.IntVar(&samples, "samples", 3, "number of samples")
	flag.IntVar(&sizeMB, "size", 10, "size of transfered data (in MB)")
	flag.Parse()

	reporter = &myReporter{}
}

var _ = AfterSuite(func() {
	reporter.printResult()
})

type measurementSeries map[string]*types.SpecMeasurement

type myReporter struct {
	Reporter

	results map[string]measurementSeries
}

var _ Reporter = &myReporter{}

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
