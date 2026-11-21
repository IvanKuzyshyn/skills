#!/usr/bin/env bash
# Run automated Go checks for a quality review and print results in labelled sections.
# Usage: run_checks.sh [--no-test] [module-root]
# Every check runs even if an earlier one fails; failures are reported, not fatal.

set -u

run_tests=1
if [ "${1:-}" = "--no-test" ]; then
  run_tests=0
  shift
fi
root="${1:-.}"
cd "$root" || { echo "cannot cd to $root"; exit 1; }

section() { printf '\n===== %s =====\n' "$1"; }

# Prints the command's output and its exit status so the reviewer can tell
# "ran clean" apart from "failed" and "not run".
run() {
  local name="$1"; shift
  section "$name"
  echo "\$ $*"
  "$@" 2>&1
  echo "[exit status: $?]"
}

section "environment"
if ! command -v go >/dev/null 2>&1; then
  echo "go toolchain not found; automated checks skipped"
  exit 0
fi
go version
if [ ! -f go.mod ]; then
  echo "no go.mod in $(pwd)"
  exit 0
fi
grep -E '^(module|go|toolchain) ' go.mod
echo "packages: $(go list ./... 2>/dev/null | wc -l | tr -d ' ')"

run "go build" go build ./...
run "go vet" go vet ./...

section "gofmt -l (unformatted files)"
# vendor/ is excluded because it is not the project's own code.
unformatted=$(find . -name '*.go' -not -path './vendor/*' -print0 | xargs -0 gofmt -l 2>&1)
if [ -n "$unformatted" ]; then echo "$unformatted"; else echo "(none)"; fi

run "go mod tidy -diff" go mod tidy -diff

if [ "$run_tests" -eq 1 ]; then
  # The race detector requires cgo; fall back so tests still run without it.
  if [ "$(go env CGO_ENABLED)" = "1" ]; then
    run "go test -race" go test -race -count=1 -timeout 10m -cover ./...
  else
    echo "(cgo disabled: running tests without -race)"
    run "go test" go test -count=1 -timeout 10m -cover ./...
  fi
else
  section "go test"
  echo "skipped (--no-test)"
fi

for tool in staticcheck golangci-lint govulncheck; do
  if command -v "$tool" >/dev/null 2>&1; then
    case "$tool" in
      staticcheck) run staticcheck staticcheck ./... ;;
      golangci-lint) run golangci-lint golangci-lint run ./... ;;
      govulncheck) run govulncheck govulncheck ./... ;;
    esac
  else
    section "$tool"
    echo "not installed; skipped"
  fi
done
