# IV. Backing services

Treat backing services as attached resources.

## Rules

1. A backing service is any service the app consumes over the network: datastores, queues, SMTP, caches, object storage, third-party APIs.
2. Each distinct backing service is a resource that can be attached to or detached from a deploy at will. Two databases are two resources.
3. The app references every service through a URL or locator plus credentials stored in config, and makes no distinction between local and third-party services. Swapping a local database for a managed one changes config only.

## Violation signals

- Connection details for any service hardcoded or assembled from constants.
- Code that branches on whether a service is local or managed (`if host == "localhost"`).
- Assumptions that a service runs on the same host: Unix sockets, `localhost` defaults relied on in production, a shared local filesystem used as a queue.

## Compliant patterns

- One env var per resource (`DATABASE_URL`, `REDIS_URL`, `SMTP_URL`, `S3_BUCKET` plus credentials).
- Clients created from config at startup.
- Health checks that report each attached resource separately.

## Legitimate exceptions

- Provider-specific SDKs are fine; the rule is about locating the resource through config, not about abstracting every vendor.
