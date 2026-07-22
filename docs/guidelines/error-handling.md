# Error Handling

## Why this document exists

This document defines **one** set of rules for error handling so that it is
uniform and, above all, **easy to investigate when something breaks in
production**.

Optimizing for debuggability is the guiding principle here. Every rule below is
justified by the same question: *when an incident happens, can an engineer get
from a symptom to the root cause quickly?*

These rules are **binding** (MUST, unless stated otherwise). All code must follow
them.

## Rules at a glance

1. **One library.** Use `github.com/friendsofgo/errors` everywhere. It preserves
   stack traces. Do not create errors with the standard `errors.New` / `fmt.Errorf`.
2. **Wrap with context, don't log while propagating.** As an error travels up,
   add context (IDs, state) to the *message*. Never log an error you are also
   returning.
3. **Log once, where the error stops.** The single place that handles an error
   terminally (an interceptor, a reconcile loop tick) logs it once — with the full
   wrap chain, stack trace, `trace_id`, and `user_id`.
4. **`Error` means "a human must look."** Reserve the `Error` level for the
   unexpected. Expected/transient failures are `Warn`; normal events are `Info`.
5. **Classify at the usecase layer.** The usecase layer tags errors with a
   business meaning (bad request, not found, ...). The transport boundary maps the
   tag to a Connect code; anything untagged is `Internal`.
6. **Don't leak internals.** `Internal` errors return only a generic message to the
   client. The detail lives in the logs.
7. **Loops log-and-continue.** Reconcile loops never add their own retry/compensation.
8. **`panic` is for programmer errors and startup only.**

The rest of this document explains each rule and shows the intended shape of the code.

## 1. Error library and stack traces

Use `github.com/friendsofgo/errors` (a `pkg/errors` fork) as the **only** error
constructor library. It attaches a stack trace at the point the error is created or
wrapped, which is the single most valuable thing for locating where a failure
originated.

- **MUST** create and wrap errors with `errors.New`, `errors.Wrap`, and
  `errors.Wrapf` from this package.
- **MUST NOT** use the standard library's `errors.New` or `fmt.Errorf` to create or
  wrap errors — they discard the stack trace and are the reason the two libraries
  currently mix.
- Inspecting errors with the standard `errors.Is` / `errors.As` is fine;
  `friendsofgo/errors` exposes the same functions and chains interoperate.

**Capture the stack once, at the origin.** Wrap an external/library error the moment
it enters our code so the stack points at the boundary where it happened. Higher up,
keep adding *context* with `Wrap`, but do not re-wrap an error without adding new
information — that only adds noise to the chain.

```go
// Good: the DB error is wrapped as it enters our code; the stack is captured here.
row, err := models.Applications(...).One(ctx, db)
if err != nil {
    return errors.Wrap(err, "querying application")
}
```

## 2. Creating and wrapping errors

As an error propagates, attach the **local context** that will matter during an
investigation — identifiers and relevant state — onto the wrap message, not into a
separate log line. One log at the boundary then carries everything.

```go
// Good: app_id and step are on the message; no local log needed.
if err := b.updateApp(ctx, appID); err != nil {
    return errors.Wrapf(err, "updating app (app_id=%s, step=%s)", appID, step)
}
```

### Message style

Wrap and error messages **MUST**:

- be written as a **gerund phrase** describing the operation: `"syncing deployments"`,
  `"updating app"`, `"querying repository"`.
- start with a lowercase letter and have **no** trailing punctuation.
- **not** include `"failed to"`. Wrapping already means "this failed"; adding it to
  every layer produces `failed to sync: failed to update: failed to query`, whereas
  gerunds read as a clean path: `syncing deployments: updating app: querying db`.

This style applies to **wrap and internal error messages**. It does **not** apply to
the classified, client-facing messages produced at the usecase layer (rule 5) or the
generic message returned to the client (rule 6): those describe a *condition*, not an
operation, and are written as plain phrases — `"application not found"`,
`"internal server error"` — never gerunds.

```go
// Good
return errors.Wrap(err, "reading docs root")
return errors.Wrapf(err, "resolving ref (ref=%s, app_id=%s)", ref, appID)

// Bad
return errors.Wrap(err, "Failed to read docs root.")   // "failed to", capitalized, punctuation
return fmt.Errorf("read docs root: %w", err)            // stdlib, no stack trace
```

## 3. Where to log: log once, at the point the error stops

The biggest obstacle to investigation is the **same failure logged in several
places** with partial context. To prevent it:

- **MUST** log an error only at the point where it is *terminally handled* — i.e.
  where you stop returning it to the caller (you swallow it, or you translate it into
  a response). In practice these points are:
  - the gRPC/Connect **interceptor** (`pkg/infrastructure/grpc/log_interceptor.go`),
    which logs every request error once, and
  - the **top of a reconcile loop tick**, which logs a failure and moves on.
- **MUST NOT** log an error that you are also returning. If you `return err`, the
  caller (ultimately the boundary) owns logging it.

