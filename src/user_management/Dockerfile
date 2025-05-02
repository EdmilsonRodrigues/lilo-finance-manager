FROM python:3.13

ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONBREAKPOINT=0

RUN pip install uv --no-cache-dir

WORKDIR /app

COPY pyproject.toml uv.lock .

RUN uv sync --group telemetry --no-cache
RUN uv run opentelemetry-bootstrap -a install

ENV OTEL_SERVICE_NAME=lfm-user-management-service
ENV OTEL_EXPORTER_OTLP_INSECURE=true
ENV OTEL_LOGS_EXPORTER=otlp
ENV OTEL_METRICS_EXPORTER=otlp
ENV OTEL_TRACES_EXPORTER=otlp
# ENV OTEL_PYTHON_LOGGING_AUTO_INSTRUMENTATION_ENABLED=true

COPY app/ ./app

EXPOSE 5000

ENTRYPOINT ["uv",  "run", "opentelemetry-instrument"]
CMD ["gunicorn", "--worker-class", "gevent", "-b", "0.0.0.0:5000", "app.app:create_app()"]
