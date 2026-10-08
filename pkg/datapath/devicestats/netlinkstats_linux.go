// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package devicestats

import (
	"github.com/cilium/cilium/pkg/datapath/linux/safenetlink"
)

type netlinkReader struct{}

// NewNetlinkReader returns a Reader backed by rtnl_link_stats64, the
// generic per-device counters available on every link type.
func NewNetlinkReader() Reader {
	return netlinkReader{}
}

func (netlinkReader) Stats(iface string) (map[string]uint64, error) {
	link, err := safenetlink.LinkByName(iface)
	if err != nil {
		return nil, err
	}

	s := link.Attrs().Statistics
	if s == nil {
		return map[string]uint64{}, nil
	}

	return map[string]uint64{
		"rx_packets": s.RxPackets,
		"tx_packets": s.TxPackets,
		"rx_bytes":   s.RxBytes,
		"tx_bytes":   s.TxBytes,
	}, nil
}
