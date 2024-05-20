# Dockerized Data Collection with OpenTelemetry, Prometheus, and Grafana

      This repository contains a Dockerized setup for collecting, storing, and visualizing telemetry data using OpenTelemetry, Prometheus, and Grafana.

## Components

### OpenTelemetry

OpenTelemetry provides a set of APIs, libraries, agents, and instrumentation to enable observability in distributed systems. In this setup, OpenTelemetry is responsible for instrumenting applications to collect telemetry data such as metrics, traces, and logs.

#### - Language-Agnostic and Platform-Agnostic Nature

      OpenTelemetry uses Protocol Buffers (protobuf) for serializing and transmitting telemetry data. Protocol Buffers are language-agnostic and platform-agnostic, meaning they can define data structures in a format independent of any specific programming language or computing platform. This allows for seamless integration and interoperability between different parts of a system, regardless of the languages or platforms they are built upon.

### Prometheus

Prometheus is a monitoring and alerting toolkit designed for reliability and scalability. It collects metrics data from instrumented targets using a pull model, stores it efficiently, and provides a powerful query language for analyzing and visualizing the data.

### Grafana

Grafana is an open-source platform for monitoring and observability. It provides a rich set of visualization options and dashboards for exploring and understanding metrics, logs, and traces data.

### Nginx

Nginx is a high-performance web server and reverse proxy that can also function as a load balancer, mail proxy, and HTTP cache. In this setup, Nginx acts as a reverse proxy to authenticate requests before forwarding them to the OpenTelemetry collector receiver.

## Setup

To run the entire data collection stack locally using Docker, follow these steps:

1. Clone this repository:

   ```bash
   git clone https://github.com/amirrhkm/metric-monitoring.git

2. Navigate to the cloned directory:

   ```bash
   cd your-project-directory

3. Start the Docker containers:

   ```bash
   docker-compose up -d
   ```
   
   This command will launch Docker containers for OpenTelemetry, Prometheus, and Grafana in detached mode.

4. Access Grafana in your web browser:

   Grafana should now be running and accessible at: http://localhost:3000

![otel drawio](https://github.com/amirrhkm/metrics-monitoring/assets/152793780/cd7dd20e-fc82-4927-b389-c398e5051ddc)
