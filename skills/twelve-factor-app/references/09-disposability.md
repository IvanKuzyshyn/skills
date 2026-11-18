# IX. Disposability

Maximize robustness with fast startup and graceful shutdown.

## Rules

1. Processes can be started or stopped at a moment's notice, which enables elastic scaling and fast deploys.
2. Startup is fast: a few seconds from launch to ready.
3. Processes shut down gracefully on SIGTERM. A web process stops accepting new requests and finishes in-flight ones. A worker returns its current job to the queue (NACK, release the lock).
4. Processes are robust against sudden termination: in-progress work is recoverable or retried automatically without data loss (crash-only design).

## Violation signals

- No SIGTERM handling, or a handler that exits immediately while requests or jobs are in flight.
- The process runs under a shell wrapper that does not forward signals (for example, `sh -c` as PID 1 in a container).
- Long startup work: warming large caches, running migrations, or blocking on slow external calls before serving.
- Jobs acknowledged before processing finishes, or job state kept only in memory.
- Multi-step writes with no transaction, idempotency key, or resumable state.

## Compliant patterns

- A SIGTERM handler that closes the listener, drains in-flight work within a timeout, closes connections, then exits.
- Readiness and liveness checks, so the platform routes traffic only to ready processes.
- Queue consumers that acknowledge after completion, with retries and idempotent handlers.
- Lazy or background warmup that does not block readiness.

## Legitimate exceptions

- Startup may be slower when the runtime itself is slow to start (some JVM apps); focus on avoidable startup work.
