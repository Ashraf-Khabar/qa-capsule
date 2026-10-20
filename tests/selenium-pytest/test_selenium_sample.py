"""
Form authentication on the-internet.herokuapp.com (Elemental Selenium demo site).

Target: https://the-internet.herokuapp.com/login
Credentials: tomsmith / SuperSecretPassword!

TC-01..TC-03 pass. TC-04 logs in with the real form, then fails on the
secure-area logout control so MCP self-healing can be tried after login.
"""

import sys

import pytest
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait

BASE_URL = "https://the-internet.herokuapp.com/login"
USERNAME = "tomsmith"
PASSWORD = "SuperSecretPassword!"
BROKEN_LOGOUT = '[data-testid="logout-button"]'


@pytest.fixture(scope="function")
def driver():
    opts = Options()
    opts.add_argument("--headless")
    opts.add_argument("--no-sandbox")
    opts.add_argument("--disable-dev-shm-usage")
    opts.add_argument("--window-size=1280,800")
    drv = webdriver.Chrome(options=opts)
    yield drv
    drv.quit()


def _open_login(driver):
    print(f"[STDOUT] Navigating to {BASE_URL}")
    driver.get(BASE_URL)


def _login(driver):
    _open_login(driver)
    driver.find_element(By.ID, "username").send_keys(USERNAME)
    driver.find_element(By.ID, "password").send_keys(PASSWORD)
    driver.find_element(By.CSS_SELECTOR, "button[type='submit']").click()
    WebDriverWait(driver, 5).until(EC.url_contains("/secure"))


class TestFormAuthentication:
    def test_login_page_heading(self, driver):
        _open_login(driver)
        heading = driver.find_element(By.CSS_SELECTOR, "h2")
        assert "Login Page" in heading.text

    def test_valid_credentials_reach_secure_area(self, driver):
        _login(driver)
        assert "Secure Area" in driver.find_element(By.CSS_SELECTOR, "h4.subheader").text

    def test_logout_link_visible_in_secure_area(self, driver):
        _login(driver)
        WebDriverWait(driver, 5).until(
            EC.presence_of_element_located((By.CSS_SELECTOR, "a[href='/logout']"))
        )
        assert driver.find_element(By.CSS_SELECTOR, "a[href='/logout']").is_displayed()

    def test_broken_logout_after_login_self_healing_demo(self, driver):
        _login(driver)
        assert "Secure Area" in driver.find_element(By.CSS_SELECTOR, "h2").text
        sys.stderr.write(f"[STDERR] Logged in. Using outdated logout selector {BROKEN_LOGOUT}\n")
        driver.find_element(By.CSS_SELECTOR, BROKEN_LOGOUT).click()
        WebDriverWait(driver, 5).until(EC.url_contains("/login"))
