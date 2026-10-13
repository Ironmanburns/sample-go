# sample-go

Public **go** sample Docker app for the local platform fleet.

- Contract: [STANDARD.md](./STANDARD.md)
- Agents: [AGENTS.md](./AGENTS.md)
- Release: Harness on `main` → `v*` (`harness/pipeline-release-main.yaml` + `scripts/harness-semver-release.sh`)
- Stable deploy: **git tag `v*` only** (`harness/pipeline-deploy-tag.yaml`)
- PR previews: `harness/pipeline-pr-preview.yaml` / teardown on close

## Local run

```bash
docker build -t sample-go:dev .
docker run --rm -p 8080:8080 sample-go:dev
curl -s localhost:8080/healthz
curl -s localhost:8080/version
```

## Versioning

Conventional Commits on PR titles → merge to `main` → Harness release pipeline
computes the next semver, bumps the language version file when applicable, pushes
an annotated `vX.Y.Z` tag (and optional GitHub Release) → deploy pipeline runs
on that tag. Merging to `main` alone does **not** ship to the stable environment.

## URLs (kind)

- Stable: `https://go-sample.burnsy.me` (NodePort `30105`)
- Preview: `https://sample-go-pr-<n>.burnsy.me` (when wired)
