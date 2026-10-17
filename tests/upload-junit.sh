#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 2 ]]; then
  echo "Usage: tests/upload-junit.sh <framework> <junit_file>"
  exit 1
fi

FRAMEWORK="$1"
JUNIT_FILE="$2"

SRE_CAPSULE_URL="${SRE_CAPSULE_URL:-${WEBHOOK_URL:-}}"
SRE_CAPSULE_API_KEY="${SRE_CAPSULE_API_KEY:-${API_KEY:-}}"
RUN_ID="${CI_PIPELINE_ID:-${GITHUB_RUN_ID:-local-$(date +%s)}}"
EXEC_ENV="${SRE_CAPSULE_EXEC_ENV:-DEV}"
EXEC_TYPE="${SRE_CAPSULE_EXEC_TYPE:-TEST-RUN}"

if [[ ! -f "${JUNIT_FILE}" ]]; then
  echo "ERROR: JUnit file not found: ${JUNIT_FILE}"
  exit 1
fi

if [[ -z "${SRE_CAPSULE_URL}" || -z "${SRE_CAPSULE_API_KEY}" ]]; then
  echo "INFO: SRE_CAPSULE_URL / SRE_CAPSULE_API_KEY not set; upload skipped."
  exit 0
fi

echo "Uploading ${FRAMEWORK} JUnit file: ${JUNIT_FILE}"
MAX_RETRIES="${SRE_CAPSULE_UPLOAD_MAX_RETRIES:-4}"
BACKOFF_SECONDS="${SRE_CAPSULE_UPLOAD_BACKOFF_SECONDS:-3}"
ATTEMPT=1

while true; do
  set +e
  HTTP_CODE="$(curl -sS -o /tmp/sre_capsule_upload_response.txt -w "%{http_code}" \
    -X POST "${SRE_CAPSULE_URL}/api/webhooks/upload?framework=${FRAMEWORK}" \
    -H "X-API-Key: ${SRE_CAPSULE_API_KEY}" \
    -H "X-Run-Id: ${RUN_ID}" \
    -H "X-Execution-Env: ${EXEC_ENV}" \
    -H "X-Execution-Type: ${EXEC_TYPE}" \
    -F "file=@${JUNIT_FILE}")"
  CURL_EXIT=$?
  set -e

  if [[ "${CURL_EXIT}" -eq 0 && "${HTTP_CODE}" =~ ^2 ]]; then
    echo "Upload accepted by SRE Capsule (HTTP ${HTTP_CODE})."
    break
  fi

  echo "Upload attempt ${ATTEMPT}/${MAX_RETRIES} failed (curl=${CURL_EXIT}, http=${HTTP_CODE})."
  if [[ -f /tmp/sre_capsule_upload_response.txt ]]; then
    echo "Response body:"
    sed -n '1,20p' /tmp/sre_capsule_upload_response.txt || true
  fi

  # Retry only transient network/server responses to keep CI simple and reliable.
  if [[ "${ATTEMPT}" -lt "${MAX_RETRIES}" ]] && [[ "${CURL_EXIT}" -ne 0 || "${HTTP_CODE}" =~ ^(429|500|502|503|504)$ ]]; then
    sleep "${BACKOFF_SECONDS}"
    ATTEMPT=$((ATTEMPT + 1))
    continue
  fi

  echo "ERROR: SRE Capsule upload failed after ${ATTEMPT} attempt(s)."
  exit 1
done
