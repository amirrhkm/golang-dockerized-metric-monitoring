package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

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

	mu sync.Mutex
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
		mu.Lock()
		defer mu.Unlock()

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

func main() {
	ctx := context.Background()

	// Resources
	aResource, err := newResource("A")
	if err != nil {
		log.Fatalf("failed to initialize A resource: %v", err)
	}
	bResource, err := newResource("B")
	if err != nil {
		log.Fatalf("failed to initialize B resource: %v", err)
	}
	cResource, err := newResource("C")
	if err != nil {
		log.Fatalf("failed to initialize C resource: %v", err)
	}

	// meterProviders
	aMeterProvider, err := initMeterProvider(ctx, aResource)
	if err != nil {
		log.Fatalf("failed to initialize A meter provider: %v", err)
	}
	defer aMeterProvider.Shutdown(ctx)

	bMeterProvider, err := initMeterProvider(ctx, bResource)
	if err != nil {
		log.Fatalf("failed to initialize B meter provider: %v", err)
	}
	defer bMeterProvider.Shutdown(ctx)

	cMeterProvider, err := initMeterProvider(ctx, cResource)
	if err != nil {
		log.Fatalf("failed to initialize C meter provider: %v", err)
	}
	defer cMeterProvider.Shutdown(ctx)

	// Gauges
	err = initGauges(aMeterProvider, "A")
	if err != nil {
		log.Fatalf("failed to initialize A metrics: %v", err)
	}

	err = initGauges(bMeterProvider, "B")
	if err != nil {
		log.Fatalf("failed to initialize B metrics: %v", err)
	}

	err = initGauges(cMeterProvider, "C")
	if err != nil {
		log.Fatalf("failed to initialize C metrics: %v", err)
	}

	handleRequest := func(pattern string, _ *metric.Int64ObservableGauge, value *int64) {
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

			mu.Lock()
			*value = val
			mu.Unlock()

			log.Printf("Updated %s to %d", parts[1], val)
			fmt.Fprintf(w, "Updated %s to %d", parts[1], val)
		})
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		handleRequest("/a/", &aGauge, &aValue)
	}()
	go func() {
		defer wg.Done()
		handleRequest("/b/", &bGauge, &bValue)
	}()
	go func() {
		defer wg.Done()
		handleRequest("/c/", &cGauge, &cValue)
	}()

	log.Println("Server is up and running on port 8008")
	log.Fatal(http.ListenAndServe(":8008", nil))

	wg.Wait()
}
