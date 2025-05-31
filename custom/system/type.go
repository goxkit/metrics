// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package system provides system metrics collection capabilities for monitoring
// memory usage, garbage collection, threads, and goroutines. It offers comprehensive
// tooling for observing the runtime behavior of Go applications in production environments.
package system

import (
	"go.opentelemetry.io/otel/metric"
)

type (
	// BasicGauges defines an interface for metrics collectors that gather
	// system-level metrics using OpenTelemetry observable gauges. It abstracts
	// the common functionality needed for different types of system metrics collectors.
	BasicGauges interface {
		// Collect registers callbacks for the metrics with the provided meter.
		// This sets up the continuous collection of metrics data from the system.
		Collect(meter metric.Meter)
	}

	// memGauges implements BasicGauges to collect memory-related metrics.
	// It contains observable gauges for various memory statistics including
	// heap allocation, garbage collection, and system memory usage.
	// These metrics are essential for monitoring memory utilization patterns and
	// identifying potential memory leaks or inefficient memory usage.
	memGauges struct {
		// System memory metrics
		ggSysBytes          metric.Int64ObservableGauge // Total bytes obtained from system
		ggAllocBytesTotal   metric.Int64ObservableGauge // Total bytes allocated, even if freed
		ggHeapAllocBytes    metric.Int64ObservableGauge // Bytes allocated and still in use
		ggFreesTotal        metric.Int64ObservableGauge // Total count of frees
		ggGcSysBytes        metric.Int64ObservableGauge // Bytes used for garbage collection system metadata
		ggHeapIdleBytes     metric.Int64ObservableGauge // Bytes in idle spans
		ggInuseBytes        metric.Int64ObservableGauge // Bytes in non-idle spans
		ggHeapObjects       metric.Int64ObservableGauge // Total number of allocated objects
		ggHeapReleasedBytes metric.Int64ObservableGauge // Bytes released to the OS
		ggHeapSysBytes      metric.Int64ObservableGauge // Bytes obtained from system for heap
		ggLastGcTimeSeconds metric.Int64ObservableGauge // Time of last garbage collection
		ggLookupsTotal      metric.Int64ObservableGauge // Total number of pointer lookups
		ggMallocsTotal      metric.Int64ObservableGauge // Total count of mallocs
		ggMCacheInuseBytes  metric.Int64ObservableGauge // Bytes in use by mcache structures
		ggMCacheSysBytes    metric.Int64ObservableGauge // Bytes used for mcache structures obtained from system
		ggMspanInuseBytes   metric.Int64ObservableGauge // Bytes in use by mspan structures
		ggMspanSysBytes     metric.Int64ObservableGauge // Bytes used for mspan structures obtained from system
		ggNextGcBytes       metric.Int64ObservableGauge // Size target for next GC cycle
		ggOtherSysBytes     metric.Int64ObservableGauge // Bytes used for other system allocations
		ggStackInuseBytes   metric.Int64ObservableGauge // Bytes in use by stack allocator
		ggGcCompletedCycle  metric.Int64ObservableGauge // Number of completed GC cycles
		ggGcPauseTotal      metric.Int64ObservableGauge // Total pause time of GC in nanoseconds
	}

	// sysGauges implements BasicGauges to collect system-level metrics.
	// It contains observable gauges for OS threads, CGo calls, and goroutines,
	// providing insights into the concurrent behavior and resource utilization
	// of a Go application.
	sysGauges struct {
		// OS and runtime metrics
		ggThreads   metric.Int64ObservableGauge // Number of OS threads created
		ggCgo       metric.Int64ObservableGauge // Number of CGO calls
		ggGRoutines metric.Int64ObservableGauge // Number of goroutines currently active
	}
)
