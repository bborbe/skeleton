---
status: completed
summary: Reordered Makefile.precommit format target so gofmt -w runs last (after golines), added the vault-cli-79daf84 comment above it, and extended CHANGELOG.md Unreleased with a fix entry
execution_id: go-skeleton-fix-gofmt-exec-006-reorder-gofmt-last-in-format-target
dark-factory-version: dev
created: "2026-08-21T10:37:34Z"
queued: "2026-08-21T10:37:34Z"
started: "2026-08-21T10:49:00Z"
completed: "2026-08-21T10:50:21Z"
---

# Reorder gofmt after golines in Makefile.precommit format target

<summary>
- The canonical `make format` target runs `gofmt -w` before golines, so golines re-breaks the tree after gofmt
- Go 1.27 gofmt rejects golines' over-indented struct-literal output, failing the precommit gofmt lint check
- Every repo using this template hits that precommit failure the moment it bumps to Go 1.27
- Moving `gofmt -w` to run last normalizes the tree so the gofmt check passes
- Mirrors the proven vault-cli fix (v0.114.1, commit 79daf84)
- No runtime behavior change — formatting determinism only
</summary>

<objective>
Make the canonical `format` target gofmt-idempotent: `gofmt -w` runs last (after goimports-reviser and golines), so golines' wrapping is normalized before the golangci-lint gofmt check — no repo fails precommit on the Go 1.27 bump for formatter-ordering reasons.
</objective>

<context>
Read `Makefile.precommit` — the `format:` target currently runs, in order: go-modtool fmt → `gofmt -w` → goimports-reviser → golines (last).
The same reorder is already proven in bborbe/vault-cli (commit 79daf84) with the comment "golines last, then gofmt last so its wrapping is normalized".
</context>

<requirements>
1. In `Makefile.precommit`, edit the `format:` target so the `gofmt -w` find line runs LAST, after the golines line.
   - Current order: go-modtool fmt → gofmt -w → goimports-reviser → golines
   - New order: go-modtool fmt → goimports-reviser → golines → gofmt -w (last)
2. Add exactly this comment line immediately above the moved `gofmt -w` line (mirroring bborbe/vault-cli commit 79daf84):
   `# golines last, then gofmt last so its wrapping is normalized and the gofmt lint check passes`
3. Change nothing else in the file — no other target, no other line.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- Do NOT modify `tools.env` or any version pin
- Keep every other line of `Makefile.precommit` byte-identical
- Repo-relative paths only
</constraints>

<verification>
Run `make precommit` -- must pass.
Run `awk '/golines/{g=NR} /gofmt -w/{f=NR} END{exit !(f&gt;g)}' Makefile.precommit` -- must exit 0 (gofmt line after golines line), and confirm the comment `# golines last, then gofmt last so its wrapping is normalized and the gofmt lint check passes` sits immediately above the gofmt line.
</verification>
