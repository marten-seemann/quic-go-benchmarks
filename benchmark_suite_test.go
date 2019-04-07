package main

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQuicGoBenchmark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Benchmark Suite", []Reporter{newBenchmarkReporter()})
}

var (
	samples int
	sizeMB  int
)

func init() {
	flag.IntVar(&samples, "samples", 3, "number of samples")
	flag.IntVar(&sizeMB, "size", 10, "size of transfered data (in MB)")
	flag.Parse()
}
