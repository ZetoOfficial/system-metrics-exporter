package cpu

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/cpu"
)

type Metrics struct {
	Usage prometheus.Gauge
}

func NewCPUMetrics() *Metrics {
	usage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "exporter_cpu_usage_percent",
		Help: "Процент использования CPU",
	})

	prometheus.MustRegister(usage)

	return &Metrics{
		Usage: usage,
	}
}

func (m *Metrics) Collect() {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Ошибка при сборе CPU: %v", err)
		return
	}
	if len(percent) > 0 {
		m.Usage.Set(percent[0])
	}
}
