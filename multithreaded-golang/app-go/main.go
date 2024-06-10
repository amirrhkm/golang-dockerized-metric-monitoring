package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"google.golang.org/grpc/credentials/insecure"
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

func initMetrics(meterProvider *sdkmetric.MeterProvider, serviceName string) (metric.Int64Counter, error) {
	meter := meterProvider.Meter(serviceName + "-meter")

	counter, err := meter.Int64Counter(
		serviceName+"_api_call",
		metric.WithDescription("Number of API calls to "+serviceName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create counter for %s: %w", serviceName, err)
	}

	return counter, nil
}

func main() {
	ctx := context.Background()

	// Resources
	posResource, err := newResource("POS")
	if err != nil {
		log.Fatalf("failed to initialize POS resource: %v", err)
	}
	hubResource, err := newResource("HUB")
	if err != nil {
		log.Fatalf("failed to initialize HUB resource: %v", err)
	}
	cdsResource, err := newResource("CDS")
	if err != nil {
		log.Fatalf("failed to initialize CLOUD resource: %v", err)
	}

	// meterProviders
	posMeterProvider, err := initMeterProvider(ctx, posResource)
	if err != nil {
		log.Fatalf("failed to initialize POS/CDS meter provider: %v", err)
	}
	defer posMeterProvider.Shutdown(ctx)

	hubMeterProvider, err := initMeterProvider(ctx, hubResource)
	if err != nil {
		log.Fatalf("failed to initialize HUB meter provider: %v", err)
	}
	defer hubMeterProvider.Shutdown(ctx)

	cdsMeterProvider, err := initMeterProvider(ctx, cdsResource)
	if err != nil {
		log.Fatalf("failed to initialize CLOUD meter provider: %v", err)
	}
	defer cdsMeterProvider.Shutdown(ctx)

	// Counters
	posCounter, err := initMetrics(posMeterProvider, "POS")
	if err != nil {
		log.Fatalf("failed to initialize POS metrics: %v", err)
	}

	hubCounter, err := initMetrics(hubMeterProvider, "HUB")
	if err != nil {
		log.Fatalf("failed to initialize HUB metrics: %v", err)
	}

	cdsCounter, err := initMetrics(cdsMeterProvider, "CDS")
	if err != nil {
		log.Fatalf("failed to initialize CLOUD metrics: %v", err)
	}

	var wg sync.WaitGroup

	handleRequest := func(pattern string, counter metric.Int64Counter, serviceName string) {
		http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			name := r.URL.Path[len(pattern):]
			data := []attribute.KeyValue{
				attribute.String("route", pattern+":name"),
				attribute.String("name", name),
			}

			counter.Add(r.Context(), 1, metric.WithAttributes(data...))
			log.Printf("Received request for %s service: %s", serviceName, name)

			fmt.Fprintf(w, "Hello from %s service, %s", serviceName, name)
		})
	}

	wg.Add(3)
	go func() {
		defer wg.Done()
		handleRequest("/pos/", posCounter, "POS")
	}()
	go func() {
		defer wg.Done()
		handleRequest("/hub/", hubCounter, "HUB")
	}()
	go func() {
		defer wg.Done()
		handleRequest("/cds/", cdsCounter, "CDS")
	}()

	wg.Wait()

	log.Println("Server is up and running on port 8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}
