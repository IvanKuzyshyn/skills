---
name: go-quality-review
description: Evaluate the quality of a Go (golang) project or package and produce a scored report with severity-ranked findings. Use this whenever the user asks to review, audit, assess, grade, or score Go code, asks "how good is this Go codebase", wants a code-health or maintainability check of a Go repo, is preparing a Go project for production or open-sourcing, or asks what to improve in their Go service, library, or CLI — even if they don't say "quality review" explicitly.
---

# Go Quality Review

Assess a Go codebase against a catalog of requirements covering correctness, error handling, concurrency, resource safety, API design, readability, testing, and build hygiene. Produce a report with a 0–100 score and actionable findings.

The aim is a review a senior Go engineer would sign off on: findings that are real, located precisely, ranked by how much they matter, and paired with a concrete fix. A short report of true problems is worth more than a long one padded with style nits.

## Workflow

### 1. Scope the review

- Confirm the target: whole module, a subset of packages, or a diff. Default to the module rooted at the current directory (look for `go.mod`).
- Exclude code the project does not own: `vendor/`, `third_party/`, and generated files (header line matching `// Code generated .* DO NOT EDIT.`).
- Map the layout: list packages (`go list ./...`), find `main` packages, note `internal/` usage, count non-test vs test files. For large modules (>~50 packages), prioritize core domain packages and public API, and state in the report which packages were sampled.

### 2. Run automated checks

Run the bundled script from the module root:

```bash
bash <skill-dir>/scripts/run_checks.sh <module-root>
```

It runs, where available: `go build`, `go vet`, `gofmt -l`, `go mod tidy -diff`, `go test -race` (falling back to no race detector when cgo is unavailable), and `staticcheck` / `golangci-lint` if installed. Missing optional tools are reported, not failed.

`go test` executes project code. If the project is not the user's own or its origin is unclear, ask before running tests (pass `--no-test` to skip them).

Treat tool output as evidence, not as the review. Tools catch a narrow class of issues; most of the requirements below need reading the code.

### 3. Review the code against the requirements

Read `references/requirements.md`. It lists every requirement with an ID, severity, what to look for, and why it matters.

Work area by area. For each requirement decide: **met**, **violated** (with evidence), or **N/A**. Useful searches to seed the review:

- Discarded errors: `_ = `, `_, _ =`, calls to `Close`/`Write`/`Encode` with no result use
- Goroutines: `go func`, `go ` + identifier — then trace how each one exits
- Package-level state: `^var ` at file scope, `func init()`
- Panics and exits: `panic(`, `log.Fatal`, `os.Exit` outside `main`
- Grab-bag packages: directories named `util`, `utils`, `common`, `helpers`, `misc`, `base`, `shared`
- Resource handling: `.Body`, `Open(`, `sql.Rows`, `time.After` inside loops, `defer` inside loops

Rules for findings:

- **Verify before reporting.** Open the file and confirm the issue in context. A `_ = f.Close()` on a read-only file is fine; on a written file it can lose data. Context changes the verdict.
- **Cite precisely**: `path/to/file.go:LINE`. If a problem is widespread, cite 2–3 representative locations and give a count.
- **Group repeated instances** of the same requirement into one finding rather than listing each occurrence.
- **Give the fix**, ideally as a short code snippet. Explain the consequence in concrete terms (what breaks, for whom, when).
- **Credit what is good.** Note areas that are clearly well done; this calibrates the score and tells the user what to preserve.

### 4. Score

Score each area from 10 and deduct per *violated requirement* (not per occurrence):

| Severity | Deduction | Widespread (≥ 3 packages or pervasive) |
|---|---|---|
| Critical | 4 | 6 |
| Major | 2 | 3 |
| Minor | 0.5 | 1 |

Floor each area at 0. Then compute the weighted total (0–100):

| Area | Weight |
|---|---|
| Error handling (ERR) | 20 |
| Concurrency (CON) | 15 |
| Resource & runtime safety (RES) | 15 |
| API & package design (API) | 15 |
| Testing (TST) | 15 |
| Readability & naming (RDB) | 10 |
| Build & tooling (BLD) | 10 |

`total = Σ(area_score × weight) / 10`

If an area is entirely N/A (e.g. no concurrency at all), drop it and rescale the remaining weights to sum to 100.

Grade: **A** ≥ 90, **B** ≥ 80, **C** ≥ 70, **D** ≥ 60, **F** < 60. Any open Critical finding caps the grade at **C** — a project that can lose data or deadlock is not a B project regardless of polish. A failing build caps the grade at **F**.

### 5. Write the report

Use this structure:

```markdown
# Go Quality Review: <module path>

**Score: <NN>/100 — Grade <X>**
Scope: <packages reviewed, what was excluded or sampled>
Go version: <from go.mod> · Tools run: <list> · Tools unavailable: <list>

## Summary
<3–5 sentences: overall state, the most important problems, the biggest strengths.>

## Scores by area
| Area | Score | Key issue |
|---|---|---|
| Error handling | 6/10 | Errors logged and returned in service layer |
| ... | | |

## Findings

### Critical
#### [CON-2] Goroutine leak in request fan-out
**Where:** `internal/fetch/fetch.go:42` (and 2 similar sites)
**Problem:** <what is wrong and the concrete consequence>
**Fix:** <how to fix, with a short snippet>

### Major
...

### Minor
...

## Strengths
- <specific things done well, with locations>

## Recommended next steps
1. <highest-impact action first>
2. ...
```

Omit empty severity sections. Keep the summary honest: if the code is good, say so and keep the findings list short.
