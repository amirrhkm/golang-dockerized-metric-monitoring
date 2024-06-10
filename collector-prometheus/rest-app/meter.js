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

const username = "admin";
const password = "admin";

const authHeader = "Basic " + Buffer.from(`${username}:${password}`).toString("base64");

const meterProvider = new MeterProvider({
  resource: new Resource({ "service.name": "my-express-app" }),
});

const metricExporter = new OTLPMetricExporter({
  url: "http://otel-collector:4317",
  headers: {
    "Authorization": authHeader
  }
});

const metricReader = new PeriodicExportingMetricReader({
  exporter: metricExporter,
  exportIntervalMillis: 60000,
});

meterProvider.addMetricReader(metricReader);

metrics.setGlobalMeterProvider(meterProvider);
