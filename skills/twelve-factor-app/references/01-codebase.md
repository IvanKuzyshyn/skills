# I. Codebase

One codebase tracked in revision control, many deploys.

## Rules

1. The app is always tracked in version control.
2. Codebase and app map one to one. Several codebases make a distributed system, where each component is its own app. Several apps sharing one codebase is a violation.
3. One codebase is deployed to many environments (production, staging, every developer's machine). Deploys may run different versions of the same codebase.

## Violation signals

- Code copied between repos or directories instead of shared through a package.
- Separate branches or forks per environment (`prod`, `staging`) that never merge back.
- One repo that builds several independently deployed apps from overlapping code, with no clear app boundary or shared library.
- Deploy scripts that pull code from more than one repo to assemble a single app.

## Compliant patterns

- One repo per app, with shared code extracted into versioned libraries pulled in through the dependency manager.
- A monorepo where each app has its own directory, build and deploy, and shared code lives in separate internal packages with explicit dependencies.
- Environments differ only by which commit is deployed and by config.

## Legitimate exceptions

- Monorepos are fine as long as each app is a distinct, separately built and deployed unit.
