import os
from flask import Flask
from urllib.parse import urljoin
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter as OTLPHTTPSpanExporter
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter as OTLPGRPCSpanExporter
from opentelemetry.instrumentation.flask import FlaskInstrumentor
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.sdk.resources import Resource

app = Flask(__name__)

otlp_endpoint = os.getenv("OTLP_ENDPOINT")
if otlp_endpoint is None:
    raise ValueError("OTLP_ENDPOINT environment variable must be set")

# Define resource attributes
resource = Resource(attributes={
    "service.name": "my-flask-app-4",
    "serverToken": "xxxx-xxxx-xxxx-xxxx",
})
#resource = Resource(attributes={})

# Set the tracer provider and a batch span processor with OTLP exporter
trace.set_tracer_provider(TracerProvider(resource=resource))
trace.get_tracer_provider().add_span_processor(
    BatchSpanProcessor(
        #OTLPHTTPSpanExporter(endpoint=urljoin(otlp_endpoint, "/v1/traces")),
        #OTLPGRPCSpanExporter(endpoint=otlp_endpoint, insecure=True),
        OTLPGRPCSpanExporter(endpoint="localhost:4317", insecure=True), # grpc
    )
)

# Instrument the Flask app to automatically generate spans
FlaskInstrumentor().instrument_app(app)

@app.route('/')
def hello():
    return 'Hello, OpenTelemetry!'

if __name__ == '__main__':
    app.run(debug=True, port=5001)
