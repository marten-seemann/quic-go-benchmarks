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

var samples int

func init() {
	flag.IntVar(&samples, "samples", 3, "number of samples")
	flag.Parse()
}
