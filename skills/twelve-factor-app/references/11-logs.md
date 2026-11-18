# XI. Logs

Treat logs as event streams.

## Rules

1. Logs are a continuous stream of time-ordered events, typically one event per line. Each process writes its stream, unbuffered, to stdout.
2. The app never routes or stores its output: no writing to or managing log files. The execution environment captures, collates, routes, and archives streams (Fluentd, platform log routers, indexing systems).

## Violation signals

- Loggers configured with file handlers or appenders (`FileHandler`, `RollingFileAppender`, `winston.transports.File`).
- Log rotation, retention, or archival logic inside the app.
- Logging libraries that ship directly to a vendor endpoint from the app as the only output.
- Buffered stdout in production (for example, Python without `PYTHONUNBUFFERED=1` or `-u`).
- Hardcoded log file paths in config or code.

## Compliant patterns

- The logger writes to stdout (or stderr), ideally as structured JSON with timestamp, level, and request or trace IDs.
- Log level from config.
- Shipping, indexing, and alerting done by the platform or a sidecar agent.

## Legitimate exceptions

- Sending traces or metrics through an agent or OpenTelemetry collector is a separate concern; this factor covers the log stream only.
