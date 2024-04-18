# Barebone Embedded OpenTelemetry Collector

This Go application embeds an OpenTelemetry Collector to process and export telemetry data. It demonstrates simplest possible configuration and usage of receivers and exporters directly within a Go application.

## Running

Set the `OTLPHTTP_ENDPOINT` environment variable to configure the OTLP HTTP exporter endpoint:

```
export OTLPHTTP_ENDPOINT="http://your-otel-endpoint:port"
```

and

```
make run
```

to run the application.

It will be listeting on default HTTP/GRPC otel-collector ports and export to the endpoint given using OTLP over HTTP.(which can also be configured from the code)
