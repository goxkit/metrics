// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package system provides detailed memory metrics collection functionality.
package system

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
)

// NewMemGauges creates a new memory metrics collector that monitors various aspects
// of the Go runtime memory usage and garbage collection. It initializes all the
// necessary observable gauges for tracking memory allocation, utilization,
// garbage collection statistics, and other related metrics.
//
// Parameters:
//   - meter: The OpenTelemetry meter used to create gauge instruments.
//
// Returns:
//   - A BasicGauges implementation for memory metrics collection.
//   - An error if any gauge creation fails.
func NewMemGauges(meter metric.Meter) (BasicGauges, error) {
	ggSysBytes, err := meter.Int64ObservableGauge("go_memstats_sys_bytes", metric.WithDescription("Number of bytes obtained from system."))
	if err != nil {
		return nil, err
	}

	ggAllocBytesTotal, err := meter.Int64ObservableGauge("go_memstats_alloc_bytes_total", metric.WithDescription("Total number of bytes allocated, even if freed."))
	if err != nil {
		return nil, err
	}

	ggHeapAllocBytes, err := meter.Int64ObservableGauge("go_memstats_heap_alloc_bytes", metric.WithDescription("Number of heap bytes allocated and still in use."))
	if err != nil {
		return nil, err
	}

	ggFreesTotal, err := meter.Int64ObservableGauge("go_memstats_frees_total", metric.WithDescription("Total number of frees."))
	if err != nil {
		return nil, err
	}

	ggGcSysBytes, err := meter.Int64ObservableGauge("go_memstats_gc_sys_bytes", metric.WithDescription("Number of bytes used for garbage collection system metadata."))
	if err != nil {
		return nil, err
	}

	ggHeapIdleBytes, err := meter.Int64ObservableGauge("go_memstats_heap_idle_bytes", metric.WithDescription("Number of heap bytes waiting to be used."))
	if err != nil {
		return nil, err
	}

	ggInuseBytes, err := meter.Int64ObservableGauge("go_memstats_heap_inuse_bytes", metric.WithDescription("Number of heap bytes that are in use."))
	if err != nil {
		return nil, err
	}

	ggHeapObjects, err := meter.Int64ObservableGauge("go_memstats_heap_objects", metric.WithDescription("Number of allocated objects."))
	if err != nil {
		return nil, err
	}

	ggHeapReleasedBytes, err := meter.Int64ObservableGauge("go_memstats_heap_released_bytes", metric.WithDescription("Number of heap bytes released to OS."))
	if err != nil {
		return nil, err
	}

	ggHeapSysBytes, err := meter.Int64ObservableGauge("go_memstats_heap_sys_bytes", metric.WithDescription("Number of heap bytes obtained from system."))
	if err != nil {
		return nil, err
	}

	ggLastGcTimeSeconds, err := meter.Int64ObservableGauge("go_memstats_last_gc_time_seconds", metric.WithDescription("Number of seconds since 1970 of last garbage collection."))
	if err != nil {
		return nil, err
	}

	ggLookupsTotal, err := meter.Int64ObservableGauge("go_memstats_lookups_total", metric.WithDescription("Total number of pointer lookups."))
	if err != nil {
		return nil, err
	}

	ggMallocsTotal, err := meter.Int64ObservableGauge("go_memstats_mallocs_total", metric.WithDescription("Total number of mallocs."))
	if err != nil {
		return nil, err
	}

	ggMCacheInuseBytes, err := meter.Int64ObservableGauge("go_memstats_mcache_inuse_bytes", metric.WithDescription("Number of bytes in use by mcache structures."))
	if err != nil {
		return nil, err
	}

	ggMCacheSysBytes, err := meter.Int64ObservableGauge("go_memstats_mcache_sys_bytes", metric.WithDescription("Number of bytes used for mcache structures obtained from system."))
	if err != nil {
		return nil, err
	}

	ggMspanInuseBytes, err := meter.Int64ObservableGauge("go_memstats_mspan_inuse_bytes", metric.WithDescription("Number of bytes in use by mspan structures."))
	if err != nil {
		return nil, err
	}

	ggMspanSysBytes, err := meter.Int64ObservableGauge("go_memstats_mspan_sys_bytes", metric.WithDescription("Number of bytes used for mspan structures obtained from system."))
	if err != nil {
		return nil, err
	}

	ggNextGcBytes, err := meter.Int64ObservableGauge("go_memstats_next_gc_bytes", metric.WithDescription("Number of heap bytes when next garbage collection will take place."))
	if err != nil {
		return nil, err
	}

	ggOtherSysBytes, err := meter.Int64ObservableGauge("go_memstats_other_sys_bytes", metric.WithDescription("Number of bytes used for other system allocations."))
	if err != nil {
		return nil, err
	}

	ggStackInuseBytes, err := meter.Int64ObservableGauge("go_memstats_stack_inuse_bytes", metric.WithDescription("Number of bytes in use by the stack allocator."))
	if err != nil {
		return nil, err
	}

	ggGcCompletedCycle, err := meter.Int64ObservableGauge("go_memstats_gc_completed_cycle", metric.WithDescription("Number of GC cycle completed."))
	if err != nil {
		return nil, err
	}

	ggGcPauseTotal, err := meter.Int64ObservableGauge("go_memstats_gc_pause_total", metric.WithDescription("Number of GC-stop-the-world caused in Nanosecond."))
	if err != nil {
		return nil, err
	}

	return &memGauges{
		ggSysBytes,
		ggAllocBytesTotal,
		ggHeapAllocBytes,
		ggFreesTotal,
		ggGcSysBytes,
		ggHeapIdleBytes,
		ggInuseBytes,
		ggHeapObjects,
		ggHeapReleasedBytes,
		ggHeapSysBytes,
		ggLastGcTimeSeconds,
		ggLookupsTotal,
		ggMallocsTotal,
		ggMCacheInuseBytes,
		ggMCacheSysBytes,
		ggMspanInuseBytes,
		ggMspanSysBytes,
		ggNextGcBytes,
		ggOtherSysBytes,
		ggStackInuseBytes,
		ggGcCompletedCycle,
		ggGcPauseTotal,
	}, nil
}

