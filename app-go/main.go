package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/attribute"
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

	instanceID string
)

type HubParams struct {
	ParamA int64 `json:"hub_param_a"`
	ParamB int64 `json:"hub_param_b"`
	ParamC int64 `json:"hub_param_c"`
}

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
		otlpmetricgrpc.WithEndpoint("otel-collector:4317"),
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
		metric.WithDescription("status of service A"),
	)
	if err != nil {
		return fmt.Errorf("failed to create gauge for param_a: %w", err)
	}

	bGauge, err = meter.Int64ObservableGauge(
		serviceName+"_param_b",
		metric.WithDescription("status of service B"),
	)
	if err != nil {
		return fmt.Errorf("failed to create gauge for param_b: %w", err)
	}

	cGauge, err = meter.Int64ObservableGauge(
		serviceName+"_param_c",
		metric.WithDescription("status of service C"),
	)
	if err != nil {
		return fmt.Errorf("failed to create gauge for param_c: %w", err)
	}

	attributes := attribute.String("instance_id", instanceID)

	log.Printf("Recording metric: hub_param_a, value: %d, instance_id: %v\n", aValue, attributes)
	log.Printf("Recording metric: hub_param_b, value: %d, instance_id: %v\n", bValue, attributes)
	log.Printf("Recording metric: hub_param_c, value: %d, instance_id: %v\n", cValue, attributes)

	_, err = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		observer.ObserveInt64(aGauge, aValue, metric.WithAttributes(attributes))
		observer.ObserveInt64(bGauge, bValue, metric.WithAttributes(attributes))
		observer.ObserveInt64(cGauge, cValue, metric.WithAttributes(attributes))
		return nil
	}, aGauge, bGauge, cGauge)

	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	return nil
}

func handleJSONRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var params HubParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if (params.ParamA != 0 && params.ParamA != 1) || (params.ParamB != 0 && params.ParamB != 1) || (params.ParamC != 0 && params.ParamC != 1) {
		http.Error(w, "Invalid parameter value", http.StatusBadRequest)
		return
	}

	aValue = params.ParamA
	bValue = params.ParamB
	cValue = params.ParamC

	log.Printf("Updated parameters: a=%d, b=%d, c=%d", aValue, bValue, cValue)
	fmt.Fprintf(w, "Updated parameters: a=%d, b=%d, c=%d", aValue, bValue, cValue)
}

func main() {
	instanceID = os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		log.Fatal("INSTANCE_ID environment variable not set")
	}

	ctx := context.Background()

	resource, err := newResource("hub")
	if err != nil {
		log.Fatalf("failed to initialize hub resource: %v", err)
	}

	meterProvider, err := initMeterProvider(ctx, resource)
	if err != nil {
		log.Fatalf("failed to initialize hub meter provider: %v", err)
	}
	defer meterProvider.Shutdown(ctx)

	err = initGauges(meterProvider, "hub")
	if err != nil {
		log.Fatalf("failed to initialize hub metrics: %v", err)
	}

	http.HandleFunc("/update", handleJSONRequest)

	log.Println("Server is up and running on port 8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}
