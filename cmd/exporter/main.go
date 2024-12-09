package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ZetoOfficial/exporter/pkg/config"
	"github.com/ZetoOfficial/exporter/pkg/cpu"
	"github.com/ZetoOfficial/exporter/pkg/disk"
	"github.com/ZetoOfficial/exporter/pkg/memory"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 120 * time.Second
	shutdownTimeout = 5 * time.Second
)

func main() {
	cfg := config.LoadConfig()
	config.ValidateConfig(cfg)

	cpuMetrics := cpu.NewCPUMetrics()
	memoryMetrics := memory.NewMemoryMetrics()
	diskMetrics := disk.NewDiskMetrics()

	metricsInterval := time.Duration(cfg.MetricsInterval) * time.Second
	ticker := time.NewTicker(metricsInterval)
	defer ticker.Stop()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cpuMetrics.Collect()
				memoryMetrics.Collect()
				diskMetrics.Collect()
			case <-done:
				return
			}
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	go func() {
		log.Printf("Экспортер запущен на %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	log.Println("Получен сигнал остановки, завершаем работу...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при завершении работы сервера: %v", err)
	}

	log.Println("Экспортер корректно завершил работу")
}
