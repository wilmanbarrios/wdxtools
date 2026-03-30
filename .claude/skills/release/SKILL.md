---
name: release
description: Determine semver bump from conventional commits since the last tag, create the version tag, and push branch + tag to trigger the GoReleaser pipeline.
disable-model-invocation: true
---

## Release wdxtools

Automate the full release flow: determine the next semver version from conventional commits, tag, and push.

### Step 1 — Preflight checks

Run these in parallel:

1. `git status` — abort if the working tree is dirty (unstaged or uncommitted changes).
2. `git tag -l --sort=-v:refname | head -1` — get the latest semver tag (e.g. `v0.1.0`).
3. `git remote -v` — confirm an `origin` remote exists.

If the tree is dirty, tell the user to commit or stash first and **stop**.

### Step 2 — Analyze commits for semver bump

Get all commits since the last tag:

```bash
git log <last-tag>..HEAD --oneline
```

Classify each commit into **releasable** (affects binaries/libraries that users consume) or **non-releasable** (docs, CI, benchmarks, tooling only):

| Commit prefix | Releasable? | Bump |
|---|---|---|
| `feat:` or `feat(…):` | Yes | **minor** |
| `fix:` or `fix(…):` | Yes | **patch** |
| `perf:` or `perf(…):` | Yes | **patch** |
| `refactor:` | Yes | **patch** |
| `build:` | Yes | **patch** |
| `docs:`, `style:`, `test:`, `ci:`, `chore:` | **No** | — |
| Any commit whose body contains `BREAKING CHANGE:` or whose prefix ends with `!` (e.g. `feat!:`) | Yes | **major** |

Use the **highest** bump found across **releasable** commits only. If no commits exist since the last tag, inform the user there is nothing to release and **stop**.

#### No releasable commits — push only

If there are commits since the last tag but **none are releasable** (all are `docs:`, `test:`, `ci:`, `chore:`, `style:`), do NOT proceed with a release. Instead:

1. Tell the user: "No client-facing changes since `<last-tag>`. These commits don't affect the published binaries:" and list them.
2. Use `AskUserQuestion` with options:
   - **Push to origin (Recommended)** — push the branch without tagging
   - **Release anyway** — proceed with a patch release (override)
   - **Abort** — cancel
3. If the user chooses "Push to origin", run `git push origin <current-branch>` and stop.
4. If the user chooses "Release anyway", continue to Step 3 with a **patch** bump.

### Step 3 — Compute the next version

Parse the latest tag (e.g. `v0.1.0`) into major.minor.patch, apply the bump, and reset lower components:

- **major** → `major+1.0.0`
- **minor** → `major.minor+1.0`
- **patch** → `major.minor.patch+1`

Prefix with `v` (e.g. `v0.2.0`).

### Step 4 — Confirm with the user

Use `AskUserQuestion` to confirm. Show:

- The commits that will be included
- The determined bump type (major/minor/patch)
- The new version tag

Options:
- **Release vX.Y.Z (Recommended)** — proceed with the computed version
- **Override version** — let the user specify a different version
- **Abort** — cancel

### Step 5 — Tag and push

Run sequentially:

```bash
git tag -a <new-version> -m "<new-version>"
git push origin <current-branch>
git push origin <new-version>
```

Always create **annotated** tags (`-a` flag) to match the existing tag convention in this repo.

### Step 6 — Confirm

Tell the user:
- The tag that was created and pushed
- That GitHub Actions will now run the GoReleaser pipeline
- Provide the URL: `https://github.com/wilmanbarrios/wdxtools/actions`
