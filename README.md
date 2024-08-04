<h2 align="center"> -=-=-= Overview =-=-=- </h2>

This project sets up a metric monitoring stack using Docker Compose. The stack includes:
- An OpenTelemetry collector
- Prometheus for metrics collection
- Grafana for metrics visualization
- OpenSearch for storing and analyzing metrics
- OpenSearch Dashboard for visualizing metrics
- DataPrepper for processing telemetry data
- OTel JS SDK instrumentation
- OTel Go SDK instrumentation

| No | Pipelines | Description |
| ---- | ---- | ---- |
| 1 | JS SDK → OpenTelemetry Collector → Prometheus → Grafana | This pipeline efficiently ingests metrics telemetry data, specifically focusing on counter-type metrics. It processes and aggregates incoming data streams in real-time, enabling structured monitoring and analysis of system performance indicators. |
| 2 | Go SDK → OpenTelemetry Collector → DataPrepper → OpenSearch | This pipeline specializes in ingesting and processing metrics telemetry data, with a focus on gauges and histograms. It efficiently handles these data types, providing real-time insights into system performance and resource utilization. By integrating gauge measurements and histogram distributions, the pipeline enables comprehensive analysis and visualization of key performance indicators, to make data-driven decisions and optimize the systems effectively. |

<h2 align="center"> -=-=-= Services =-=-=- </h2>

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
- **Ports**: Exposes port **9090**.

### grafana
- **Purpose**: Visualizes metrics collected by Prometheus.
- **Image**: Uses the `grafana/grafana:9.0.1` image.
- **Volumes**: 
  - Configuration file from `./grafana/grafana.ini`.
  - Provisioning files from `./grafana/provisioning/`.
- **Ports**: Exposes port **3000**.

### opensearch
- **Purpose**: Stores and analyzes metrics.
- **Image**: Uses the `opensearchproject/opensearch:latest` image.
- **Environment**:
  - Configures cluster and node settings.
  - Sets JVM heap sizes.
  - Disables security plugins for simplicity.
- **Volumes**: Mounts data volume for persistent storage.
- **Ports**: 
  - **9200**: REST API
  - **9600**: Performance Analyzer

### opensearch-dashboards
- **Purpose**: Visualizes metrics stored in OpenSearch.
- **Image**: Uses the `opensearchproject/opensearch-dashboards:latest` image.
- **Environment**: Configures connection to OpenSearch.
- **Ports**: Exposes port **5601**.

### data-prepper
- **Purpose**: Processes telemetry data and sends it to OpenSearch.
- **Image**: Uses the `opensearchproject/data-prepper:latest` image.
- **Volumes**: Mounts the pipeline configuration file from `./dataprepper/pipelines.yaml`.
- **Ports**: 
  - **21890**, **21891**, **21892**, **4900**
- **Depends on**: `opensearch`, `otel-collector`.

### app-node
- **Purpose**: Runs the REST application along with OTel Javascript SDK instrumentation.
- **Build**: Uses the Dockerfile located in `./app-node`.
- **Ports**: Exposes port 8080.
- **Environment Variables**: Sets the `PORT` to **8080**.

### app-go
- **Purpose**: Runs the Go Server along with OTel Go SDK instrumentation.
- **Build**: Uses the Dockerfile located in `./app-go`.
- **Ports**: Exposes port 8008.
- **Environment Variables**: Sets the `PORT` to **8008**.

## Usage
1. **Clone**: To get repo, run:
   ```sh
   git clone https://github.com/amirrhkm/golang-dockerized-metric-monitoring.git
   ```

2. **Run**: To run the services, navigate to the directory containing `docker-compose.yaml` (either collector-prometheus-pipeline or collector-opensearch-pipeline) and run:
   ```sh
   docker-compose up -d
   ```
   
3. **Test**: To send metric into OpenTelemetry collector endpoint;
   
   - For `app-node` service, navigate to directory `app-node` and run:
   
   Install npm modules and get dependency:
   ```sh
   npm install
   ```
   send API request to http://localhost:8080
   ```sh
   node trigger-node.js
   ```

   - For `app-go` service, send JSON API request to `localhost:8008/update` using Postman which Caddy will redirect into their designated containers.
   ```JSON
   {
      hub_param_a: 1,
      hub_param_b: 1,
      hub_param_c: 1,
   }
   ```

5. **Access Services**:
- REST App: http://localhost:8008 (Go) & http://localhost:8080 (JS)
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000
- OpenSearch: http://localhost:9200
- OpenSearch Dashboards: http://localhost:5601

   #### Heatmap visualization example utilising sum of 3 parameters
   ![Heatmap visualization example](https://github.com/user-attachments/assets/28bf5da0-a38a-4ba9-bee2-83a8e38fb6df)

6. **Shut Down**: To stop and remove the containers, run:
   ```sh
   docker-compose down -v
   ```

<h2 align="center"> -=-=-= Configurations =-=-=- </h2>

### REST Node.js App
- The REST application source code should placed in `./app-node` directory.

### OpenTelemetry Collector
- The OpenTelemetry collector configuration file should be placed in the `./collector/otelcol-config.yaml` file.

### Prometheus
- The Prometheus configuration file should be placed in the `./prometheus/prometheus.yaml` file.

### Grafana
- The Grafana configuration file and provisioning files should be placed in the `./grafana` directory.

###  DataPrepper
- The DataPrepper pipeline configuration file should be placed in the `./dataprepper/pipelines.yaml` file.

## Simple Architecture Diagram

### 1. Collector-Prometheus Pipeline
![Collector-Prometheus-Pipeline](https://github.com/user-attachments/assets/0e82a200-94a7-4417-bb30-b4d9f167727c)

### 2. Collector-OpenSearch Pipeline
![OpenSearch-Pipeline drawio](https://github.com/user-attachments/assets/6801194d-3297-4f83-8fe4-a0bc44a173bf)




