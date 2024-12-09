package memory

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct {
	Total prometheus.Gauge
	Used  prometheus.Gauge
}

func NewMemoryMetrics() *Metrics {
	total := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "exporter_memory_total_bytes",
		Help: "Общий объем оперативной памяти в байтах",
	})

	used := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "exporter_memory_used_bytes",
		Help: "Используемый объем оперативной памяти в байтах",
	})

	prometheus.MustRegister(total)
	prometheus.MustRegister(used)

	return &Metrics{
		Total: total,
		Used:  used,
	}
}

func (m *Metrics) Collect() {
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Ошибка при сборе памяти: %v", err)
		return
	}
	m.Total.Set(float64(vm.Total))
	m.Used.Set(float64(vm.Used))
}
