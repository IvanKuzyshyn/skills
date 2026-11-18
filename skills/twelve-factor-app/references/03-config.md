# III. Config

Store config in the environment.

## Rules

1. Config is strictly separate from code. Config is everything likely to vary between deploys: resource handles, credentials, per-deploy values like the canonical hostname. Litmus test: the codebase could be open-sourced right now without leaking a credential.
2. Config lives in environment variables, not in code constants or in untracked config files scattered in different formats.
3. Env vars are granular controls managed per deploy, never grouped into named environments (`development`, `production`, `joes-staging`).

Internal app config that does not vary between deploys (routes, dependency injection wiring) is not config in this sense and belongs in code.

## Violation signals

- Hardcoded hostnames, URLs, ports, credentials, API keys or tokens in source.
- Committed `.env` files containing real values, or secrets in committed YAML/JSON/properties files.
- Per-environment config files or blocks (`config/production.yml`, `settings_prod.py`, `application-prod.properties`) selected by an environment name.
- Branching on the environment name (`if NODE_ENV == "production"`, `if env == "staging"`) to choose values, instead of a dedicated variable per setting.
- Defaults that point at production resources.

## Compliant patterns

- Every deploy-varying value read from an env var, gathered in one config module and validated at startup; the process fails fast on missing or invalid values.
- A committed `.env.example` with names and dummy values; real values injected by the platform or a secret manager.
- Behavior controlled by specific variables (`CACHE_TTL_SECONDS`, `ENABLE_X`) rather than by environment name.

## Legitimate exceptions

- Secrets mounted as files by the platform (for example, Kubernetes secret volumes) keep config out of code and per deploy, which meets the intent, provided the path comes from config.
- Safe local-development defaults for non-secret values are acceptable when production never relies on them.
- Framework environment-name flags used only for framework internals (for example, enabling debug tooling) are minor; report Partial only if they select config values.
