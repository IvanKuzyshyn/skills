# Go Quality Requirements

Each requirement has an ID, a severity, what to check, and why it matters. Severity reflects the typical impact; downgrade or upgrade one level when context clearly warrants it (e.g. an ignored error in a test helper is less serious than one in a payment path) and say so in the finding.

- **Critical** — can cause data loss, corruption, crashes, leaks, security exposure, or nondeterministic behavior in production.
- **Major** — makes the code materially harder to change correctly, or exposes callers to avoidable misuse.
- **Minor** — friction for readers and maintainers; no direct runtime impact.

## Contents

1. [Error handling (ERR)](#error-handling-err)
2. [Concurrency (CON)](#concurrency-con)
3. [Resource & runtime safety (RES)](#resource--runtime-safety-res)
4. [API & package design (API)](#api--package-design-api)
5. [Testing (TST)](#testing-tst)
6. [Readability & naming (RDB)](#readability--naming-rdb)
7. [Build & tooling (BLD)](#build--tooling-bld)

---

## Error handling (ERR)

### ERR-1 · Critical · No silently discarded errors
**Check:** Error results assigned to `_`, or calls whose error result is dropped entirely — especially `Close` on writable files, `Flush`, `Write`, `Encode`, `Commit`, `Rollback`, `Exec`.
**Why:** A dropped error turns a detectable failure into silent corruption or data loss. `Close` on a written file is often the first point where a failed write is reported.
**Acceptable:** Discards that are genuinely irrelevant (closing a read-only file, `fmt.Println` to stdout) — ideally with a brief comment when non-obvious.

### ERR-2 · Major · Each error is handled exactly once
**Check:** An error both logged and returned, or logged at several layers as it propagates.
**Why:** Duplicate handling produces repeated, misleading log lines and makes it unclear who is responsible for acting on the failure. Either handle it (log, retry, degrade, translate) or return it — not both.

### ERR-3 · Major · Errors carry context and are inspected structurally
**Check:** Bare `return err` across package or layer boundaries where the caller can't tell which operation failed; use of `%v` instead of `%w` when wrapping an error callers may inspect; comparison by `err.Error()` string matching or `==` against an error that may be wrapped.
**Why:** Without context, an error like `EOF` at the top of a call stack is undiagnosable. String matching breaks when messages change; `==` breaks when anything wraps the error. Use `fmt.Errorf("...: %w", err)` and `errors.Is` / `errors.As`.

### ERR-4 · Major · Panic and exit are reserved for unrecoverable states
**Check:** `panic` used for ordinary failures (bad input, I/O errors, not-found); `log.Fatal*` or `os.Exit` outside `main` (or outside `run()`-style functions called only by `main`).
**Why:** Library code that panics or exits takes the decision away from the caller and skips deferred cleanup. Panics are appropriate for programmer errors and impossible states, and `Must*` helpers used at init time.

### ERR-5 · Major · Results are not used when an error is returned
**Check:** Code that reads a returned value before or without checking the accompanying error; functions that return a partially valid value alongside a non-nil error without documenting it.
**Why:** By convention, a non-nil error means other results are undefined. Violating either side of that convention causes nil dereferences and use of half-built objects.

### ERR-6 · Minor · Error values follow conventions
**Check:** Error strings capitalized or ending in punctuation; `error` not the last return value; exported sentinel errors or error types that no caller needs to distinguish; custom error types returned as concrete pointer types (`*MyErr`) from functions whose signature returns `error`, which can produce non-nil interfaces holding nil pointers.
**Why:** Error strings are composed into larger messages, so casing and punctuation look wrong mid-sentence. Every exported error becomes API that must be maintained. The typed-nil case is a real bug (upgrade to Critical if it is reachable).

---

## Concurrency (CON)

### CON-1 · Critical · No data races
**Check:** `go test -race` failures; shared maps, slices, or struct fields read and written from multiple goroutines without a mutex, channel, or atomic; closures in goroutines capturing variables that are mutated afterwards.
**Why:** Races cause intermittent, unreproducible corruption and crashes (concurrent map writes abort the program).

### CON-2 · Critical · Every goroutine has a defined way to stop
**Check:** For each `go` statement, identify what makes it return: a cancelled context, a closed channel, a finite input, or completion. Look for sends on unbuffered channels whose receiver may have returned early (e.g. on first error or timeout), and infinite loops without an exit signal.
**Why:** A goroutine blocked forever is a memory and resource leak that grows with traffic and eventually takes the process down. If you can't explain when a goroutine exits, neither can the next maintainer.

### CON-3 · Major · Context is propagated correctly
**Check:** Functions doing I/O or long-running work that don't accept `context.Context`; `ctx` not the first parameter; `context.Background()`/`TODO()` used deep in request paths instead of the caller's context; contexts stored in struct fields; `cancel` functions not called (vet `lostcancel`).
**Why:** Without propagation, cancellation and deadlines stop at the first function that ignores them, so timed-out requests keep consuming resources.

### CON-4 · Major · Concurrency is the caller's choice
**Check:** Library or package APIs that start background goroutines implicitly (in constructors or `init`), or that return results only via channels where a plain synchronous call would do.
**Why:** A synchronous API can always be made concurrent by the caller with `go`; the reverse is awkward and hides lifecycle. Components that do need background work should expose explicit `Start`/`Run(ctx)` and `Close`/`Stop` so the owner controls the lifetime.

### CON-5 · Major · Channel ownership is clear
**Check:** Channels closed by a receiver, closed from more than one place, or possibly closed twice; sends that may happen after close.
**Why:** Sending on or closing a closed channel panics. The sender (owner) should be the only one that closes.

### CON-6 · Major · Synchronization primitives are used correctly
**Check:** `sync.Mutex`/`WaitGroup` values copied (vet `copylocks`); `wg.Add` called inside the goroutine instead of before `go`; locks held across blocking I/O or channel operations; missing `Unlock` on early-return paths (prefer `defer mu.Unlock()`).
**Why:** Each of these leads to races, deadlocks, or `Wait` returning too early.

---

## Resource & runtime safety (RES)

### RES-1 · Critical · Acquired resources are released
**Check:** Files, `http.Response.Body`, `sql.Rows`, `sql.Tx`, network connections, tickers (`time.NewTicker`) without a matching `Close`/`Stop` on every path. `defer x.Close()` placed before the error check on acquisition. `defer` inside a loop that iterates many times.
**Why:** Leaked file descriptors and connections exhaust pools and OS limits under load. Deferred calls in a loop don't run until the function returns.

### RES-2 · Critical · No avoidable panics at runtime
**Check:** Writes to possibly-nil maps; unchecked type assertions (`x.(T)` without `, ok`) on values that aren't guaranteed; indexing without bounds checks on external input; dereferencing pointers that may be nil per the function's contract.
**Why:** Each is a crash triggered by input or state rather than by a programmer error that tests would catch.

### RES-3 · Critical · No common security defects
**Check:** SQL or shell commands built with string concatenation or `fmt.Sprintf` from external input; `math/rand` used for tokens, keys, or IDs that must be unguessable; secrets or credentials hard-coded in source; `InsecureSkipVerify: true` outside tests; user-controlled file paths joined without cleaning/validation; HTML rendered with `text/template`.
**Why:** These are directly exploitable.

### RES-4 · Major · Network I/O has timeouts
**Check:** `http.DefaultClient` or `http.Get` in production paths; `http.Client{}` without `Timeout`; `http.Server` without `ReadHeaderTimeout` (and usually `ReadTimeout`/`WriteTimeout`/`IdleTimeout`); `net.Dial` without deadline.
**Why:** Without timeouts, a slow or malicious peer can hold goroutines and connections open indefinitely.

### RES-5 · Minor · Avoidable allocations in hot paths
**Check:** Only where a path is clearly performance-sensitive (request handling loops, parsing): repeated string concatenation in loops instead of `strings.Builder`; slices grown in loops when the final size is known; `time.After` inside a `select` in a loop; regexps compiled per call.
**Why:** These are cheap to fix and add up. Don't flag cold paths — premature optimization reduces clarity for nothing.

---

## API & package design (API)

### API-1 · Major · Packages are named for what they provide
**Check:** Package names like `util`, `utils`, `common`, `helpers`, `misc`, `base`, `shared`, `lib`, `types`, `models` used as grab-bags; package names that are plural, mixed-case, or contain underscores; package name differing from its directory without reason.
**Why:** A package name is part of every identifier a caller uses. Grab-bag packages accumulate unrelated code, attract dependencies from everywhere, and become the source of import cycles.

### API-2 · Major · Exported surface is minimal and deliberate
**Check:** Identifiers exported without any external user; implementation packages not under `internal/` in modules that others import; exported struct fields that callers must not mutate.
**Why:** Every exported name is a compatibility promise. Unexporting later is a breaking change; exporting later is free.

### API-3 · Major · No mutable package-level state
**Check:** Package-level `var`s that are mutated after init (registries, caches, configuration, clients, loggers); `init()` functions with side effects beyond registering with a well-known registry (e.g. `database/sql` drivers).
**Why:** Global mutable state couples every user of the package, makes tests order-dependent and non-parallelizable, and hides dependencies. Pass dependencies explicitly through constructors or parameters. Constants, sentinel errors, and compiled regexps are fine.

### API-4 · Major · Interfaces are small and defined by consumers
**Check:** Large interfaces (many methods) declared beside their only implementation; constructors returning interfaces instead of concrete types; interfaces created solely for mocking with a single implementation in non-test code; parameters of concrete types where only one or two methods are used and a small interface would decouple.
**Why:** The consumer knows which behavior it needs; declaring the interface there keeps it minimal and avoids coupling. Returning concrete types lets callers access everything without type assertions and lets you add methods without breaking anyone.

### API-5 · Major · Dependencies flow one way and `main` stays thin
**Check:** `main` packages containing business logic that can't be tested; lower-level packages importing higher-level ones (e.g. storage importing HTTP handlers); workarounds for import cycles (duplicated types, `interface{}` passing).
**Why:** Logic in `main` can't be imported or unit-tested. Inverted dependencies make every change ripple through the codebase.

### API-6 · Minor · Signatures are hard to misuse
**Check:** Functions with many parameters, especially adjacent parameters of the same type that can be swapped silently; boolean flag parameters whose meaning is unclear at the call site; growing lists of optional parameters.
**Why:** `Copy(dst, src string)` compiles just as well as `Copy(src, dst string)`. Prefer distinct types, a config struct, or functional options when parameters accumulate.

### API-7 · Minor · Zero values are useful
**Check:** Types that require a constructor where a usable zero value would be easy (e.g. lazily initializing a map); types whose zero value is silently broken without a constructor being clearly required.
**Why:** A useful zero value removes a whole class of "forgot to call New" bugs and simplifies embedding.

---

## Testing (TST)

### TST-1 · Major · Non-trivial packages have passing tests
**Check:** Packages containing logic but no `_test.go` files; `go test ./...` failures; tests that are skipped unconditionally.
**Why:** Untested logic can't be changed with confidence. Prioritize: packages with branching business logic and parsing matter more than thin wiring.

### TST-2 · Major · Tests are deterministic and isolated
**Check:** `time.Sleep` used to wait for goroutines; dependence on wall-clock time, map iteration order, or test execution order; tests sharing mutable package state; tests requiring network access or external services without a build tag or skip condition.
**Why:** Flaky tests get ignored, and an ignored test suite protects nothing.

### TST-3 · Minor · Multiple cases are table-driven with subtests
**Check:** Copy-pasted test functions differing only in inputs; loops over cases without `t.Run`, so a failure doesn't identify which case failed.
**Why:** Tables make it cheap to add cases and make coverage gaps visible; subtests give each case a name and allow running one in isolation.

### TST-4 · Minor · Failures are diagnosable
**Check:** Assertions that report only "failed" or "unexpected result" without inputs, got, and want; helpers missing `t.Helper()`; `t.Fatal` used where `t.Error` would let remaining checks run.
**Why:** A failure message should let someone fix the bug without re-running the test under a debugger.

### TST-5 · Minor · Test tooling is used idiomatically
**Check:** Manual temp-dir creation and cleanup instead of `t.TempDir()`/`t.Cleanup()`; `httptest` not used for HTTP handlers/clients; heavy mocking frameworks where a small fake would do.
**Why:** Standard tooling cleans up reliably and keeps tests readable.

---

## Readability & naming (RDB)

### RDB-1 · Minor · Names fit their scope
**Check:** Long names for short-lived locals (`indexCounter` for a loop index) or one-letter names for package-level or long-lived identifiers; names that encode the type (`userMap`, `nameString`); stutter between package and identifier (`config.ConfigLoader`); inconsistent initialisms (`Id`, `Url`, `Http` instead of `ID`, `URL`, `HTTP`); receiver names that are long, `this`/`self`, or inconsistent across methods.
**Why:** In Go the package name is always present at call sites, so it should do part of the naming work. Name length should scale with the distance between declaration and use.

### RDB-2 · Minor · Exported identifiers are documented
**Check:** Exported types, functions, constants, and packages without doc comments; doc comments that don't start with the identifier name; comments that restate the code rather than explaining intent, constraints, or non-obvious behavior; commented-out code.
**Why:** Doc comments are the API reference (`go doc`, pkg.go.dev). Comments that restate code go stale; comments that explain *why* stay useful.

### RDB-3 · Minor · Control flow stays shallow
**Check:** `else` after a branch that returns; deeply nested `if` blocks where guard clauses would keep the main path at the left margin; very long functions doing several unrelated things.
**Why:** Early returns keep the successful path readable top to bottom and make error paths obvious.

### RDB-4 · Minor · Code is consistent within the project
**Check:** Several styles for the same thing (multiple logging libraries, mixed error-wrapping styles, both `interface{}` and `any` in new code, inconsistent constructor naming).
**Why:** Inconsistency forces readers to wonder whether differences are meaningful.

---

## Build & tooling (BLD)

### BLD-1 · Critical · Module builds and passes vet
**Check:** `go build ./...` and `go vet ./...` results.
**Why:** A failing build means nothing else can be trusted. Vet findings are high-signal bugs (printf mismatches, copied locks, unreachable code, lost cancels).

### BLD-2 · Major · Module definition is clean
**Check:** Missing `go.mod`; `go mod tidy` would change `go.mod`/`go.sum`; `replace` directives pointing to local paths committed in a module others depend on; `go` directive far behind supported Go releases; dependencies with known vulnerabilities if `govulncheck` is available.
**Why:** An untidy module produces irreproducible builds and confusing dependency graphs; local `replace` directives break builds for everyone else.

### BLD-3 · Minor · Code is formatted
**Check:** Files listed by `gofmt -l`.
**Why:** Formatting is non-negotiable in the Go ecosystem; unformatted code produces noisy diffs and signals missing tooling.

### BLD-4 · Minor · Static analysis is part of the workflow
**Check:** No linter configuration (e.g. `.golangci.yml`) and no CI step running vet/tests; staticcheck or golangci-lint findings if those tools ran.
**Why:** Automated checks catch regressions before review and free reviewers to focus on design.

---

*Informed by Effective Go, Go Code Review Comments, the Go Proverbs, and Dave Cheney's Practical Go.*
