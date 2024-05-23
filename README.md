# Docker Compose Setup for Metric Monitoring

## Overview

This project sets up a metric monitoring stack using Docker Compose. The stack includes:
- An OpenTelemetry collector
- Prometheus for metrics collection
- Grafana for metrics visualization
- OTel Node.js instrumentation
- OTel Go instrumentation

## Services

### otel-collector
- **Purpose**: Collects telemetry data.
- **Image**: Uses the `otel/opentelemetry-collector-contrib:latest` image.
- **Command**: Uses `otelcol-config.yaml` for configuration.
- **Volumes**: Mounts the configuration file from `./collector/otelcol-config.yaml`.
- **Ports**: Exposes various ports for different extensions and receivers.
   - **1888**: pprof extension
   - **8888**: Prometheus metrics exposed by the collector
   - **8889**: Prometheus exporter metrics
   - **13133**: health_check extension
   - **4317**: OTLP gRPC receiver
   - **4318**: OTLP HTTP receiver
   - **55679**: zpages extension

### prometheus
- **Purpose**: Scrapes and stores metrics.
- **Image**: Uses the `quay.io/prometheus/prometheus:v2.34.0` image.
- **Command**: Configured via `prometheus.yaml`.
- **Volumes**: Mounts the configuration file from `./prometheus/prometheus.yaml`.
- **Ports**: Exposes port 9090.

### grafana
- **Purpose**: Visualizes metrics collected by Prometheus.
- **Image**: Uses the `grafana/grafana:9.0.1` image.
- **Volumes**: 
  - Configuration file from `./grafana/grafana.ini`.
  - Provisioning files from `./grafana/provisioning/`.
- **Ports**: Exposes port 3000.

### app-node
- **Purpose**: Runs the REST application along with OTel Nodejs instrumentation.
- **Build**: Uses the Dockerfile located in `./app-node`.
- **Ports**: Exposes port 8008.
- **Environment Variables**: Sets the `PORT` to 8008.

## Usage

1. **Run**: To run the services, navigate to the directory containing `docker-compose.yaml` and run:
   ```sh
   docker-compose up -d
   ```

2. **Access Services**:
- REST App: http://localhost:8008 (Go) & http://localhost:8080 (Node.js)
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000

3. **Shut Down**: To stop and remove the containers, run:
   ```sh
   docker-compose down
   ```

## Configuration

### REST Node.js App
- The REST application source code should placed in `./app-node` directory.

### OpenTelemetry Collector
- The OpenTelemetry collector configuration file should be placed in the `./collector/otelcol-config.yaml` file.

### Prometheus
- The Prometheus configuration file should be placed in the `./prometheus/prometheus.yaml` file.

### Grafana
- The Grafana configuration file and provisioning files should be placed in the `./grafana` directory.

## Simple Architecture Diagram
![metric-monitoring-diagram](https://github.com/amirrhkm/metrics-monitoring/assets/152793780/dd1ce8c1-2b41-4c04-b505-324b9f7c8d9b)



