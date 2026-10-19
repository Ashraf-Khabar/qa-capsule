#!/usr/bin/env bash
# Upload a JUnit XML report to SRE Capsule (skips gracefully when secrets or file are missing).
set -euo pipefail

FILE="${1:-}"
FRAMEWORK="${2:-JUnit}"
RUN_ID="${3:-local}"

WEBHOOK_URL="${WEBHOOK_URL:-${SRE_CAPSULE_URL:-}}"
API_KEY="${API_KEY:-${API_KEY_LEGACY:-}}"

if [ -z "$WEBHOOK_URL" ] || [ -z "$API_KEY" ]; then
  echo "ERROR: SRE Capsule upload skipped — SRE_CAPSULE_URL or API_KEY is empty."
  echo "Add repository secrets SRE_CAPSULE_URL (public https base URL, no path) and the project API key."
  echo "Nothing was sent, so the dashboard will stay empty."
  exit 1
fi

# The UI shows https://host/api/webhooks/<ci>. Ingest expects the origin only.
WEBHOOK_URL="${WEBHOOK_URL%/}"
WEBHOOK_URL="${WEBHOOK_URL%%/api/webhooks*}"

if [ -z "$FILE" ] || [ ! -f "$FILE" ]; then
  echo "JUnit report not found: ${FILE:-<unset>}"
  exit 0
fi

EXEC_ENV="${SRE_CAPSULE_EXEC_ENV:-STAGING}"
EXEC_TYPE="${SRE_CAPSULE_EXEC_TYPE:-TEST-RUN}"

set +e
HTTP_STATUS=$(curl -sS -o /tmp/sre-ingest.json -w "%{http_code}" \
  -X POST "${WEBHOOK_URL}/api/webhooks/upload?framework=${FRAMEWORK}" \
  -H "X-API-Key: ${API_KEY}" \
  -H "X-Run-Id: ${RUN_ID}" \
  -H "X-Commit-Sha: ${GITHUB_SHA:-}" \
  -H "X-Branch: ${GITHUB_REF_NAME:-}" \
  -H "X-Execution-Env: ${EXEC_ENV}" \
  -H "X-Execution-Type: ${EXEC_TYPE}" \
  -F "file=@${FILE}" \
  2>/tmp/sre-ingest.err)
CURL_EXIT=$?
set -e

if [ "$CURL_EXIT" -ne 0 ]; then
  echo "ERROR: curl exit ${CURL_EXIT} — GitHub could not reach SRE Capsule."
  echo "A runner cannot open localhost or a private address. SRE_CAPSULE_URL must be a public https origin, with no /api/webhooks path."
  cat /tmp/sre-ingest.err || true
  exit 1
fi

echo "SRE Capsule HTTP status: ${HTTP_STATUS}"
cat /tmp/sre-ingest.json || true

if [ "$HTTP_STATUS" -ge 400 ] 2>/dev/null; then
  echo "ERROR: SRE Capsule ingestion returned HTTP ${HTTP_STATUS}"
  exit 1
fi

echo "Results queued for processing."
