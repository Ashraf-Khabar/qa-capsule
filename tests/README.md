# Automated tests (multi-framework)

**QA Capsule Community** — realistic demo suites for CI/CD ingestion, incident correlation, and **MCP self-healing**. Each pipeline runs against **public demo sites** and **fails on purpose** on the last scenario so JUnit upload + MCP `heal_incident` can be exercised from the dashboard or IDE.

**Docs:** [Test frameworks catalog](https://qa-capsule.github.io/qa-capsule-community/integration/test-frameworks/) · [MCP self-healing](https://qa-capsule.github.io/qa-capsule-community/guides/mcp-self-healing-testing/)

## Demo targets

| Framework | Site / API | Passing | Intentional failure |
|-----------|------------|---------|---------------------|
| Cypress | [saucedemo.com](https://www.saucedemo.com/) | Login + cart | Broken checkout locator |
| Playwright | saucedemo.com | Login + cart | Broken checkout locator |
| Robot | saucedemo.com + [reqres.in](https://reqres.in/) | E2E + API | Broken locator + wrong HTTP status |
| Selenium | [the-internet.herokuapp.com](https://the-internet.herokuapp.com/login) | Form auth | Broken submit selector |
| Pytest API | reqres.in | Live REST | Contract drift assertions |
| Newman | reqres.in | Live REST | Contract drift folder |
| JUnit Java | reqres.in | `HttpClient` | Contract drift + runtime error |

Credentials (public demos): `standard_user` / `secret_sauce` (Swag Labs), `tomsmith` / `SuperSecretPassword!` (The Internet).

## Prerequisites

- Python 3.10+
- Node.js 18+ (Cypress, Playwright, Newman)
- Java 21+ (JUnit)
- `bash` (Git Bash on Windows)

## Layout

```
tests/
├── run-framework.sh
├── upload-junit.sh
├── robotframework/     saucedemo_checkout.robot, api_health.robot
├── playwright/         saucedemo_checkout.spec.js
├── cypress/            saucedemo_checkout.cy.js
├── pytest/             reqres.in API tests
├── selenium-pytest/    the-internet.herokuapp.com login
├── newman/             reqres.in Postman collection
└── junit-java/         reqres.in HttpClient tests
```

## Run locally

```bash
chmod +x tests/run-framework.sh tests/*/run.sh tests/upload-junit.sh
./tests/run-framework.sh robot
./tests/run-framework.sh playwright
./tests/run-framework.sh cypress
./tests/run-framework.sh pytest
./tests/run-framework.sh selenium-py
./tests/run-framework.sh newman
./tests/run-framework.sh junit-java
```

### Upload to QA Capsule

```bash
export QA_CAPSULE_URL="http://localhost:9000"
export QA_CAPSULE_API_KEY="your-project-api-key"
./tests/run-framework.sh robot
```

### API host override (Robot)

```bash
export API_HEALTH_HOST=reqres.in
./scripts/run-tests.sh
```

### DOM capture (MCP self-healing)

Every **browser** framework embeds page HTML between `[QA_CAPSULE_DOM_SNAPSHOT_START/END]` markers on failure:

| Framework | Hook |
|-----------|------|
| Robot | `tests/robotframework/resources/dom_capture_listener.py` (`--listener` in CI) |
| Cypress | `tests/cypress/cypress/support/qa-capsule.js` |
| Playwright | `tests/playwright/fixtures/qa-capsule-fixture.js` |
| Selenium | `tests/selenium-pytest/conftest.py` |

API suites (Pytest/Newman/JUnit) have no browser DOM — healing uses live URL fetch or error text.

## CI/CD (GitHub Actions)

All framework workflows are **manual** (`workflow_dispatch`). Each run:

1. Executes real tests against the demo site/API
2. Uploads JUnit XML to QA Capsule (`continue-on-error` on tests)
3. **Fails the workflow** (red pipeline) so failures are visible in GitHub

Use **MCP `heal_incident`** or the dashboard **Propose fix** action to analyze incidents after upload.

| Workflow | Framework |
|----------|-----------|
| `e2e-tests-robot.yml` | Robot Framework |
| `e2e-tests-playwright.yml` | Playwright |
| `e2e-tests-cypress.yml` | Cypress |
| `e2e-tests-selenium-pytest.yml` | Selenium |
| `api-tests-pytest.yml` | Pytest |
| `api-tests-newman.yml` | Newman |
| `api-tests-junit-java.yml` | JUnit Java |

### Secrets

| Secret | Purpose |
|--------|---------|
| `QA_CAPSULE_URL` | Control plane base URL |
| `QA_CAPSULE_API_*_KEY` | Per-framework project API key |

Shared upload script: `.github/scripts/upload-junit-to-qacapsule.sh`.
