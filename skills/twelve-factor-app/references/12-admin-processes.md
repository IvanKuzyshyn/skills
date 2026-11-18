# XII. Admin processes

Run admin/management tasks as one-off processes.

## Rules

1. The app distinguishes its regular process formation from one-off admin tasks: migrations, REPL consoles, one-time scripts.
2. Admin processes run against a release, in the same environment as long-running processes: same codebase, same config, same dependency isolation. Admin code ships with the app.

## Violation signals

- Migrations or data fixes run from a developer laptop against production with local code or local config.
- Admin scripts kept outside the repo, or needing dependencies the app does not declare.
- Migrations run automatically at the start of every web process, so many replicas race to migrate.
- Admin actions exposed only as hidden HTTP endpoints in the web process.
- Manual SQL run directly against production as a routine step.

## Compliant patterns

- Migrations and maintenance scripts committed to the repo and run with the same image and isolation (`docker run <image> migrate`, a Kubernetes Job from the release image, `heroku run`).
- Migrations as a release-phase step or a one-off job that runs once per release.
- A REPL opened in a one-off process from the current release.

## Legitimate exceptions

- Running migrations on startup is acceptable for a single-instance service with a migration lock; flag it as Partial if the service can scale out.
