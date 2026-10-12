---
name: pr
description: Use when drafting, creating, or updating a pull request in the Ngadu repository, including requests to prepare a branch for review or standardize a PR title and description.
---

# Pull requests

Use this repository's PR format for a reviewer who has not seen the chat. Write in English by default; follow an explicit user language request.

## Title

Use Conventional Commits style: `type(scope): description`. Scope is optional. Use a short imperative description without a trailing period.

| Type | Purpose |
| --- | --- |
| `feat` | Application feature |
| `fix` | Bug fix |
| `docs` | Documentation or agent guidance |
| `chore` | Repository maintenance or tooling |
| `refactor` | Restructuring without a behavior change |
| `test` | Test additions or corrections |
| `build` | Build or dependency configuration |
| `perf` | Performance improvement |

Choose the type for the overall result. Use scopes such as `frontend`, `backend`, or `docs` only when they clarify the change. Mark an actual breaking contract change with `!` before the colon and describe its impact in Notes. This is a local PR title convention, not a GitHub requirement.

## Description

Every PR body has these sections, in this order:

```markdown
## Summary

Explain the concrete problem and resulting behavior in one or two sentences.

## Validation

- `command actually run` — passed, failed, or blocked, with relevant details.
- Checks not run — identify them and explain why.
```

Use real results in place of the example text. Include the not-run bullet only for checks known to have been skipped. For documentation-only changes, explain that application checks are unnecessary; if they were skipped, state that reason. When execution or results are unknown, say they were not provided rather than inferring that a check ran or was skipped.

Add these sections after Validation when applicable:

- `## Notes`: material limitations, diagnostics despite a zero exit status, breaking changes, migrations, setup steps, or relevant tradeoffs. Include screenshots for visible UI changes when available; state when visual verification was unavailable.
- `## Related issues`: use `Closes #123` only when the PR fully resolves that issue; use `Related to #123` otherwise. Include only verified, existing issues. An issue is not required.

Keep small PRs short. Describe the final diff and its purpose; update the title and body together when scope changes. Summarize meaningful changes instead of enumerating every file or recounting the chat.

## Prepare and publish

1. Read the applicable `AGENTS.md` and inspect the complete diff against the intended base branch. Confirm the base from the remote default or the user's request. Keep unrelated work out of the PR.
2. Follow the checks in the root `AGENTS.md` for the changed files and review the diff for correctness, unintended content, and secrets. Report diagnostics and missing checks accurately; a zero exit code does not erase diagnostic output. Local hooks check the working tree, so inspect partially staged changes separately.
3. A request to draft text authorizes drafting. Push branches or create/update GitHub PRs only when the user has authorized that action in the current task. Reuse the authorized task's existing open PR when appropriate. This skill itself grants no publishing, merge, or messaging permission.
4. When publishing is authorized, use an available GitHub tool or `gh`. For multiline bodies with `gh`, write the exact Markdown to a temporary file and pass `--body-file`; preserve actual newlines.
5. Verify the PR URL, base/head branches, and final diff. When available, attach every created PR to the chat using `attach_artifact` with artifact type `pull_request`. Return the PR link and any material validation limitation. Leave merging to an explicit user request.

## Example

Title: `chore: add shared Git pre-commit checks`

```markdown
## Summary

Run lint and typecheck before commits affecting the frontend or Git hooks. Contributors enable the shared hook once per clone with `bun run setup:hooks`; Markdown-only commits skip application checks.

## Validation

- `sh -n .githooks/pre-commit` — passed.
- Temporary Git repository scenarios — documentation-only changes skip checks; failing lint or typecheck blocks the hook.
- `bun run lint` and `bun run typecheck` — passed.
- `bun run build` — exited successfully with a prerender diagnostic described below.

## Notes

The existing `/matches/[id]` route reports `useParams()` outside Suspense during prerendering.
```

The example illustrates formatting; substitute evidence from the current change rather than copying its results.
