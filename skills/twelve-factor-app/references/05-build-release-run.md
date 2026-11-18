# V. Build, release, run

Strictly separate build and run stages.

## Rules

1. Three distinct stages:
   - **Build** turns a specific commit into an executable bundle: fetch dependencies, compile binaries and assets.
   - **Release** combines a build with a deploy's config.
   - **Run** launches processes against a selected release.

   Code cannot change at runtime.
2. Every release has a unique ID (timestamp, incrementing number, image digest) and is immutable. Any change creates a new release, which makes rollback possible.
3. Push complexity into build and keep run minimal. Run happens unattended (restarts, reboots), so it must have as few moving parts as possible.

## Violation signals

- Dependencies installed, assets compiled, or code generated at process start (`npm install` or `pip install` in the start command or entrypoint).
- Pulling code at runtime (pulling from version control on servers, downloading plugins at startup).
- Deploys that edit files in place on servers instead of shipping a new artifact.
- Mutable artifact tags as the only identifier (`:latest`), with no way to name or roll back to a previous release.
- Config baked into the build artifact, so the same build cannot be promoted across environments.

## Compliant patterns

- CI builds one artifact per commit (container image, binary, package) tagged with the commit SHA or digest.
- The platform combines that artifact with per-deploy config to create a release; rollback means selecting an earlier release.
- The start command only launches the process.

## Legitimate exceptions

- Interpreted languages that load source at runtime are fine as long as the source comes from the immutable build artifact.