Because context is carried on the wrap message (rule 2), the single boundary log
still contains every ID and every step — plus the stack trace, `trace_id`, and
`user_id`. Nothing is lost by not logging locally.

```go
// Bad: logs here AND returns → the boundary logs it again (double log).
if err != nil {
    slog.ErrorContext(ctx, "updating app", "app_id", appID, "error", err)
    return err
}

// Good: attach context, return; the boundary logs it exactly once.
if err != nil {
    return errors.Wrapf(err, "updating app (app_id=%s)", appID)
}
```

## 4. Log levels

Levels carry meaning. Keeping `Error` pure is what makes it a usable alert signal.
The interceptor already derives the level from the Connect code; that mapping is the
project-wide standard.

- **`Error`** — unexpected; a human should look. Internal failures, bugs, broken
  external dependencies. Connect codes `Internal`, `Unknown`, `Unavailable`,
  `DeadlineExceeded`, `Unimplemented`, `DataLoss`.
- **`Warn`** — expected but noteworthy. Client errors (`InvalidArgument`,
  `NotFound`, `PermissionDenied`, ...) and transient failures a loop will retry.
- **`Info`** — normal events: a successful request, a loop making normal progress.

A per-item failure inside a loop that the next tick will retry is `Warn`, not `Error`
— it is expected and self-healing.

## 5. Classifying errors and mapping them to transport codes

Deciding that a failure is "not found" or "bad request" requires knowing the business
intent, so **classification is the usecase layer's job**.

- **Infrastructure and domain layers** return plain wrapped errors. They do not know
  or decide HTTP/Connect semantics.
- **The usecase layer** tags an error with its business meaning using the helper in
  `pkg/usecase/apiserver/errors.go` (`newError(ErrorType..., "message", cause)`).
  Use `errors.Is` to detect a known low-level condition and translate it into a tag.
- **The transport boundary** (`handleUseCaseError` in
  `pkg/infrastructure/grpc/api_service.go`) decomposes the tag and maps it to a
  Connect code. **Anything untagged maps to `Internal`** — an untagged error is by
  definition something we did not anticipate.

```go
// usecase layer: detect a known condition and classify it.
app, err := s.appRepo.GetApp(ctx, id)
if errors.Is(err, repository.ErrNotFound) {
    return nil, newError(ErrorTypeNotFound, "application not found", err)
}
if err != nil {
    return nil, errors.Wrap(err, "getting application") // stays Internal at the boundary
}
```

Add new `ErrorType` values when a genuinely new category of client-facing outcome
appears — keep the set small and meaningful.

## 6. What the client sees

Split what the client receives from what we record:

- **Classified errors** (bad request, not found, forbidden, already exists) **MUST**
  return their message to the client — the user can act on it.
- **`Internal` errors MUST return only a generic message** (e.g. "internal server
  error"). The wrapped detail — DB errors, file paths, internal state — stays in the
  logs only, never in the response. The full detail is already captured by the single
  boundary log (rule 3), correlated by `user_id`, `procedure`, and timestamp.

## 7. Reconcile loops

NeoShowcase is built on reconciliation loops (see `docs/architecture.md`). Errors in a
loop are transient by design: the next tick reads desired state from the DB and tries
again.

- **MUST** log a failure and continue — do **not** abort the whole system on one
  error.
- When processing many items in one tick, a single item's failure **MUST NOT** abort
  the remaining items. Log that item (`Warn` if the next tick will retry) and keep
  going.
- **MUST NOT** add bespoke retry, backoff, or compensation logic. The loop *is* the
  retry mechanism; adding more only hides problems and duplicates machinery.

```go
for _, app := range apps {
    if err := s.reconcileApp(ctx, app); err != nil {
        // log-and-continue: the next tick retries this app.
        slog.WarnContext(ctx, "reconciling app", "app_id", app.ID, "error", err)
        continue
    }
}
```

## 8. panic

- **MUST** return an `error` for anything that can happen at runtime — request
  handling, loop bodies, I/O.
- `panic` is allowed **only** for:
  - **programmer errors** — an invariant that cannot be violated unless the code is
    wrong (e.g. an impossible `switch` branch), and
  - **unrecoverable startup/initialization failures**, where the process cannot
    sensibly run (config that fails to load, a required dependency that fails to
    construct).
- **MUST NOT** `panic` inside a goroutine for control flow — an unrecovered goroutine
  panic takes the whole process down.

## 9. Structured logging

- **MUST** use `slog` with the `Context` variants (`slog.ErrorContext`,
  `slog.WarnContext`, `slog.InfoContext`) so that `trace_id` and `user_id`, injected
  by the interceptor, travel with every line.
- **MUST** pass the error as a structured attribute: `slog.ErrorContext(ctx, "…",
  "error", err)`. Do not format the error into the message string.
- Prefer stable, queryable attribute keys (`app_id`, `build_id`, `procedure`) over
  interpolated prose.
