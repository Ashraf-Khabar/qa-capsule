"""Edge case for locator healing. Not a login and not a logout.

The board shows five lookalike controls. Four mention shard 7 or say
Acknowledge. The test must acknowledge only the active sev1 row
"disk latency on shard 7".

The outdated selector names a maintenance window that is not on the page,
so a keyword scan cannot guess. The page HTML is captured on failure.
The right control is #ack-lat-shard-7.
"""

import sys
from pathlib import Path

import pytest
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By

BOARD = Path(__file__).resolve().parent / "fixtures" / "shard_alerts.html"
# Deliberately unrelated to every id on the page.
OUTDATED_ACK = '[data-testid="confirm-maintenance-window"]'


@pytest.fixture
def driver():
    opts = Options()
    opts.add_argument("--headless")
    opts.add_argument("--no-sandbox")
    opts.add_argument("--disable-dev-shm-usage")
    opts.add_argument("--window-size=1280,800")
    drv = webdriver.Chrome(options=opts)
    yield drv
    drv.quit()


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
