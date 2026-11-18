# II. Dependencies

Explicitly declare and isolate dependencies.

## Rules

1. The app builds and runs in a deterministic environment: a new developer needs only the language runtime and dependency manager, then one build command.
2. The app never relies on implicitly installed system-wide packages. It declares every dependency exactly in a manifest and uses isolation at runtime. Declaration and isolation are both required.
3. The app never relies on implicitly installed system tools (for example, shelling out to `curl` or ImageMagick). Tools the app needs are vendored or declared.

## Violation signals

- No dependency manifest, or a manifest without a committed lockfile (`package-lock.json`, `poetry.lock`, `go.sum`, `Cargo.lock`, `Gemfile.lock`).
- Unpinned or floating versions (`*`, `latest`, open ranges) for runtime dependencies.
- `subprocess`, `exec`, `system`, or backtick calls to binaries not installed by the build (`curl`, `convert`, `ffmpeg`).
- Setup docs that say "install X globally", or `apt-get install` steps outside the build definition.
- Imports that resolve only because a package is installed globally.

## Compliant patterns

- Manifest plus lockfile, installed with a frozen or locked install (`npm ci`, `pip install --require-hashes`, `go mod download`).
- Runtime isolation: virtualenv, `node_modules`, `bundle exec`, static linking, or a container image that installs everything it needs.
- System tools installed explicitly in the build (for example, in the Dockerfile) or replaced by a library.

## Legitimate exceptions

- A container image built from the repo counts as the isolation boundary, as long as every system package it uses is declared in the image definition.
