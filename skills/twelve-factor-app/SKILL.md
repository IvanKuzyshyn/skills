---
name: twelve-factor-app
description: Design or evaluate deployable services (web apps, APIs, workers, schedulers) against the Twelve-Factor App methodology (codebase, dependencies, config, backing services, build/release/run, processes, port binding, concurrency, disposability, dev/prod parity, logs, admin processes). Use when the user asks for a 12-factor or twelve-factor review, a cloud-readiness, deployability, or horizontal-scaling audit of a service codebase, or is designing how a new service will be configured, deployed, scaled, and operated. Do NOT use for libraries, SDKs, end-user CLIs, mobile, desktop, or front-end-only apps, database schema or API contract design, general code-quality reviews (use idiomatic-go for Go code quality), or a single how-to question such as adding graceful shutdown or reading one setting from env — answer those directly.
---

# Twelve-Factor App

The twelve factors define a clean contract between an app and its execution platform. This skill applies them in two modes: **evaluate** an existing codebase, or **design** a service so it complies from the start.

## The factors

Each line is the rule. Open the reference file only when you need detail: violation signals, compliant patterns, or exceptions.

| # | Factor | Rule | Reference |
|---|---|---|---|
| I | Codebase | One codebase in version control per app, many deploys | [01-codebase.md](references/01-codebase.md) |
| II | Dependencies | Declare all dependencies explicitly and isolate them; no implicit system packages or tools | [02-dependencies.md](references/02-dependencies.md) |
| III | Config | Everything that varies between deploys lives in env vars, never in code or env-named groups | [03-config.md](references/03-config.md) |
| IV | Backing services | Every network service is an attached resource, located by config, swappable without code changes | [04-backing-services.md](references/04-backing-services.md) |
| V | Build, release, run | Strictly separate stages; releases are immutable with unique IDs; keep run minimal | [05-build-release-run.md](references/05-build-release-run.md) |
| VI | Processes | Stateless, share-nothing processes; persistent state lives in backing services | [06-processes.md](references/06-processes.md) |
| VII | Port binding | The app is self-contained and exports its service by binding to a port | [07-port-binding.md](references/07-port-binding.md) |
| VIII | Concurrency | Scale out by process type; never daemonize or write PID files | [08-concurrency.md](references/08-concurrency.md) |
| IX | Disposability | Fast startup, graceful shutdown on SIGTERM, robust against sudden death | [09-disposability.md](references/09-disposability.md) |
| X | Dev/prod parity | Minimize time, personnel and tools gaps; same backing service types and versions everywhere | [10-dev-prod-parity.md](references/10-dev-prod-parity.md) |
| XI | Logs | Write an unbuffered event stream to stdout; never route or store logs from the app | [11-logs.md](references/11-logs.md) |
| XII | Admin processes | Run one-off tasks as separate processes against the same release and config | [12-admin-processes.md](references/12-admin-processes.md) |

## Choose the mode

- **Evaluate**: an existing codebase and the user wants a review, audit, or readiness check.
- **Design**: a new service, or a planned change to how a service is built, configured, deployed or scaled.
- **Targeted question**: the user asks how to satisfy one factor in their service (for example, handling SIGTERM or moving a setting to env). Answer it directly using that factor's reference file; don't run a full evaluation or ask for a mode.
- If the request fits none of these clearly, ask which one the user wants.
- If the target is not a deployable service (a library, SDK, end-user CLI, or front-end-only app), say that most factors don't apply and evaluate only the ones that do, such as dependencies and config.

## Where the result goes

- If the skill is used while writing or editing a design document (the user is working on a doc, or names one), write the result into that document, in the section the user indicates or a new "Twelve-Factor" section.
- Otherwise, answer in chat.

## Evaluate mode

1. Identify the app boundary: which deployable apps live in the repo and what process types each runs (web, worker, scheduler, CLI). Evaluate each app separately if there are several.
2. For each factor, search for the violation signals listed in its reference file. Read the reference file for any factor where you find a signal or the evidence is ambiguous.
3. Classify each factor:
   - **Compliant**: evidence shows the rule is followed.
   - **Partial**: followed in some places, broken in others.
   - **Violation**: concrete evidence the rule is broken.
   - **Unknown**: the repo cannot show it (for example, deployment platform configuration kept elsewhere). Say what would settle it.
   - **N/A**: the factor does not apply (for example, port binding for a batch job).
4. Never report a violation without evidence. Cite `file:line` for every Partial or Violation.
5. Check whether an apparent violation is a legitimate exception listed in the reference file before reporting it.

Output:

```markdown
## Twelve-Factor evaluation: <app>

| Factor | Status | Evidence | Fix |
|---|---|---|---|
| III. Config | Violation | `src/db.ts:12` hardcodes the database URL | Read `DATABASE_URL` from the environment |
| ... | | | |

### Priority fixes
1. <highest-impact violation>: why it matters, what to change.
```

Order priority fixes by risk: leaked credentials and data loss first, then anything that blocks horizontal scaling or safe deploys, then the rest.

## Design mode

1. Gather the facts the design depends on: language and runtime, process types, backing services, deployment target, scaling and state needs. Ask for anything missing instead of assuming.
2. For each factor, make a concrete decision for this service, using the compliant patterns in the reference files. Name the actual mechanism (for example, "sessions in Redis with TTL", "migrations as a release-phase job"), not the principle.
3. Mark factors that do not apply and say why.
4. Call out tradeoffs and the places where the design deliberately deviates, with the reason.

Output:

```markdown
## Twelve-Factor design: <service>

| Factor | Decision | Notes |
|---|---|---|
| III. Config | All config from env vars, validated at startup; fail fast on missing values | Secrets injected by the platform |
| ... | | |

### Tradeoffs and open questions
- ...
```
