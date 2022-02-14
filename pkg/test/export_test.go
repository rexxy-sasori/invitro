package test

import (
	"sync"
	"testing"

	mc "github.com/eth-easl/loader/pkg/metric"
	"github.com/stretchr/testify/assert"
)

func TestCheckOverload(t *testing.T) {
	exporter := mc.NewExporter()
	exporter.ReportExecution(
		mc.ExecutionRecord{
			Latency: 50,
			Runtime: 5,
		},
	)
	exporter.ReportExecution(
		mc.ExecutionRecord{
			Latency: 1,
			Runtime: 1,
			Timeout: true,
		},
	)
	exporter.ReportExecution(
		mc.ExecutionRecord{
			Latency: 1,
			Runtime: 1,
			Failed:  true,
		},
	)
	exporter.ReportExecution(
		mc.ExecutionRecord{
			Latency: 1,
			Runtime: 1,
		},
	)
	exporter.ReportExecution(
		mc.ExecutionRecord{
			Latency: 1,
			Runtime: 1,
		},
	)
	assert.True(t, exporter.CheckOverload(0.6))
	assert.False(t, exporter.CheckOverload(0.7))
}

func TestConcurrentReporting(t *testing.T) {
	exporter := mc.NewExporter()
	var wg sync.WaitGroup

	// This function increments a named counter
	// in a loop.
	doReport := func(t, n int) {
		for i := 0; i < n; i++ {
			if t == 0 {
				exporter.ReportExecution(mc.ExecutionRecord{})
			} else {
				exporter.ReportInvocation(mc.MinuteInvocationRecord{})
			}
		}
		wg.Done()
	}

	wg.Add(3)
	go doReport(0, 10_000)
	go doReport(0, 10_000)
	go doReport(1, 10_000)
	wg.Wait()

	assert.Equal(t, 20_000, exporter.GetLantencyRecordLen())
	assert.Equal(t, 10_000, exporter.GetInvocationRecordLen())
}

func TestGetLatenciesInOrder(t *testing.T) {
	exporter := mc.NewExporter()

	exporter.ReportExecution(
		mc.ExecutionRecord{
			Timestamp: 1000_000,
			Latency:   0,
		},
	)
	exporter.ReportExecution(
		mc.ExecutionRecord{
			Timestamp: 1000_000_000,
			Latency:   2,
		},
	)
	exporter.ReportExecution(
		mc.ExecutionRecord{
			Timestamp: 1000,
			Latency:   1,
		},
	)

	assert.EqualValues(t, []float64{1, 0, 2}, exporter.GetLatenciesInOrder())
}

// func TestIsLatencyStationary(t *testing.T) {
// 	exporter := mc.NewExporter()

// 	for i := 0; i < 100; i++ {
// 		exporter.ReportLantency(
// 			mc.LatencyRecord{
// 				Timestamp: int64(i),
// 				Latency:   int64(i + 5000),
// 			},
// 		)
// 	}
// 	assert.False(t, exporter.IsLatencyStationary(0.05))
// 	assert.True(t, exporter.IsLatencyStationary(0.98))
// }