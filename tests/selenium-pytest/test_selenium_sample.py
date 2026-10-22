"""
Selenium + Pytest suite.

TC-01..TC-03 check login on the-internet.herokuapp.com.
The self-healing failure is not login or logout. It opens a local alert
board with five lookalike controls and clicks a selector that is not on
the page. The only correct control is the active sev1 row
"disk latency on shard 7" (#ack-lat-shard-7).
"""

import sys
from pathlib import Path

import pytest
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait

BASE_URL = "https://the-internet.herokuapp.com/login"
USERNAME = "tomsmith"
PASSWORD = "SuperSecretPassword!"
BOARD = Path(__file__).resolve().parent / "fixtures" / "shard_alerts.html"
# Deliberately unrelated to every id on the page.
OUTDATED_ACK = '[data-testid="confirm-maintenance-window"]'


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


def test_acknowledge_only_active_sev1_disk_latency_on_shard_7(driver):
    """Acknowledge only the active sev1 alert titled disk latency on shard 7.

    Leave shard 17, shard 70, the preview draft, and the history link alone.
    """
    driver.get(BOARD.as_uri())
    assert "Shard alert board" in driver.title
    assert driver.find_element(By.ID, "ack-lat-shard-7").is_displayed()
    sys.stderr.write(
        "[STDERR] Acknowledging the active sev1 disk latency on shard 7 "
        f"with outdated selector {OUTDATED_ACK}\n"
    )
    driver.find_element(By.CSS_SELECTOR, OUTDATED_ACK).click()
