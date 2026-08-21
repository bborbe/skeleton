---
status: completed
spec: [002-healthz-endpoint]
summary: 'Created healthz JSON handler returning {"status":"ok"} with Content-Type: application/json, wired through factory into main.go route, with Ginkgo/Gomega tests and CHANGELOG entry'
execution_id: 0ad976be-df7c-5d1b-874b-b49b8aae253b-exec-008-spec-002-json-healthz-endpoint
dark-factory-version: dev
created: "2026-08-21T15:35:00Z"
queued: "2026-08-21T15:39:01Z"
started: "2026-08-21T15:39:03Z"
completed: "2026-08-21T15:49:11Z"
branch: dark-factory/healthz-endpoint
---

# JSON healthz endpoint

<summary>
- `/healthz` now returns a JSON body `{"status":"ok"}` instead of the plain-text `OK`.
- The route path and HTTP method (`GET`) stay unchanged so the k8s liveness probe still targets it.
- The response carries `Content-Type: application/json` so monitoring tooling can parse the body.
- A new handler lives in `pkg/handler/`, wired through `pkg/factory/`, replacing the inline `libhttp.NewPrintHandler("OK")` in `main.go`.
- A Ginkgo v2/Gomega unit test (httptest) asserts the status code, exact body, and content type.
- No new dependencies, no manifest changes, no other routes touched.

</summary>

<objective>
Make `GET /healthz` return HTTP 200 with the exact JSON body `{"status":"ok"}` and `Content-Type: application/json`, reachable unauthenticated whenever the process is up. The route path stays `/healthz` and the k8s liveness probe in `k8s/skeleton-deploy.yaml` continues to target it unchanged.

</objective>

<context>
Read `CLAUDE.md` for project conventions. Read these files before changing anything:

- `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/main.go` — current router wiring (line ~92-103). The `/healthz` route is `router.Path("/healthz").Handler(libhttp.NewPrintHandler("OK"))`.
- `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/pkg/factory/factory.go` — existing factory functions (`CreateTestLoglevelHandler`, `CreateSentryAlertHandler`) are the pattern to follow.
- `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/pkg/handler/test-loglevel.go` and `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/pkg/handler/sentry-alert.go` — established handler signatures (`func New<Name>Handler(...) http.Handler`).
- `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/pkg/handler/test-loglevel_test.go` and `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/pkg/handler/sentry-alert_test.go` — established Ginkgo test layout.
- `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/k8s/skeleton-deploy.yaml` — the liveness probe targets `GET /healthz` on port 9090 (line ~62-70). DO NOT modify.

Reference docs (host plugin docs are read-only at `/home/node/.claude/plugins/marketplaces/coding/docs/` on the host, but inside the YOLO container they live at `/home/node/.claude/plugins/marketplaces/coding/docs/`):

- `/home/node/.claude/plugins/marketplaces/coding/docs/go-http-handler-refactoring-guide.md` — handler naming + factory wiring rules.
- `/home/node/.claude/plugins/marketplaces/coding/docs/go-testing-guide.md` — Ginkgo v2/Gomega conventions.

Critical library API verified by reading source (`/home/node/go/pkg/mod/github.com/bborbe/http@v1.26.16/`):

- `http.NewPrintHandler(format string, a ...any) http.Handler` — currently used; produces plain text via `fmt.Fprintf` with NO Content-Type. This is what we replace.
- `http.NewJSONHandler(jsonHandler JSONHandler) WithError` — uses `json.NewEncoder(resp).Encode(result)` which APPENDS A TRAILING `\n`. **Do not use this** — the spec requires the body to exactly equal `{"status":"ok"}` with no extra bytes.
- `http.ApplicationJSONContentType = "application/json"` (in `http_content-types.go`) and `http.ContentTypeHeaderName = "Content-Type"` (in `http_headers.go`).
- `http.JSONHandler` / `http.JSONHandlerFunc` interface signatures for cases that DO want the encoder newline.

So this handler must call `json.Marshal` and write the raw bytes to the response — `NewJSONHandler` is the wrong tool because of the trailing newline.

</context>

<requirements>

## 1. Create the handler file

Create `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/pkg/handler/healthz.go` with this exact content shape (license header per repo convention — copy from `pkg/handler/test-loglevel.go`):

