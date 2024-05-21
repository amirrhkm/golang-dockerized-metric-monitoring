package main

import (
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var counter metric.Int64Counter

func main() {
	initMeter()

	// Acquire a meter from the global meter provider
	meter := otel.Meter("my-service-meter")

	// Create a new counter instrument
	var err error
	counter, err = meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Number of API calls."),
	)
	if err != nil {
		log.Fatalf("failed to create counter: %v", err)
	}

	// Set up the HTTP handler
	http.HandleFunc("/user/", userHandler)

	log.Println("Server is up and running on port 8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/user/"):]
	data := []attribute.KeyValue{
		attribute.String("route", "/user/:name"),
		attribute.String("name", name),
	}
	// Convert the attributes to AddOption
	options := []metric.AddOption{
		metric.WithAttributes(data...),
	}
	counter.Add(r.Context(), 1, options...)
	log.Printf("Received request for name: %s", name)
	fmt.Fprintf(w, "Hello %s", name)
}
