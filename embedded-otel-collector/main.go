package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	otlpexporter "go.opentelemetry.io/collector/exporter/otlpexporter"
	otlphttpexporter "go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/collector/service"
)

var (
	DefaultOTLPHTTPEndpoint = "http://localhost:4317"
)

func generateFactories() otelcol.Factories {
	receiversFactory, err := receiver.MakeFactoryMap(otlpreceiver.NewFactory())
	if err != nil {
		log.Fatalf("failed to make receiver factory map: %s", err)
	}

	exportersFactory, err := exporter.MakeFactoryMap(otlphttpexporter.NewFactory(), otlpexporter.NewFactory())
	if err != nil {
		log.Fatalf("failed to make exporter factory map: %s", err)
	}

	return otelcol.Factories{
		Receivers: receiversFactory,
		Exporters: exportersFactory,
	}
}

func generateConfigMap() *confmap.Conf {
	otlphttp_endpoint := os.Getenv("OTLPHTTP_ENDPOINT")
	if otlphttp_endpoint == "" {
		otlphttp_endpoint = DefaultOTLPHTTPEndpoint
	}

	return confmap.NewFromStringMap(map[string]interface{}{
		"receivers": map[string]interface{}{
			"otlp": map[string]interface{}{
				"protocols": map[string]interface{}{
					"grpc": nil,
					"http": nil,
				},
			},
		},
		"exporters": map[string]interface{}{
			"otlphttp": map[string]interface{}{
				"endpoint": otlphttp_endpoint,
			},
		},
		"service": map[string]interface{}{
			"pipelines": map[string]interface{}{
				"traces": map[string]interface{}{
					"receivers": []interface{}{"otlp"},
					"exporters": []interface{}{"otlphttp"},
				},
			},
		},
	})
}

func main() {
	ctx := context.Background()

	factories := generateFactories()
	configMap := generateConfigMap()

	otelCfgSettings, err := otelcol.Unmarshal(configMap, factories)
	if err != nil {
		log.Fatalf("failed to make otel config: %s", err)
	}

	otelCfg := otelcol.Config{
		Receivers: otelCfgSettings.Receivers.Configs(),
		Exporters: otelCfgSettings.Exporters.Configs(),
		Service:   otelCfgSettings.Service,
	}

	if err := otelCfg.Validate(); err != nil {
		log.Fatalf("failed to validate otel config: %s", err)
	}

	svc, err := service.New(ctx, service.Settings{
		Receivers:  receiver.NewBuilder(otelCfg.Receivers, factories.Receivers),
		Exporters:  exporter.NewBuilder(otelCfg.Exporters, factories.Exporters),
		Connectors: connector.NewBuilder(otelCfg.Connectors, factories.Connectors),
	}, otelCfg.Service)

	if err != nil {
		log.Fatalf("failed to create Otel service: %s", err)
	}

	// start the service
	if err := svc.Start(ctx); err != nil {
		log.Fatalf("failed to start Otel service: %s", err)
	}

	// wait for a signal to stop the service
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()

	// Stop the service
	if err := svc.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shut down service cleanly: %s", err)
	}

}
