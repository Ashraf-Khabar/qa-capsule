---
icon: material/console
---

# Utility binaries

Besides the main server (`cmd/srecapsule`) and developer CLI (`cmd/cli`), the repository ships small **operator tools**.

---

## Server (`srecapsule`)

```bash
go run ./cmd/srecapsule/main.go
# or Docker: docker compose up -d --build
```

| Env | Purpose |
|---|---|
| `SRECAPSULE_DATA_DIR` | SQLite + artifacts directory |
| `SRECAPSULE_JWT_SECRET` | JWT signing |
| `APP_ENV=development` | Dev JWT fallback (not for production) |
| `SRECAPSULE_INGEST_WORKERS` | Async ingest worker count |
| `SRECAPSULE_MCP_TOKEN` | Protect `/mcp` |

---

## Developer CLI (`srecapsule` / `cmd/cli`)

```bash
go build -o bin/srecapsule ./cmd/cli
bin/srecapsule run [flags] -- <your test command>
```

Documented in [Artifacts & CLI](../guides/artifacts-and-cli.md).

---

## JUnit agent (`cmd/agent`)

Uploads a JUnit XML file to the webhook upload endpoint.

```bash
go run ./cmd/agent/main.go -h
```

Typical CI usage: produce `results.xml`, then agent POSTs to `POST /api/webhooks/upload` with project API key.

See [JUnit XML upload](../integration/junit-xml-upload.md).

---

## listusers (`cmd/listusers`)

Lists users from the SQLite database (offline admin tool).

```bash
go run ./cmd/listusers/main.go
go run ./cmd/listusers/main.go /path/to/srecapsule.db
```

Respects `SRECAPSULE_DATA_DIR` and project root detection via `core.EnsureProjectRoot()`.

---

## resetpass (`cmd/resetpass`)

Resets a user password when locked out of the UI.

```bash
go run ./cmd/resetpass/main.go admin 'NewSecurePass123!'
go run ./cmd/resetpass/main.go --no-change-required admin 'TempPass1!'
```

Flags:

| Flag | Effect |
|---|---|
| `--no-change-required` / `-n` | Skip forced password change on next login |

---

## Security notes

- Utility commands read **`srecapsule.db` directly** — run only on trusted hosts.
- Do not expose the database file in shared backups without encryption.
- Prefer UI admin reset when SMTP and audit trail matter.
