# VIII. Concurrency

Scale out via the process model.

## Rules

1. Any running program is one or more processes.
2. Processes are first-class, following the Unix model for service daemons.
3. Workloads are split into process types (web, worker, scheduler), and the app scales out by running more processes of a type. Internal threads or async are fine, but the app must also scale across processes on many machines. The set of process types and their counts is the *process formation*.
4. Processes never daemonize or write PID files. The process manager (systemd, the cloud platform, Foreman locally) handles output, crash restarts, and shutdowns.

## Violation signals

- Long background work run inside the web process (threads or in-process queues) that cannot be scaled separately.
- Forking into the background, double-fork daemonization, or writing PID files.
- Custom supervisor loops inside the app that restart child processes.
- Designs that only scale by giving a single process more CPU or memory.
- Singleton assumptions (only one instance may run) without a lock or leader election.

## Compliant patterns

- A `Procfile`, container entrypoints, or deployment definitions with one process type per workload.
- Each process runs in the foreground and exits on fatal errors, leaving restarts to the platform.
- Independent scaling per process type (separate replica counts or autoscalers).

## Legitimate exceptions

- A single process type is fine for small services, as long as it can run as several instances.
- Pre-fork servers (Gunicorn, Puma) managing their own worker pool in the foreground are compliant.