```go
// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"net/http"

	libhttp "github.com/bborbe/http"
)

// NewHealthzHandler creates an HTTP handler that serves the canonical liveness
// response. It returns HTTP 200 with body `{"status":"ok"}` and
// Content-Type: application/json. The handler is dependency-free: it never
// touches Kafka, BoltDB, Sentry, or any other subsystem, so it can only ever
// fail by reporting "the process is down" (process exit / listener gone).
func NewHealthzHandler() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add(libhttp.ContentTypeHeaderName, libhttp.ApplicationJSONContentType)
		// Use json.Marshal + write raw bytes instead of json.NewEncoder.Encode
		// so the response body is exactly {"status":"ok"} with no trailing newline.
		body, err := json.Marshal(struct {
			Status string `json:"status"`
		}{Status: "ok"})
		if err != nil {
			// Marshal of a static struct literal cannot fail in practice; if it
			// ever does, returning an empty 500 would break the liveness probe.
			// Fall back to the exact bytes which are guaranteed valid JSON.
			body = []byte(`{"status":"ok"}`)
		}
		_, _ = resp.Write(body)
	})
}
```

Notes:
- The `_, _ = resp.Write(body)` intentionally discards errors (matches the `libhttp.WriteAndGlog` style in `pkg/handler/sentry-alert.go` line ~31). For an unauthenticated static ~15-byte write to a live listener, a failed `Write` only happens if the client hung up — nothing to recover to.
- No imports for `context`, `errors`, `libkv`, `libsentry` — the handler must remain dependency-free so the liveness probe never reflects subsystem health (spec Non-goal #1 + Failure Mode row 3).
- No `//nolint` directives.

## 2. Add the factory function

In `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/pkg/factory/factory.go`, append a new factory function following the existing pattern (no changes to the two existing entries):

```go
// CreateHealthzHandler creates an HTTP handler that serves the canonical
// `/healthz` liveness response (HTTP 200, body `{"status":"ok"}`,
// Content-Type: application/json).
func CreateHealthzHandler() http.Handler {
	return handler.NewHealthzHandler()
}
```

The existing `import ("net/http"; libsentry "github.com/bborbe/sentry"; "github.com/bborbe/go-skeleton/pkg/handler")` block already covers what this function needs — no new imports.

## 3. Wire the route in main.go

In `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/main.go`, replace the existing `/healthz` line in `createHTTPServer` (currently line ~93):

```go
// Before
router.Path("/healthz").Handler(libhttp.NewPrintHandler("OK"))

// After
router.Path("/healthz").Handler(factory.CreateHealthzHandler())
```

Leave `/readiness` alone (Non-goal #2 — that route stays as `libhttp.NewPrintHandler("OK")`). Do not change the k8s probe path (Non-goal #3). Do not touch any other route. Do not modify `application` struct fields, `Run`, or `service.Run` — this is a route-only change.

After the edit, check whether `libhttp` is still referenced elsewhere in `main.go`; if not, remove the unused import. It is still used by `libhttp.NewGarbageCollectorHandler()` and `libhttp.NewServer(...)` (lines ~100, ~105), so it stays.

## 4. Add the handler test

Create `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/pkg/handler/healthz_test.go` with this exact content shape (license header per repo convention — copy from `pkg/handler/test-loglevel_test.go`):

```go
// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/go-skeleton/pkg/handler"
)

var _ = Describe("HealthzHandler", func() {
	var httpHandler http.Handler

	BeforeEach(func() {
		httpHandler = handler.NewHealthzHandler()
	})

	It("returns HTTP 200", func() {
		req := httptest.NewRequest("GET", "/healthz", nil)
		resp := httptest.NewRecorder()

		httpHandler.ServeHTTP(resp, req)

		Expect(resp.Code).To(Equal(http.StatusOK))
	})

	It("returns the exact JSON body", func() {
		req := httptest.NewRequest("GET", "/healthz", nil)
		resp := httptest.NewRecorder()

		httpHandler.ServeHTTP(resp, req)

		// Exactly {"status":"ok"} — no surrounding whitespace, no trailing newline.
		Expect(resp.Body.String()).To(Equal(`{"status":"ok"}`))
	})

	It("sets Content-Type to application/json", func() {
		req := httptest.NewRequest("GET", "/healthz", nil)
		resp := httptest.NewRecorder()

		httpHandler.ServeHTTP(resp, req)

		Expect(resp.Header().Get("Content-Type")).To(Equal("application/json"))
	})
})
```

Three `It` blocks — one per Acceptance Criterion row (status, exact body, content type). Using `It` rather than `DescribeTable` keeps each assertion in its own reporter node, which matches the project's existing test style (`sentry-alert_test.go`, `test-loglevel_test.go`). All three blocks must pass.

## 5. Update CHANGELOG.md

In `/home/node/.cache/github-dark-factory-agent/work/0ad976be-df7c-5d1b-874b-b49b8aae253b/CHANGELOG.md`, add a new bullet under the existing `## Unreleased` section (lines ~19-23). Follow the existing `- feat:` / `- fix:` / `- chore:` prefix style. The new bullet:

```
- feat: serve `/healthz` as JSON (`{"status":"ok"}` with `Content-Type: application/json`) so consumers and monitoring tooling can parse the liveness response; route path and k8s liveness probe unchanged
```

Do not bump the version header. Do not add a new `##` version section.

</requirements>

<constraints>

- The `/healthz` route path, HTTP method (`GET`), and k8s liveness-probe role MUST NOT change (spec Acceptance Criteria #4 + Non-goal #3 + Failure Mode row 4).
- The handler MUST remain dependency-free — no DB, Kafka, Sentry, or other subsystem calls (spec Non-goal #1 + Failure Mode row 3). A future change that wired a dependency in would silently turn the probe into a readiness check.
- Response body MUST be exactly `{"status":"ok"}` — NO trailing newline, NO surrounding whitespace, NO extra fields. This is why `libhttp.NewJSONHandler` is the wrong tool (it uses `json.NewEncoder.Encode` which appends `\n`); the handler MUST use `json.Marshal` + raw `resp.Write`.
- Response `Content-Type` MUST be exactly `application/json`. The spec accepts headers that "start with" `application/json` per Acceptance Criterion #3; the constant `libhttp.ApplicationJSONContentType` evaluates to that exact string.
- Do NOT modify `/readiness`, `/metrics`, `/resetdb`, `/resetbucket/{BucketName}`, `/setloglevel/{level}`, `/gc`, `/testloglevel`, `/sentryalert`, the `application` struct, `Run`, `service.Run`, `libkafka`, `libboltkv`, `libmetrics`, `libsentry`, `libhttp.NewServer`, or anything in `k8s/`. The PR is scoped to the healthz response.
- No new dependencies — do not edit `go.mod` or `go.sum`.
- Handler naming follows the repo's frozen convention: `func New<Name>Handler(...) http.Handler` in `pkg/handler/`, wired by `func Create<Name>Handler(...) http.Handler` in `pkg/factory/`.
- Tests use Ginkgo v2 / Gomega with `httptest`. No `//nolint` without explanation. Errors wrapped via `github.com/bborbe/errors` if any are wrapped (the healthz handler has no fallible call paths except the static `json.Marshal`, which is handled inline).
- `make precommit` MUST pass (format + generate + lint + test + security checks).
- Existing tests MUST still pass.
- Do NOT commit — dark-factory handles git.

</constraints>

<verification>

Run these commands in order. Each must succeed before moving to the next:

1. `go build ./...` — confirms `main.go` still compiles after the route wiring change and that `libhttp` import is still needed (or removed if not).
2. `go test -mod=mod ./pkg/handler/` — must pass and include the three new `HealthzHandler` `It` blocks (run with `-v` if you need to see the Ginkgo reporter output).
3. `go test -mod=mod ./pkg/factory/` — must pass; existing factory test suite unchanged but confirms the package still compiles after the new factory function lands.
4. `make precommit` — full format + generate + lint + test + security pipeline. Must exit 0.
5. Manually inspect the final response shape with a quick httptest smoke in a scratch file (or via `go test -run TestSuite -v` with a focused `--focus="HealthzHandler"` Ginkgo run): confirm `resp.Body.String()` is the 15 bytes `{"status":"ok"}` and `resp.Header().Get("Content-Type")` is `application/json`.
6. Confirm the k8s probe path is unchanged: `grep -n "/healthz" k8s/skeleton-deploy.yaml` must still show the liveness probe pointing at `/healthz`.

If `go test -mod=mod ./pkg/handler/` fails on the exact-body assertion (Acceptance Criterion #2), the most likely cause is the `json.NewEncoder` newline issue described in `<context>` — re-check that the handler uses `json.Marshal` + `resp.Write`, NOT `libhttp.NewJSONHandler`.

If `make precommit` fails on lint with a Content-Type or header-name complaint, the constant form `libhttp.ContentTypeHeaderName` (`"Content-Type"`) is the canonical way to reference it; do not switch to a raw string literal unless the linter specifically rejects the constant (it doesn't — `sentry-alert.go` and `test-loglevel.go` use libhttp constants without issue).

</verification>
