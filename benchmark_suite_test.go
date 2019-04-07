package main

import (
	"flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCrypto(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Benchmark Suite")
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
