/**
 * SRE Capsule test fixture for Playwright.
 *
 * On failure, captures page HTML and embeds DOM markers in the test error message
 * (JUnit failure body) plus a named attachment for sre-capsule-reporter.js.
 */

const { test: baseTest, expect } = require('@playwright/test');

const MAX_HTML_CHARS = 24_000;

function buildSnapshot(html) {
  const trimmed = html.slice(0, MAX_HTML_CHARS);
  return (
    `[SRE_CAPSULE_DOM_SNAPSHOT_START]\n` +
    `${trimmed}\n` +
    `[SRE_CAPSULE_DOM_SNAPSHOT_END]`
  );
}

const test = baseTest.extend({
  page: async ({ page }, use, testInfo) => {
    await use(page);

    if (testInfo.status === testInfo.expectedStatus) {
      return;
    }

    try {
      const html = await page.content();
      if (!html) {
        return;
      }
      const snapshot = buildSnapshot(html);
      await testInfo.attach('sre-capsule-dom-snapshot', {
        body: Buffer.from(html, 'utf8'),
        contentType: 'text/html',
      });
      if (testInfo.errors.length > 0) {
        testInfo.errors[0].message = `${testInfo.errors[0].message || ''}\n${snapshot}`;
      } else {
        testInfo.errors.push({ message: snapshot });
      }
    } catch (_) {
      // Page may already be closed — best effort only.
    }
  },
});

module.exports = { test, expect };
