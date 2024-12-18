# Exporter на Go

## Описание

Приложение-экспортер на Go, которое собирает информацию о:

- Использовании CPU
- Объёме оперативной памяти
- Использовании дискового пространства

Метрики предоставляются в формате Prometheus и доступны на корневом пути `/metrics`.

## Структура проекта

```
exporter/
├── cmd/
│   └── exporter/
│       └── main.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── cpu/
│   │   └── cpu.go
│   ├── memory/
│   │   └── memory.go
│   └── disk/
│       └── disk.go
├── go.mod
└── README.md
```

## Переменные окружения

- `EXPORTER_HOST`: Хост для запуска экспортера (по умолчанию `0.0.0.0`)
- `EXPORTER_PORT`: Порт для запуска экспортера (по умолчанию `8080`)
- `EXPORTER_METRICS_INTERVAL`: Интервал сбора метрик в секундах (по умолчанию `10`)

## Метрики

### CPU

- `exporter_cpu_usage_percent`: Процент использования CPU

### Память

- `exporter_memory_total_bytes`: Общий объём оперативной памяти в байтах
- `exporter_memory_used_bytes`: Используемый объём оперативной памяти в байтах

### Диски

- `exporter_disk_total_bytes{mountpoint="<path>"}`
- `exporter_disk_used_bytes{mountpoint="<path>"}`

## Запросы PromQL

### Использование CPU

```promql
exporter_cpu_usage_percent
```

### Памяти всего и используется

**Общий объём памяти:**

```promql
exporter_memory_total_bytes
```

**Используемый объём памяти:**

```promql
exporter_memory_used_bytes
```

### Объём дисков всего и используется

**Общий объём дисков:**

```promql
exporter_disk_total_bytes
```

**Используемый объём дисков:**

```promql
exporter_disk_used_bytes
```

**Пример с фильтрацией по точке монтирования:**

```promql
exporter_disk_used_bytes{mountpoint="/"}
```

## Запуск

### Локально

1. **Установите зависимости:**

   ```bash
   go mod tidy
   ```

2. **Соберите приложение:**

   ```bash
   go build -o exporter ./cmd/exporter
   ```

3. **Запустите экспортера:**

   ```bash
   EXPORTER_HOST=0.0.0.0 EXPORTER_PORT=8080 EXPORTER_METRICS_INTERVAL=10 ./exporter
   ```

### Использование с Docker

Для запуска Графаны + Прометеуса можно просто поднять контейнер:

```bash
docker-compose up --build
```

## Интеграция с Prometheus

Отредактируйте prometheus.yml

```yaml
scrape_configs:
  - job_name: "go_exporter"
    static_configs:
      - targets: ["<EXPORTER_HOST>:<EXPORTER_PORT>"]
```

Замените `<EXPORTER_HOST>` и `<EXPORTER_PORT>` на соответствующие значения.
