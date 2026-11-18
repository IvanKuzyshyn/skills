# VII. Port binding

Export services via port binding.

## Rules

1. The app is completely self-contained and exports HTTP by binding to a port and listening on it. It does not depend on a webserver injected at runtime (a PHP module inside Apache, a WAR inside Tomcat). The webserver is a declared dependency inside the app.
2. Any protocol can be exported this way, so one app can become a backing service for another through a URL in the consumer's config.

## Violation signals

- The app only runs when deployed into an external container (WAR/EAR deployed to an app server, a module loaded into Apache or IIS).
- The listening port hardcoded instead of read from config (`PORT`).
- The app depends on a reverse proxy on the same host for basic serving.

## Compliant patterns

- An embedded server started by the app (Go `net/http`, Node `http`, Spring Boot embedded server, Uvicorn or Gunicorn declared as dependencies).
- Port and bind address from config, defaulting to a sensible local value.
- A routing layer or load balancer in front, managed by the platform, not by the app.

## Legitimate exceptions

- Workloads that serve no network requests (batch jobs, queue consumers) do not bind a port; mark N/A.
- An application server launched from the app's own declared dependencies (for example, Gunicorn in front of a WSGI app) is compliant.
