FROM python:3.12-slim-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends chromium chromium-driver ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && command -v chromium \
    && command -v chromedriver

WORKDIR /work
COPY tests/selenium-pytest/requirements.txt /tmp/selenium-req.txt
COPY tests/pytest/requirements.txt /tmp/pytest-req.txt
RUN pip install --no-cache-dir -r /tmp/selenium-req.txt -r /tmp/pytest-req.txt

COPY tests/runner/server.py /opt/runner/server.py
COPY tests/runner/runner_chrome.py /opt/runner/runner_chrome.py
ENV PYTHONPATH=/opt/runner
ENV TEST_ROOT=/work
ENV CHROME_BIN=/usr/bin/chromium
ENV CHROMEDRIVER=/usr/bin/chromedriver
ENV PORT=8091

EXPOSE 8091
CMD ["python", "/opt/runner/server.py"]
