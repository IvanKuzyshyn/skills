# VI. Processes

Execute the app as one or more stateless processes.

## Rules

1. The app runs as one or more processes, from a single script up to many process types.
2. Processes are stateless and share nothing. Persistent data goes to a stateful backing service. Memory or local disk may be used only as a brief, single-transaction cache; never assume it survives to the next request or job, since another process may handle it and restarts wipe local state.
3. Sticky sessions are a violation. Session data belongs in a datastore with expiry (Redis, Memcached).

## Violation signals

- In-memory sessions (default memory session stores in Express, Flask, Spring, and similar).
- User uploads or generated files written to local disk and served later.
- In-process caches, maps, or counters treated as the source of truth across requests (rate limits, locks, job state, websocket membership).
- Load balancer configuration with session affinity or sticky cookies.
- Assets compiled into the filesystem at runtime instead of during build.
- In-process schedulers that assume a single instance (cron libraries inside the web process with no distributed lock).

## Compliant patterns

- Sessions in Redis/Memcached, or stateless signed tokens.
- Uploads streamed to object storage; local temp files deleted within the same request or job.
- Shared state (rate limits, locks, presence) in a backing service.
- Scheduled work in a dedicated process type or platform scheduler, with locking if more than one instance can run.

## Legitimate exceptions

- Per-process caches of immutable or easily rebuilt data (compiled templates, config, read-through caches) are fine when correctness does not depend on them.
