# X. Dev/prod parity

Keep development, staging, and production as similar as possible.

## Rules

1. Minimize three gaps:
   - **Time gap**: code reaches production in hours or minutes, not weeks.
   - **Personnel gap**: the people who write code deploy it and watch it in production.
   - **Tools gap**: development and production use the same stack.
2. All deploys use the same type and version of each backing service. Resist lightweight local substitutes (SQLite locally and PostgreSQL in production; a memory cache locally and Redis in production), because small incompatibilities break code in production.

## Violation signals

- Different database engines between local, test, and production configs.
- In-memory or filesystem substitutes for queues, caches, or object storage outside unit tests.
- Local setup running different service versions than production (for example, `postgres:12` locally, 16 in production).
- No CI/CD pipeline, or manual, infrequent release steps.
- Code paths that run only in production (`if production: use_real_queue()`).

## Compliant patterns

- A `docker-compose.yml` or dev container that runs the same service types and pinned versions as production.
- Integration tests against real backing services (Testcontainers or similar).
- Continuous deployment from the main branch, with developers owning releases.

## Legitimate exceptions

- Unit tests may use fakes or in-memory doubles; parity applies to running deploys and integration tests.
- A managed production service (for example, RDS) with a local container of the same engine and version is compliant.
