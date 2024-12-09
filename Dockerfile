FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/exporter ./cmd/exporter
COPY ./pkg ./pkg

RUN go build -o exporter ./cmd/exporter

EXPOSE 8080

ENV EXPORTER_HOST=0.0.0.0
ENV EXPORTER_PORT=8080
ENV EXPORTER_METRICS_INTERVAL=10

CMD ["./exporter"]