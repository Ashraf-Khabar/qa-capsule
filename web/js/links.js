/**
 * Encrypted navigation. The address is /XXXX/YYYY/<hash>. Page names and record
 * ids stay inside the hash, not in the two segments.
 */
import { fetchWithAuth, parseApiJson } from './api.js';

const minted = new Map();
let navSeq = 0;

/** Token from /XXXX/YYYY/<hash> or the older /l/<hash> form. */
export function secureLinkTokenFromPath(pathname) {
    const path = pathname || '/';
    const legacy = path.match(/^\/l\/([^/]+)$/);
    if (legacy) return decodeURIComponent(legacy[1]);
    const dynamic = path.match(/^\/([a-z0-9]{4})\/([a-z0-9]{4})\/([^/]+)$/);
    if (dynamic) return decodeURIComponent(dynamic[3]);
    return '';
}

function destKey(dest) {
    return JSON.stringify({
        kind: dest.kind || '',
        view: dest.view || '',
        topic: dest.topic || '',
        project: dest.project || '',
        run_id: dest.run_id || '',
        incident_id: dest.incident_id || 0,
        comment_id: dest.comment_id || 0
    });
}

async function mintToken(dest) {
    const key = destKey(dest);
    if (minted.has(key)) return minted.get(key);
    const pending = fetchWithAuth('/api/links', {
        method: 'POST',
        body: JSON.stringify(dest)
    }).then(async (res) => {
        const { data } = await parseApiJson(res);
        if (!res.ok || !data?.token) {
            minted.delete(key);
            throw new Error((data && data.error) || 'Could not build the link');
        }
        return { token: data.token, path: data.path || `/l/${data.token}` };
    });
    minted.set(key, pending);
    return pending;
}

/** Put the current screen on /XXXX/YYYY/<hash>. */
export async function publishSecureLink(dest) {
    if (window.__openingSecureLink || !dest?.kind) return;
    const seq = ++navSeq;
    const here = window.location.pathname || '/';
    const replace = !secureLinkTokenFromPath(here);
    try {
        const mintedLink = await mintToken(dest);
        if (seq !== navSeq || window.__openingSecureLink) return;
        const path = mintedLink.path;
        if (!path || window.location.pathname === path) return;
        if (replace) history.replaceState({ sreLink: seq }, '', path);
        else history.pushState({ sreLink: seq }, '', path);
    } catch (_) {
        /* navigation still works if the address cannot be updated */
    }
}

window.addEventListener('popstate', () => {
    const token = secureLinkTokenFromPath(window.location.pathname || '');
    if (token && typeof window.openSecureLink === 'function') {
        window.openSecureLink(token);
    }
});
