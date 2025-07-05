package handler

import (
	"sync"
)

type Metrics struct {
	TotalRequests        int
	RequestsPerProvider  map[string]int
	TotalResponseTimeMs  int64
	ResponseTimePerProvider map[string]int64
	mutex               sync.Mutex
}

var metrics = &Metrics{
	RequestsPerProvider:     make(map[string]int),
	ResponseTimePerProvider: make(map[string]int64),
}

func (m *Metrics) Add(provider string, durationMs int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.TotalRequests++
	m.RequestsPerProvider[provider]++
	m.TotalResponseTimeMs += durationMs
	m.ResponseTimePerProvider[provider] += durationMs
}

func (m *Metrics) Snapshot() (total int, perProvider map[string]int, avgMs float64, avgPerProvider map[string]float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	total = m.TotalRequests
	perProvider = make(map[string]int)
	avgPerProvider = make(map[string]float64)
	for k, v := range m.RequestsPerProvider {
		perProvider[k] = v
		if v > 0 {
			avgPerProvider[k] = float64(m.ResponseTimePerProvider[k]) / float64(v)
		}
	}
	if total > 0 {
		avgMs = float64(m.TotalResponseTimeMs) / float64(total)
	}
	return
}
