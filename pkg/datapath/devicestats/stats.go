// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package devicestats

// Reader fetches statistics for a network device
type Reader interface {
	Stats(iface string) (map[string]uint64, error)
}
