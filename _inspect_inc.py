import os, sqlite3
db = os.path.join(os.environ["TEMP"], "sre-incident.db")
print("DB", db, os.path.getsize(db))
c = sqlite3.connect(db)
print("TABLES", c.execute("SELECT name FROM sqlite_master WHERE type='table'").fetchall())
print("IDS", c.execute("SELECT id, substr(name,1,80) FROM incidents ORDER BY id DESC LIMIT 10").fetchall())
row = c.execute(
    """SELECT id,
              length(COALESCE(error_message,'')),
              length(COALESCE(console_logs,'')),
              length(COALESCE(error_logs,'')),
              substr(COALESCE(error_message,''),1,200)
       FROM incidents WHERE id=5"""
).fetchone()
print("ROW", row)
blob = c.execute(
    "SELECT COALESCE(error_message,'') || '\n' || COALESCE(console_logs,'') || '\n' || COALESCE(error_logs,'') FROM incidents WHERE id=5"
).fetchone()[0]
for needle in ("SRE_CAPSULE_DOM", "Secure Area", "Login Page", "/logout", "logout-button", "sign-in-button"):
    print(needle, blob.find(needle))
print("--- error_logs head ---")
logs = c.execute("SELECT COALESCE(error_logs,'') FROM incidents WHERE id=5").fetchone()[0]
print(logs[:1500])
print("--- console head ---")
print(c.execute("SELECT COALESCE(console_logs,'') FROM incidents WHERE id=5").fetchone()[0][:800])
meta = c.execute("PRAGMA table_info(incidents)").fetchall()
print("COLS", [m[1] for m in meta])
em = c.execute("SELECT error_message FROM incidents WHERE id=5").fetchone()[0]
i = em.find("SRE_CAPSULE_DOM_SNAPSHOT_START")
j = em.find("SRE_CAPSULE_DOM_SNAPSHOT_END")
print("marker", i, j, "len", len(em))
chunk = em[i:j+40] if i>=0 else em[-500:]
# keep only a short structural view
import re
hrefs = re.findall(r'href="[^"]+"', chunk)
print("HREFS", hrefs[:20])
print("HAS_LOGOUT_ANCHOR", 'href="/logout"' in chunk or "href='/logout'" in chunk)
print("CHUNK_LEN", len(chunk))
print(chunk[:1200])
print("--- TAIL ---")
print(chunk[-800:])
