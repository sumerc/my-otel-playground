import os
import logging
from flask import Flask
from urllib.parse import urljoin
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter as OTLPHTTPSpanExporter
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter as OTLPGRPCSpanExporter
from opentelemetry.instrumentation.flask import FlaskInstrumentor
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.sdk.resources import Resource

otlp_endpoint = os.getenv("OTLP_ENDPOINT")
if otlp_endpoint is None:
    raise ValueError("OTLP_ENDPOINT environment variable must be set")


# log exporter config start (still experimental in the SDK, but useful for testing)
# normally the standart way is to do this from an OTEL collector where you receive a
# log in any format and export it
from opentelemetry.instrumentation.logging import LoggingInstrumentor
from opentelemetry.sdk._logs import LoggerProvider, LoggingHandler
from opentelemetry.sdk._logs.export import BatchLogRecordProcessor
from opentelemetry._logs import set_logger_provider
from opentelemetry.exporter.otlp.proto.http._log_exporter import (
    OTLPLogExporter,
)

# generate the root logger
LoggingInstrumentor().instrument(set_logging_format=True) # should be done before logger init
logger = logging.getLogger(__name__)

logger_provider = LoggerProvider(
    resource=Resource.create(
        {
            "service.name": "cipo",
            "service.instance.id": "instance-1",
        }
    ),
)
set_logger_provider(logger_provider)
exporter = OTLPLogExporter(endpoint=urljoin(otlp_endpoint, "/v1/logs"))
logger_provider.add_log_record_processor(BatchLogRecordProcessor(exporter))
handler = LoggingHandler(level=logging.NOTSET, logger_provider=logger_provider)

# Attach OTLP handler to root logger
logger.addHandler(handler)

# log exporter config end

app = Flask(__name__)

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
        OTLPHTTPSpanExporter(endpoint=urljoin(otlp_endpoint, "/v1/traces")),
        #OTLPHTTPSpanExporter(endpoint=otlp_endpoint),
        #OTLPGRPCSpanExporter(endpoint=otlp_endpoint, insecure=True),
        #OTLPGRPCSpanExporter(endpoint="localhost:4317", insecure=True), # grpc
    )
)

# Instrument the Flask app to automatically generate spans
FlaskInstrumentor().instrument_app(app)

@app.route('/')
def hello():
    logger.info("hello from Flask!")
    return 'Hello, OpenTelemetry!'

if __name__ == '__main__':
    app.run(debug=True, port=5001)
