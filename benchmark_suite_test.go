package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestQuicGoBenchmark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Benchmark Suite", []Reporter{newBenchmarkReporter()})
}

var (
	samples int
	sizeMB  int

	netemAvailable bool
)

var _ = BeforeSuite(clearNetem)
var _ = AfterSuite(clearNetem)

func init() {
	flag.IntVar(&samples, "samples", 3, "number of samples")
	flag.IntVar(&sizeMB, "size", 10, "size of transfered data (in MB)")
	flag.Parse()

	if _, err := exec.LookPath("tc"); err != nil {
		fmt.Println("WARNING: This benchmark suite requires netem!")
	} else {
		netemAvailable = true
	}
}

type networkCondition struct {
	Description string
	Command     string
}

var conditions = []networkCondition{
	{Description: "direct transfer"},
	{Description: "5ms RTT", Command: "tc qdisc add #device root netem delay 2.5ms"},
	{Description: "10ms RTT", Command: "tc qdisc add #device root netem delay 5ms"},
	{Description: "25ms RTT", Command: "tc qdisc add #device root netem delay 12.5ms"},
	{Description: "50ms RTT", Command: "tc qdisc add #device root netem delay 25ms"},
	{Description: "100ms RTT", Command: "tc qdisc add #device root netem delay 50ms"},
}

func execNetem(cmd string) string {
	if len(cmd) == 0 {
		return ""
	}
	r := strings.NewReplacer("#device", "dev lo")
	cmd = r.Replace(cmd)
	command := exec.Command("/bin/sh", "-c", "sudo "+cmd)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))
	return string(session.Out.Contents())
}

func clearNetem() {
	if !netemAvailable {
		return
	}
	status := execNetem("tc qdisc show #device")
	if strings.Contains(status, "netem") {
		execNetem("tc qdisc del #device root")
	}
}
