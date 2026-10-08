// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package devicestats

import (
	"log/slog"

	"github.com/cilium/hive/cell"
	"github.com/cilium/statedb"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/cilium/cilium/pkg/datapath/tables"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/metrics"
)

var Cell = cell.Module(
	"device-metrics",
	"Exports generic link statistics for Cilium-managed devices to Prometheus",
	cell.Invoke(registerCollector),
)

// deviceLister supplies the names of interfaces to collect stats for. It
// decouples the collector from statedb so unit tests can use a fixed list.
type deviceLister interface {
	Names() []string
}

type statedbDeviceLister struct {
	db    *statedb.DB
	table statedb.Table[*tables.Device]
}

func (l statedbDeviceLister) Names() []string {
	var names []string
	for dev := range l.table.List(l.db.ReadTxn(), tables.DevicesBySelected(true)) {
		names = append(names, dev.Name)
	}
	return names
}

type collector struct {
	reader  Reader
	devices deviceLister
	desc    *prometheus.Desc
}

func newCollector(reader Reader, devices deviceLister) *collector {
	return &collector{
		reader:  reader,
		devices: devices,
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(metrics.Namespace, "device", "stat"),
			"Generic link statistic value, labeled by device and stat name",
			[]string{"device", "stat"}, nil,
		),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, dev := range c.devices.Names() {
		stats, err := c.reader.Stats(dev)
		if err != nil {
			continue
		}
		for name, value := range stats {
			ch <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, float64(value), dev, name)
		}
	}
}

func registerCollector(logger *slog.Logger, db *statedb.DB, table statedb.Table[*tables.Device]) {
	c := newCollector(NewNetlinkReader(), statedbDeviceLister{db: db, table: table})
	if err := metrics.Register(c); err != nil {
		logger.Error("Failed to register device metrics collector", logfields.Error, err)
	}
}
