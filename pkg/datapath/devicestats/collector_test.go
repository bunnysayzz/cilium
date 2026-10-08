// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package devicestats

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

type fakeReader map[string]map[string]uint64

func (f fakeReader) Stats(iface string) (map[string]uint64, error) {
	return f[iface], nil
}

type fakeDevices []string

func (f fakeDevices) Names() []string { return f }

func TestCollector(t *testing.T) {
	reader := fakeReader{
		"eth0": {"rx_packets": 10, "tx_packets": 20},
	}
	c := newCollector(reader, fakeDevices{"eth0"})

	want := `
# HELP cilium_device_stat Generic link statistic value, labeled by device and stat name
# TYPE cilium_device_stat counter
cilium_device_stat{device="eth0",stat="rx_packets"} 10
cilium_device_stat{device="eth0",stat="tx_packets"} 20
`
	if err := testutil.CollectAndCompare(c, strings.NewReader(want)); err != nil {
		t.Fatal(err)
	}
}
