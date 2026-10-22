"""Pytest hooks — capture Selenium page HTML on failure for SRE Capsule MCP healing."""

import pytest

MAX_HTML_CHARS = 24_000


def _build_snapshot(html: str) -> str:
    trimmed = (html or "")[:MAX_HTML_CHARS]
    return (
        f"[SRE_CAPSULE_DOM_SNAPSHOT_START]\n"
        f"{trimmed}\n"
        f"[SRE_CAPSULE_DOM_SNAPSHOT_END]"
    )


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_makereport(item, call):
    outcome = yield
    report = outcome.get_result()
    if report.when != "call" or not report.failed:
        return

    driver = item.funcargs.get("driver")
    if driver is None:
        return

    try:
        snapshot = _build_snapshot(driver.page_source)
        report.sections.append(("Captured output", snapshot))
        if report.longrepr:
            report.longrepr = f"{report.longrepr}\n{snapshot}"
    except Exception:
        pass