// Collect registers callbacks for memory metrics collection.
// It reads memory statistics from the Go runtime and reports them through the
// observable gauges. The callback function will be invoked periodically by the
// OpenTelemetry SDK to gather the latest memory statistics.
//
// Parameters:
//   - meter: The OpenTelemetry meter used to register callbacks.
func (m *memGauges) Collect(meter metric.Meter) {
	// Define a callback function that will be called periodically to collect metrics
	cb := func(_ context.Context, observer metric.Observer) error {
		// Retrieve the current memory statistics from the Go runtime
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		// Record all memory metrics using the observer
		observer.ObserveInt64(m.ggSysBytes, int64(stats.Sys))                   // Total memory obtained from OS
		observer.ObserveInt64(m.ggAllocBytesTotal, int64(stats.TotalAlloc))     // Total bytes allocated (even if freed)
		observer.ObserveInt64(m.ggHeapAllocBytes, int64(stats.HeapAlloc))       // Bytes allocated and in use
		observer.ObserveInt64(m.ggFreesTotal, int64(stats.Frees))               // Total number of frees
		observer.ObserveInt64(m.ggGcSysBytes, int64(stats.GCSys))               // Memory used for GC metadata
		observer.ObserveInt64(m.ggHeapIdleBytes, int64(stats.HeapIdle))         // Heap memory waiting to be used
		observer.ObserveInt64(m.ggInuseBytes, int64(stats.HeapInuse))           // Heap memory in use
		observer.ObserveInt64(m.ggHeapObjects, int64(stats.HeapObjects))        // Number of allocated objects
		observer.ObserveInt64(m.ggHeapReleasedBytes, int64(stats.HeapReleased)) // Heap memory returned to OS
		observer.ObserveInt64(m.ggHeapSysBytes, int64(stats.HeapSys))           // Heap memory obtained from OS
		observer.ObserveInt64(m.ggLastGcTimeSeconds, int64(stats.LastGC))       // Time of last GC
		observer.ObserveInt64(m.ggLookupsTotal, int64(stats.Lookups))           // Number of pointer lookups
		observer.ObserveInt64(m.ggMallocsTotal, int64(stats.Mallocs))           // Total number of mallocs
		observer.ObserveInt64(m.ggMCacheInuseBytes, int64(stats.MCacheInuse))   // Bytes in mcache structures
		observer.ObserveInt64(m.ggMCacheSysBytes, int64(stats.MCacheSys))       // MCacheSys bytes from system
		observer.ObserveInt64(m.ggMspanInuseBytes, int64(stats.MSpanInuse))     // Bytes in mspan structures
		observer.ObserveInt64(m.ggMspanSysBytes, int64(stats.MSpanSys))         // MSpanSys bytes from system
		observer.ObserveInt64(m.ggNextGcBytes, int64(stats.NextGC))             // Target heap size of next GC
		observer.ObserveInt64(m.ggOtherSysBytes, int64(stats.OtherSys))         // Other system allocations
		observer.ObserveInt64(m.ggStackInuseBytes, int64(stats.StackSys))       // Stack system bytes
		observer.ObserveInt64(m.ggGcCompletedCycle, int64(stats.NumGC))         // Number of completed GC cycles
		observer.ObserveInt64(m.ggGcPauseTotal, int64(stats.PauseTotalNs))      // Total GC pause time in nanoseconds

		return nil
	}

	// Register the callback with the meter
	// We ignore the returned registration to avoid verbosity
	_, _ = meter.RegisterCallback(cb)
}
