//Filename: /pkg/metrics.go
package pkg

import (
	"sync"
	"time"
)

type Metrics struct {
	mu             sync.Mutex
	MessageCount   int
	TotalLatency   time.Duration
	MaxLatency     time.Duration
	MinLatency     time.Duration
	DroppedPackets int
}

func NewMetrics() *Metrics {
	return &Metrics{
		MinLatency: time.Hour, // Initialize with large value
	}
}

func (m *Metrics) Record(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.MessageCount++
	m.TotalLatency += latency

	if latency > m.MaxLatency {
		m.MaxLatency = latency
	}
	if latency < m.MinLatency {
		m.MinLatency = latency
	}
}

func (m *Metrics) RecordDrop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DroppedPackets++
}

func (m *Metrics) AverageLatency() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.MessageCount == 0 {
		return 0
	}
	return m.TotalLatency / time.Duration(m.MessageCount)
}