/**
 * Initializes the OpenTelemetry MeterProvider and sets up metric exporting.
 * @module meter
 */

"use strict";

const { Resource } = require("@opentelemetry/resources");
const { metrics } = require("@opentelemetry/api");
const {
  OTLPMetricExporter,
} = require("@opentelemetry/exporter-metrics-otlp-grpc");
const {
  MeterProvider,
  PeriodicExportingMetricReader,
} = require("@opentelemetry/sdk-metrics");

/**
 * Username for HTTP basic authentication.
 * @type {string}
 */
const username = "admin";

/**
 * Password for HTTP basic authentication.
 * @type {string}
 */
const password = "admin";

/**
 * HTTP Authorization header value for basic authentication.
 * @type {string}
 */
const authHeader = "Basic " + Buffer.from(`${username}:${password}`).toString("base64");

/**
 * Creates a new OpenTelemetry MeterProvider.
 * @type {MeterProvider}
 */
const meterProvider = new MeterProvider({
  resource: new Resource({ "service.name": "my-express-app" }),
});

/**
 * Initializes an OTLP Metric Exporter with basic authentication headers.
 * @type {OTLPMetricExporter}
 */
const metricExporter = new OTLPMetricExporter({
  url: "http://0.0.0.0:4317",
  headers: {
    "Authorization": authHeader
  }
});

/**
 * Logs errors encountered during metric exporting.
 */
metricExporter.on("error", (error) => {
  console.error("Error exporting metrics:", error);
});

/**
 * Creates a PeriodicExportingMetricReader for exporting metrics periodically.
 * @type {PeriodicExportingMetricReader}
 */
const metricReader = new PeriodicExportingMetricReader({
  exporter: metricExporter,
  exportIntervalMillis: 60000,
});

/**
 * Adds the metric reader to the global meter provider.
 */
meterProvider.addMetricReader(metricReader);

/**
 * Sets the global meter provider for the OpenTelemetry API.
 */
metrics.setGlobalMeterProvider(meterProvider);
