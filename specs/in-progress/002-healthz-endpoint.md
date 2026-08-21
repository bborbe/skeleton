---
status: prompted
approved: "2026-08-21T15:30:15Z"
generating: "2026-08-21T15:33:35Z"
prompted: "2026-08-21T15:36:56Z"
branch: dark-factory/healthz-endpoint
---

## Summary

- Add a liveness endpoint `/healthz` that answers GET with HTTP 200 and a JSON body exactly matching `{"status":"ok"}`.
- The existing `/healthz` route (currently a plain-text `OK` body) is replaced by this JSON response; the route path and its role as the k8s liveness probe target are unchanged.
- The endpoint is unauthenticated, static, and dependency-free: it reports process liveness only, never the health of Kafka, BoltDB, or Sentry.
- Ships as a pure code change via PR — no k8s manifest changes, no config changes, no new dependencies.

## Problem

go-skeleton already serves a `/healthz` route used as the k8s liveness probe, but the response body is a plain-text `OK`. Consumers and monitoring tooling that parse health responses get no structured payload, and the skeleton — a template other services copy — models a non-standard liveness contract. A canonical JSON `{"status":"ok"}` response makes the liveness contract explicit and machine-checkable.

## Goal

`GET /healthz` returns HTTP 200 with a JSON body exactly `{"status":"ok"}` and `Content-Type: application/json`, reachable without authentication whenever the process is up. The route path, the k8s liveness probe targeting it, and all other routes are unchanged.

## Non-goals

- No readiness-style dependency checks (Kafka, BoltDB, Sentry) in the response.
- No changes to `/readiness`, `/metrics`, or any other route.
- No changes to k8s manifests or probe configuration.
- No authentication, rate limiting, or access logging on `/healthz`.
- No body fields beyond `status`.
- A separate JSON path (e.g. `/health.json`) or keeping the plain-text body alongside a JSON one was considered and rejected — the contract is one canonical JSON response on the existing probe path.

## Acceptance Criteria

- [ ] `GET /healthz` returns HTTP 200 — evidence: HTTP response code is 200.
- [ ] Response body exactly matches `{"status":"ok"}` — evidence: HTTP response body equals `{"status":"ok"}` (no extra fields, no surrounding whitespace).
- [ ] Response `Content-Type` header starts with `application/json` — evidence: HTTP response header matches `application/json`.
- [ ] No change to k8s probe manifests — evidence: `git diff k8s/` is empty after the change.

## Verification

## Container-executable (runs inside the YOLO container at prompt time)

- `make precommit` exits 0 — format + generate + lint + test + security checks.
- `make test` exits 0 — full suite including the new handler test.
- `go test -mod=mod ./pkg/handler/` exits 0 — targeted handler suite asserting the `/healthz` response (status 200, exact body, content type).

## Desired Behavior

1. `GET /healthz` returns HTTP status 200 whenever the process is up.
2. The response body is the exact JSON document `{"status":"ok"}`.
3. The response `Content-Type` header is `application/json`.
4. The endpoint requires no auth, no headers, no query parameters, and no request body.
5. The route is served at exactly `/healthz` and remains the k8s liveness probe target (probe config untouched).
6. The handler makes no backend calls (DB, Kafka, Sentry) — liveness reports process liveness only.
7. A Ginkgo v2/Gomega handler test (httptest) asserts status 200, the exact body, and the content type.

## Constraints

- The `/healthz` route path and its liveness-probe role must not change.
- Handler code follows the repo's frozen conventions: HTTP handlers in `pkg/handler/` as `New<Name>Handler() http.Handler`; DI wiring via factory functions in `pkg/factory/` as `Create<Name>Handler()` (bcom-style service framework).
- Tests use Ginkgo v2/Gomega with `httptest` (Counterfeiter for mocks); errors wrapped via `github.com/bborbe/errors`; no `//nolint` without explanation.
- Assumption: the k8s liveness probe checks HTTP status only, so changing the response body does not affect probe behavior.
- `make precommit` must pass; no new dependencies.

## Failure Modes

| Trigger | Expected behavior | Recovery |
|---|---|---|
| Handler returns non-200 or wrong body after change | Targeted handler test fails; liveness probe would flap in future deploys | Fix handler to return 200 + exact body; re-run `go test -mod=mod ./pkg/handler/` until green |
| Content-Type missing or wrong | Handler test asserting the header fails | Set `application/json` on the response; re-run the targeted test |
| Handler wired to backend deps (DB/Kafka) | Liveness starts reflecting subsystem health; test with deps absent fails | Keep handler dependency-free (static response); verify via code review + handler test |
| Route path or method changed (e.g. GET→POST) | k8s liveness probe (GET) stops reaching the route while the manifest is unchanged | Restore `GET /healthz`; confirm the probe path in the manifest is still `/healthz` |
| Out-of-scope files changed (other routes, manifests) | `git diff k8s/` non-empty; unrelated route tests fail | Revert out-of-scope changes; keep the PR scoped to the healthz change |

## Security / Abuse

`/healthz` is an unauthenticated public endpoint serving a static ~15-byte response. It accepts no user input, performs no state mutation, touches no filesystem or network, and reflects no attacker-controlled data, so the abuse surface is minimal: no injection, no amplification (tiny static body), no sensitive-data exposure (the body is fixed and contains no versions, config, or environment information). No auth or rate limiting is warranted — the endpoint is the liveness probe path and must stay reachable by the unauthenticated k8s probe.

## Suggested Decomposition

Not needed. This is a single-layer, single-behavior spec (one endpoint response + handler + test); a single prompt covers it.

## Do-Nothing Option

If not done, `/healthz` keeps returning plain-text `OK`. The liveness probe still works (status 200), so the operational cost is low — but the skeleton keeps modeling a non-JSON health contract that downstream parsers cannot rely on, and every service cloned from it inherits that example. Doing it now is a small change; deferring leaves the wrong contract in the template indefinitely.
