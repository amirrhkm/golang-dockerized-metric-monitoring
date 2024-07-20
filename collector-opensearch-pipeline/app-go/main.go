package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	aValue     int64
	bValue     int64
	cValue     int64
	instanceID string
	mu         sync.Mutex
)

type HubParams struct {
	ParamA *int64 `json:"hub_param_a,omitempty"`
	ParamB *int64 `json:"hub_param_b,omitempty"`
	ParamC *int64 `json:"hub_param_c,omitempty"`
}

func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(instanceID),
		))
}

func newManualReader() *sdkmetric.ManualReader {
	return sdkmetric.NewManualReader(
		sdkmetric.WithTemporalitySelector(func(kind sdkmetric.InstrumentKind) metricdata.Temporality {
			return metricdata.DeltaTemporality
		}),
	)
}

func initMeterProvider(res *resource.Resource, reader *sdkmetric.ManualReader) (*sdkmetric.MeterProvider, error) {
	return sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(res),
	), nil
}

func initMetricExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	authHeader := "Basic " + "admin:admin"

	return otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("otel-collector:4317"),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()),
		otlpmetricgrpc.WithHeaders(map[string]string{"Authorization": authHeader}),
	)
}

func collectAndExportMetrics(ctx context.Context, reader *sdkmetric.ManualReader, exporter sdkmetric.Exporter) {
	for {
		rm := &metricdata.ResourceMetrics{}
		if err := reader.Collect(ctx, rm); err != nil {
			log.Printf("Failed to collect metrics: %v", err)
			continue
		}

		if err := exporter.Export(ctx, rm); err != nil {
			log.Printf("Failed to export metrics: %v", err)
		}

		time.Sleep(60 * time.Second)
	}
}

func initHistogram(ctx context.Context, meterProvider *sdkmetric.MeterProvider, serviceName string) {
	meter := meterProvider.Meter(serviceName + "-meter")

	histogram, err := meter.Int64Histogram("hub-utilization",
		metric.WithDescription("Sum of parameters a, b, and c"))
	if err != nil {
		fmt.Printf("Failed to create histogram: %v\n", err)
	}

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				sum := aValue + bValue + cValue
				histogram.Record(ctx, sum)
				log.Printf("Interval recorded with sum = %d, current value of (a,b,c) = %d, %d, %d", sum, aValue, bValue, cValue)
				mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var params HubParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		mu.Lock()
		if params.ParamA != nil {
			aValue = *params.ParamA
		}
		if params.ParamB != nil {
			bValue = *params.ParamB
		}
		if params.ParamC != nil {
			cValue = *params.ParamC
		}
		mu.Unlock()

		log.Printf("Updated parameters: a=%d, b=%d, c=%d", aValue, bValue, cValue)
		fmt.Fprintf(w, "Updated parameters: a=%d, b=%d, c=%d", aValue, bValue, cValue)
	})
}

func main() {
	instanceID = os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		log.Fatal("INSTANCE_ID environment variable not set")
	}

	ctx := context.Background()

	resource, err := newResource()
	if err != nil {
		log.Fatalf("failed to initialize hub resource: %v", err)
	}

	reader := newManualReader()

	meterProvider, err := initMeterProvider(resource, reader)
	if err != nil {
		log.Fatalf("failed to initialize hub meter provider: %v", err)
	}
	defer meterProvider.Shutdown(ctx)

	exporter, err := initMetricExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize metric exporter: %v", err)
	}

	go collectAndExportMetrics(ctx, reader, exporter)
	go initHistogram(ctx, meterProvider, "hub")

	log.Println("Server is up and running on port 8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}
