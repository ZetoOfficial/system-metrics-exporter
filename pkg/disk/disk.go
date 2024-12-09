package disk

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/disk"
)

type Metrics struct {
	Total *prometheus.GaugeVec
	Used  *prometheus.GaugeVec
}

func NewDiskMetrics() *Metrics {
	total := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "exporter_disk_total_bytes",
		Help: "Общий объем диска в байтах",
	}, []string{"mountpoint"})

	used := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "exporter_disk_used_bytes",
		Help: "Используемый объем диска в байтах",
	}, []string{"mountpoint"})

	prometheus.MustRegister(total)
	prometheus.MustRegister(used)

	return &Metrics{
		Total: total,
		Used:  used,
	}
}

func (d *Metrics) Collect() {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Ошибка при получении разделов диска: %v", err)
		return
	}

	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			log.Printf("Ошибка при получении использования диска для %s: %v", p.Mountpoint, err)
			continue
		}
		d.Total.WithLabelValues(p.Mountpoint).Set(float64(usage.Total))
		d.Used.WithLabelValues(p.Mountpoint).Set(float64(usage.Used))
	}
}
