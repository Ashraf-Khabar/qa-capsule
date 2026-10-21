---
icon: material/book-open-variant
---

# SRE theory and SRE Capsule

Each section states the practice, then the screen, control, and role in SRE Capsule. The same text is the first topic of the in-app **Documentation** page.

SRE Capsule records what happened. It does not invent a cause. The incident **Runbook** button does not call Slack, Jira, or a cluster. Those calls happen only when a workflow you already enabled runs on ingest.

---

## Detect

An incident starts when a signal says the service or the test is no longer doing what it should. The response clock starts at detection, not when a person later notices.

**In SRE Capsule.** CI posts the failure to the gateway webhook (`X-API-Key`, `X-Run-Id`). The row appears on **Telemetry Stream**. *Detected* is `created_at`.

---

## Acknowledge

Acknowledgement means a named person has taken the incident. Severity says how badly it hurts users. A second acknowledgement does not replace the first owner.

**In SRE Capsule.** A lead or manager opens the incident, picks **sev1** to **sev4**, and clicks **Acknowledge**. If nobody does that for 15 minutes, the row says *Waiting N min for an owner*.

---

## Mitigate, then resolve

Mitigation stops the user impact. Resolution closes the record. They are different moments.

**In SRE Capsule.** **Mark mitigated** writes the mitigate time. **Resolve** on the pipeline card sets the incident resolved. The timeline reads Detected, Acknowledged, Mitigated, Resolved.

---

## One incident for the same cause

A spike should be one incident with a count, not a new page for every identical failure. The original rows stay for the audit.

**In SRE Capsule.** On Telemetry Stream, failures of the same gateway and the same test (or the same fingerprint) inside **6 hours** fold into one row. The badge lists the signal numbers and times. Nothing is deleted from the database.

---

## Change failure and the error budget

DORA change-failure rate is failed changes divided by all changes in the window. An error budget is how many of those failures you agreed to spend.

```
allowed  = deployments × SLO
burn     = failed deployments / allowed
```

| Burn | Status on the DORA card |
|---|---|
| No runs | Nothing was deployed in the window |
| Under 0.8 | Within budget |
| 0.8 up to 1 | Close to budget |
| 1 or more | Freeze deploys |

*Freeze deploys* is a warning on the screen. SRE Capsule does not block the pipeline.

**In SRE Capsule.** A manager opens **DORA**. The error-budget table shows runs, failed runs, budget left, and burn. The SLO is the allowed failure rate for that gateway, from **1% to 50%**. **Save** stores it. The default is **15%**.

```http
POST /api/dora/slo
Authorization: Bearer <jwt>
Content-Type: application/json

{ "project": "checkout-api", "slo_percent": 15 }
```

Manager role required. An unknown gateway returns 404. A percent outside 1–50 is rejected.

---

## The change next to the incident

Most incidents are easier to read next to the deploy that just landed. A change far from the incident is context, not a cause.

**In SRE Capsule.** Each pipeline card shows the nearest deploy: short commit, branch, outcome, and minutes before the incident. The link prefers the incident’s own pipeline run. Otherwise it uses the latest run that started at or before the incident. If that deploy is within **30 minutes** and the role can open workflows, **Open workflow** jumps to that gateway’s DAG.

---

## Runbook

A runbook is a known sequence for a known symptom. You confirm it. A word in a log does not fire a production action by itself.

**In SRE Capsule.** **Runbook** on the incident asks the plugin engine which catalog steps match (502, OOM, flaky, timeout, regression). The confirmation names those steps. Saving records the dry-run on the incident, with the note that no external system was called. A symptom with no matching step is left alone.

Applying a template onto a gateway’s workflow (the older **Apply** path) is a separate action. That writes `sre_workflow_json`. The incident button does not.

---

## Postmortem

A postmortem is the chronology of what was recorded, written after the incident is closed. A missing fact stays missing.

**In SRE Capsule.** After resolve, **Write postmortem** drafts from the timeline, comments, and the nearest deploy. If AI is enabled, the model may rewrite that draft and is instructed not to add a cause, a person, or a time that is not already there. If AI is off or does not answer, the factual draft is saved. **Save draft** keeps your edit.

---

## Maintenance

A maintenance window says the team already knows why the gateway may fail. Hiding those incidents would erase the record. Marking them keeps the queue honest.

**In SRE Capsule.** On **DORA**, under the error budget, a manager sets gateway, reason, and **1 to 24 hours**, then **Start**. Telemetry Stream shows a banner. Incidents that start during the window stay in the list and carry the reason. **End** closes the window. One open window per gateway.

```http
POST /api/maintenance
Authorization: Bearer <jwt>
Content-Type: application/json

{ "project": "checkout-api", "reason": "database upgrade", "hours": 4 }
```

```http
POST /api/maintenance
{ "action": "end", "project": "checkout-api" }
```

Manager or platform admin. `GET /api/maintenance` is any signed-in user, limited to gateways they can see.

---

## Where to click

| Practice | Screen | Who |
|---|---|---|
| Detect, acknowledge, mitigate, same cause, runbook, postmortem | Telemetry Stream | Lead or manager for the buttons. Readers see the record. |
| Error budget, SLO, maintenance | DORA | Manager |
| AI rewrite of the postmortem | AI Configuration, then Telemetry Stream | Manager enables the provider. Lead writes the draft. |

---

## Related

- [Incident Lifecycle](incident-lifecycle.md)
- [Runbooks & DORA](runbooks-dora.md)
- [Dashboard Operations](dashboard-operations.md)
