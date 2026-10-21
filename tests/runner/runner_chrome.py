"""Points Selenium at the container's Chromium when CHROME_BIN is set."""

import os

import selenium.webdriver
from selenium.webdriver.chrome.service import Service

_chrome = selenium.webdriver.Chrome


def Chrome(options=None, service=None, **kwargs):
    binary = os.environ.get("CHROME_BIN", "")
    driver = os.environ.get("CHROMEDRIVER", "")
    if options is not None and binary:
        options.binary_location = binary
    if service is None and driver:
        service = Service(driver)
    return _chrome(options=options, service=service, **kwargs)


selenium.webdriver.Chrome = Chrome
