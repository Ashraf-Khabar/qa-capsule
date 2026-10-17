/**
 * SRE Capsule support file for Cypress.
 *
 * On failure, captures page HTML synchronously (cy.state('window')) so the
 * snapshot is embedded in the Mocha/JUnit failure message before upload.
 * Using cy.window() after a failed cy.get() often never runs — Cypress aborts
 * the command queue.
 */

const MAX_HTML_CHARS = 24_000;

function buildSnapshot(win) {
  const html = win.document.documentElement.outerHTML.slice(0, MAX_HTML_CHARS);
  return (
    `[SRE_CAPSULE_DOM_SNAPSHOT_START]\n` +
    `${html}\n` +
    `[SRE_CAPSULE_DOM_SNAPSHOT_END]`
  );
}

function appendSnapshotToFailedTest(test) {
  if (!test || test.state !== 'failed') {
    return;
  }
  try {
    const win = cy.state('window');
    if (!win || !win.document) {
      return;
    }
    const snapshot = buildSnapshot(win);
    if (test.err) {
      test.err.message = `${test.err.message || ''}\n${snapshot}`;
    }
    // Best-effort Node log (may not run if the command queue is dead).
    try {
      cy.task('sreLog', snapshot, { log: false });
    } catch (_) {
      // ignore
    }
  } catch (_) {
    // ignore
  }
}

afterEach(function () {
  appendSnapshotToFailedTest(this.currentTest);
});
