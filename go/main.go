package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func initConn() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		"0.0.0.0:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	log.Println("Successfully connected to collector receiver")
	return conn, nil
}

func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String("my-service"),
			semconv.ServiceVersionKey.String("0.1.0"),
		))
}

func initMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:admin"))

	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
		otlpmetricgrpc.WithHeaders(map[string]string{"Authorization": authHeader}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}
	log.Println("Successfully initialized meter provider with basic authentication")

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}

func initMeter() (func(context.Context) error, error) {
	ctx := context.Background()

	conn, err := initConn()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize connection: %w", err)
	}

	res, err := newResource()
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	shutdown, err := initMeterProvider(ctx, res, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize meter provider: %w", err)
	}

	return shutdown, nil
}

var counter metric.Int64Counter

func main() {
	shutdown, err := initMeter()
	if err != nil {
		log.Fatalf("failed to initialize meter: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown meter: %v", err)
		}
	}()

	meter := otel.Meter("my-service-meter")

	counter, err = meter.Int64Counter(
		"api.req_counter",
		metric.WithDescription("Number of API calls."))
	if err != nil {
		log.Fatalf("failed to create counter: %v", err)
	}

	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/user/"):]
		data := []attribute.KeyValue{
			attribute.String("route", "/user/:name"),
			attribute.String("name", name),
		}

		counter.Add(r.Context(), 1, metric.WithAttributes(data...))
		log.Printf("Received request for name: %s", name)

		fmt.Fprintf(w, "Hello %s", name)
	})

	log.Println("Server is up and running on port 8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}
