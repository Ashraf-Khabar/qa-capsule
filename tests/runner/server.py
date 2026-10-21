#!/usr/bin/env python3
"""Runs one pytest node from the mounted test tree.

This process is meant to stay on the Docker network. It refuses any id that
is not a relative pytest path.
"""

import json
import os
import re
import subprocess
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

ROOT = os.environ.get("TEST_ROOT", "/work")
TOKEN = os.environ.get("SRE_VERIFY_TOKEN", "")
PORT = int(os.environ.get("PORT", "8091"))
NODE_RE = re.compile(r"^[\w./-]+\.py(?:::[A-Za-z_][\w]*){1,3}$")
SNAPSHOT_RE = re.compile(
    r"SRE_CAPSULE_DOM_SNAPSHOT_START.*?SRE_CAPSULE_DOM_SNAPSHOT_END",
    re.S,
)


def public_output(text):
    text = SNAPSHOT_RE.sub("[page snapshot omitted]", text or "")
    text = text.strip()
    if len(text) > 1500:
        text = text[-1500:]
    return text


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/healthz":
            self._json(200, {"status": "ok"})
            return
        self.send_error(404)

    def do_POST(self):
        if self.path != "/run":
            self.send_error(404)
            return
        if TOKEN and self.headers.get("X-Verify-Token") != TOKEN:
            self._json(401, {"error": "the test runner rejected the request"})
            return
        length = int(self.headers.get("Content-Length", "0") or "0")
        if length <= 0 or length > 4000:
            self._json(400, {"error": "invalid body"})
            return
        try:
            body = json.loads(self.rfile.read(length).decode("utf-8"))
        except (UnicodeDecodeError, json.JSONDecodeError):
            self._json(400, {"error": "invalid body"})
            return
        node = str(body.get("node") or "")
        if not NODE_RE.match(node) or ".." in node:
            self._json(400, {"error": "refusing this test id"})
            return
        file_part = node.split("::", 1)[0]
        root_real = os.path.realpath(ROOT)
        abs_file = os.path.realpath(os.path.join(ROOT, file_part))
        if not abs_file.startswith(root_real + os.sep) or not os.path.isfile(abs_file):
            self._json(404, {"error": "the test file was not found"})
            return
        try:
            proc = subprocess.run(
                ["pytest", "-q", "--tb=line", "-p", "runner_chrome", node],
                cwd=ROOT,
                capture_output=True,
                text=True,
                timeout=140,
            )
        except subprocess.TimeoutExpired:
            self._json(200, {"passed": False, "exit_code": 124, "output": "The test timed out."})
            return
        output = public_output((proc.stdout or "") + "\n" + (proc.stderr or ""))
        self._json(200, {"passed": proc.returncode == 0, "exit_code": proc.returncode, "output": output})

    def _json(self, code, payload):
        raw = json.dumps(payload).encode("utf-8")
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(raw)))
        self.end_headers()
        self.wfile.write(raw)

    def log_message(self, fmt, *args):
        return


if __name__ == "__main__":
    ThreadingHTTPServer(("0.0.0.0", PORT), Handler).serve_forever()
