# Docker Compose Setup for Metric Monitoring

## Overview

This project sets up a metric monitoring stack using Docker Compose. The stack includes:
- An OpenTelemetry collector
- Prometheus for metrics collection
- Grafana for metrics visualization
- A REST application service

## Services

### otel-collector
- **Purpose**: Collects telemetry data.
- **Image**: Uses the `otel/opentelemetry-collector-contrib:latest` image.
- **Command**: Uses `otel-collector-config.yaml` for configuration.
- **Volumes**: Mounts the configuration file from `./collector/otel-collector-config.yaml`.
- **Ports**: Exposes various ports for different extensions and receivers.

### prometheus
- **Purpose**: Scrapes and stores metrics.
- **Image**: Uses the `quay.io/prometheus/prometheus:v2.34.0` image.
- **Command**: Configured via `prometheus.yml`.
- **Volumes**: Mounts the configuration file from `./prometheus/prometheus.yaml`.
- **Ports**: Exposes port 9090.

### grafana
- **Purpose**: Visualizes metrics collected by Prometheus.
- **Image**: Uses the `grafana/grafana:9.0.1` image.
- **Volumes**: 
  - Configuration file from `./grafana/grafana.ini`.
  - Provisioning files from `./grafana/provisioning/`.
- **Ports**: Exposes port 3000.

### rest-app
- **Purpose**: Runs the REST application.
- **Build**: Uses the Dockerfile located in `./rest-app`.
- **Ports**: Exposes port 8008.
- **Environment Variables**: Sets the `PORT` to 8008.

## Networks

### net
- **Driver**: Uses the bridge driver to create an isolated network for the services.

## Usage

1. **Run**: To run the services, navigate to the directory containing `docker-compose.yaml` and run:
   ```sh
   docker-compose up -d
   ```

2. **Access Services**:
- REST App: http://localhost:8008
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000

3. **Shut Down**: To stop and remove the containers, run:
   ```sh
   docker-compose down
   ```

## Configuration

### REST App
- The REST application source code should placed in `./rest-app` directory.

### OpenTelemetry Collector
- The OpenTelemetry collector configuration file should be placed in the `./collector/otel-collector-config.yaml` file.

### Prometheus
- The Prometheus configuration file should be placed in the `./prometheus/prometheus.yaml` file.

### Grafana
- The Grafana configuration file and provisioning files should be placed in the `./grafana` directory.

## Simple Architecture Diagram
![metric-monitoring-diagram](https://github.com/amirrhkm/metrics-monitoring/assets/152793780/233d01dd-ec18-4697-b64b-e6d2f692c57c)
