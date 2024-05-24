### Run custom collector:

./my-otelcol1 --config=../otelcol-config.yaml --set=service.telemetry.logs.level=DEBUG

### Run telemetrygen

telemetrygen traces --otlp-endpoint "0.0.0.0:4317" --otlp-insecure --traces 1
