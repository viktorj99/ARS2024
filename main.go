package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/configuration"
	"projekat/handlers"
	"projekat/repository"
	"projekat/service"
	"time"

	_ "projekat/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"golang.org/x/time/rate"
)

// @title Configuration API
// @version 1.0
// @description This is a sample server for a configuration service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	cfg := configuration.GetConfiguration()

	// Initialize OpenTelemetry
	ctx := context.Background()
	exp, err := newExporter(cfg.JaegerEndpoint)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	tp := newTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer := tp.Tracer("config-service")

	repo, err := repository.NewConfigConsulRepository(tracer)
	if err != nil {
		log.Fatalf("Failed to create Consul repository: %v", err)
	}

	repoGroup, err := repository.NewConfigGroupConsulRepository()
	if err != nil {
		log.Fatalf("Failed to create Consul repository: %v", err)
	}

	idempotencyRepo, err := repository.NewIdempotencyRepository()
	if err != nil {
		log.Fatalf("Failed to create Idempotency repository: %v", err)
	}

	handlers.SetIdempotencyRepository(idempotencyRepo)

	serviceConfig := service.NewConfigService(repo, tracer)
	serviceConfigGroup := service.NewConfigGroupService(repoGroup, tracer)

	handlerConfig := handlers.NewConfigHandler(serviceConfig, tracer)
	handlerConfigGroup := handlers.NewConfigGroupHandler(serviceConfigGroup, serviceConfig, tracer)

	router := mux.NewRouter()

	limiter := rate.NewLimiter(0.167, 10)

	router.Handle("/configs/{name}/{version}", handlers.RateLimit(limiter, Count(http.HandlerFunc(handlerConfig.GetConfig), "GET", "/configs/{name}/{version}"))).Methods("GET")
	router.Handle("/configs", handlers.RateLimit(limiter, handlers.IdempotencyMiddleware(Count(http.HandlerFunc(handlerConfig.AddConfig), "POST", "/configs")))).Methods("POST")
	router.Handle("/configs/{name}/{version}", handlers.RateLimit(limiter, Count(http.HandlerFunc(handlerConfig.DeleteConfig), "DELETE", "/configs/{name}/{version}"))).Methods("DELETE")

	router.Handle("/configGroups/{name}/{version}", handlers.RateLimit(limiter, Count(http.HandlerFunc(handlerConfigGroup.GetConfigGroup), "GET", "/configGroups/{name}/{version}"))).Methods("GET")
	router.Handle("/configGroups", handlers.RateLimit(limiter, handlers.IdempotencyMiddleware(Count(http.HandlerFunc(handlerConfigGroup.AddConfigGroup), "POST", "/configGroups")))).Methods("POST")
	router.Handle("/configGroups/{name}/{version}", handlers.RateLimit(limiter, Count(http.HandlerFunc(handlerConfigGroup.DeleteConfigGroup), "DELETE", "/configGroups/{name}/{version}"))).Methods("DELETE")
	router.Handle("/configGroups/{name}/{version}", handlers.RateLimit(limiter, handlers.IdempotencyMiddleware(Count(http.HandlerFunc(handlerConfigGroup.AddConfigToGroup), "POST", "/configGroups/{name}/{version}")))).Methods("POST")
	router.Handle("/configGroups/{groupName}/{groupVersion}/{configName}/{configVersion}", handlers.RateLimit(limiter, Count(http.HandlerFunc(handlerConfigGroup.DeleteConfigFromGroup), "DELETE", "/configGroups/{groupName}/{groupVersion}/{configName}/{configVersion}"))).Methods("DELETE")
	router.Handle("/configGroups/{groupName}/{groupVersion}/{labels}", handlers.RateLimit(limiter, Count(http.HandlerFunc(handlerConfigGroup.GetConfigsFromGroupByLabels), "GET", "/configGroups/{groupName}/{groupVersion}/{labels}"))).Methods("GET")
	router.Handle("/configGroups/{groupName}/{groupVersion}/{labels}", handlers.RateLimit(limiter, Count(http.HandlerFunc(handlerConfigGroup.DeleteConfigsFromGroupByLabels), "DELETE", "/configGroups/{groupName}/{groupVersion}/{labels}"))).Methods("DELETE")

	// Swagger documentation route
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	router.Path("/metrics").Handler(MetricsHandler())

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stop
	signal.Stop(stop)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func newExporter(address string) (*jaeger.Exporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(address)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("config-service"),
	)

	return sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exp),
		sdktrace.WithResource(r),
	)
}
