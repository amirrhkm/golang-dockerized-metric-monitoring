package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	aGauge metric.Int64ObservableGauge
	bGauge metric.Int64ObservableGauge
	cGauge metric.Int64ObservableGauge

	aValue int64
	bValue int64
	cValue int64
)

func newResource(serviceName string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("0.1.0"),
		))
}

func initMeterProvider(ctx context.Context, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:admin"))

	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("0.0.0.0:4317"),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()),
		otlpmetricgrpc.WithHeaders(map[string]string{"Authorization": authHeader}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)

	return meterProvider, nil
}

func initGauges(meterProvider *sdkmetric.MeterProvider, serviceName string) error {
	meter := meterProvider.Meter(serviceName + "-meter")

	var err error
	aGauge, err = meter.Int64ObservableGauge(
		serviceName+"_param_a",
		metric.WithDescription("Tracks the status of param_a (0 or 1)"),
	)
	if err != nil {
		return fmt.Errorf("failed to create gauge for param_a: %w", err)
	}

	bGauge, err = meter.Int64ObservableGauge(
		serviceName+"_param_b",
		metric.WithDescription("Tracks the status of param_b (0 or 1)"),
	)
	if err != nil {
		return fmt.Errorf("failed to create gauge for param_b: %w", err)
	}

	cGauge, err = meter.Int64ObservableGauge(
		serviceName+"_param_c",
		metric.WithDescription("Tracks the status of param_c (0 or 1)"),
	)
	if err != nil {
		return fmt.Errorf("failed to create gauge for param_c: %w", err)
	}

	_, err = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		observer.ObserveInt64(aGauge, aValue)
		observer.ObserveInt64(bGauge, bValue)
		observer.ObserveInt64(cGauge, cValue)
		return nil
	}, aGauge, bGauge, cGauge)

	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	return nil
}

func handleRequest(pattern string, value *int64) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		val, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil || (val != 0 && val != 1) {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}

		*value = val

		log.Printf("Updated %s to %d", parts[1], val)
		fmt.Fprintf(w, "Updated %s to %d", parts[1], val)
	})
}

func main() {
	ctx := context.Background()

	// Resources
	resource, err := newResource("hub")
	if err != nil {
		log.Fatalf("failed to initialize hub resource: %v", err)
	}

	// MeterProvider
	meterProvider, err := initMeterProvider(ctx, resource)
	if err != nil {
		log.Fatalf("failed to initialize hub meter provider: %v", err)
	}
	defer meterProvider.Shutdown(ctx)

	// Gauges
	err = initGauges(meterProvider, "hub")
	if err != nil {
		log.Fatalf("failed to initialize hub metrics: %v", err)
	}

	// Handle requests
	handleRequest("/a/", &aValue)
	handleRequest("/b/", &bValue)
	handleRequest("/c/", &cValue)

	log.Println("Server is up and running on port 8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}
