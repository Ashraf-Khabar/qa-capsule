---
icon: material/rocket-launch
---

# Documentation Deployment Pipeline

This site is built with [MkDocs Material](https://squidfunk.github.io/mkdocs-material/) and deployed automatically to **GitHub Pages** via GitHub Actions.

---

## Workflow File

`.github/workflows/docs.yml`

### Triggers

| Event | When |
|---|---|
| `push` to `main` | Changes under `docs/**`, `mkdocs.yml`, or the workflow file |
| `workflow_dispatch` | Manual run from the Actions tab |

### Jobs

1. **build-and-deploy** — Installs Python 3.12, runs `mkdocs build --strict`, pushes the generated `site/` folder to the **`gh-pages`** branch.

---

## One-Time Repository Setup

Point GitHub Pages at the branch that contains the built MkDocs HTML (not `main`, which only has Markdown sources and the README).

1. Open your repository on GitHub.
2. Go to **Settings → Pages**.
3. Under **Build and deployment → Source**, select **Deploy from a branch**.
4. Set **Branch** to **`gh-pages`** and folder **`/ (root)`**.
5. Save.
6. Run **Actions → Deploy Documentation → Run workflow** once (or push a change under `docs/`).

!!! warning "Do not publish from `main`"
    If Pages uses the **`main`** branch, visitors see the repository **README**, not the MkDocs site. The workflow builds HTML into `site/` locally in CI and publishes it only on **`gh-pages`**.

!!! note "MkDocs, not pydoc"
    This project uses **[MkDocs Material](https://squidfunk.github.io/mkdocs-material/)** (`mkdocs build`). Python's built-in **pydoc** is not used here.

---

## Local Development

```bash
# Install dependencies
pip install -r docs/requirements.txt

# Serve with live reload
mkdocs serve

# Open http://127.0.0.1:8000
```

### Strict build (same as CI)

```bash
mkdocs build --strict
```

Fix any broken links or missing files before pushing to `main`.

### Mermaid diagrams show as raw code

MkDocs Material only renders ` ```mermaid ` blocks when `mkdocs.yml` configures `pymdownx.superfences` with a `mermaid` custom fence (`class: mermaid`). This repo includes that setting. After changing `mkdocs.yml`, restart `mkdocs serve` and hard-refresh the browser (Ctrl+F5). Requires `mkdocs-material>=9.5`.

---

## Project Structure

```
docs/
├── index.md                    # Home page
├── setup/                      # Installation & config
├── integration/                # CI/CD guides
├── guides/                     # Dashboard & operations
├── plugins/                    # Plugin engine
├── api/                        # REST API reference
├── deployment/                 # This file
└── requirements.txt            # Python dependencies

mkdocs.yml                      # Site configuration & navigation
.github/workflows/docs.yml      # CI/CD deploy workflow
```

---

## Custom Domain (Optional)

To use a custom domain (e.g. `docs.yourcompany.com`):

1. Add a `CNAME` file in `docs/` with your domain.
2. Configure DNS `CNAME` record pointing to `sre-capsule.github.io`.
3. Enable **Enforce HTTPS** in GitHub Pages settings.

---

## Troubleshooting Failed Deployments

### Site shows the README instead of MkDocs

**Symptom:** https://sre-capsule.github.io/sre-capsule-community/ looks like the GitHub README (no navigation tabs, no search bar).

**Cause:** Pages is publishing from **`main`**, which has no built `index.html` — GitHub falls back to rendering `README.md`.

**Fix:**

1. Ensure **Deploy Documentation** workflow completed successfully (creates/updates `gh-pages`).
2. **Settings → Pages → Branch** → **`gh-pages`** / **`/ (root)`** (not `main`).
3. Hard-refresh the browser (Ctrl+F5).

### Deploy job fails with HTTP 404 (legacy Actions artifact mode)

If you still use the old two-job workflow with `actions/deploy-pages`, a 404 means Pages source was not **GitHub Actions**. The current workflow publishes to **`gh-pages`** instead — use branch **`gh-pages`** in Pages settings.

### Other errors

| Error | Cause | Fix |
|---|---|---|
| `Resource not accessible by integration` | Pages permissions missing | Workflow must have `pages: write` and `id-token: write` |
| `No GitHub Pages deployment found` | Pages source not set to Actions | Settings → Pages → Source → GitHub Actions |
| `mkdocs: command not found` | Dependencies not installed | Ensure `pip install -r docs/requirements.txt` runs before build |
| Strict mode warning | Broken internal link | Run `mkdocs build --strict` locally and fix links |
| Site not updating | Path filter skipped workflow | Push a change under `docs/` or run `workflow_dispatch` |
| 404 on site URL | Wrong `site_url` in mkdocs.yml | Match `site_url` to your actual Pages URL |

---

## Manual Deploy Trigger

If you need to redeploy without code changes:

1. Go to **Actions** tab on GitHub.
2. Select **Deploy Documentation**.
3. Click **Run workflow** → **Run workflow**.
