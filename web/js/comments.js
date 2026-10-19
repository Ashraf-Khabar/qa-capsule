/**
 * Incident comments: rich text, @mentions, and the resolve note.
 */
import { fetchWithAuth } from './api.js';
import { notify } from './ui.js';
import { avatarMarkup } from './avatars.js';

let directory = [];
let directoryLoaded = false;
let pendingAction = null;

const SIZE_OK = /^(12|14|16|18|20|24|28)px$/;

export function initCommentUI() {
    const modal = document.getElementById('comment-modal');
    if (!modal || modal.dataset.bound === '1') return;
    modal.dataset.bound = '1';
    modal.querySelectorAll('[data-comment-close]').forEach((el) => {
        el.addEventListener('click', closeCommentModal);
    });
    document.getElementById('comment-modal-submit')?.addEventListener('click', submitCommentModal);
    const toolbar = modal.querySelector('.rte-toolbar');
    toolbar?.addEventListener('mousedown', (event) => {
        if (event.target.closest('button, select, input')) event.preventDefault();
    });
    toolbar?.addEventListener('click', (event) => {
        const btn = event.target.closest('[data-cmd]');
        if (!btn) return;
        applyCommand(btn.dataset.cmd);
    });
    document.getElementById('rte-color')?.addEventListener('input', (event) => {
        focusEditor();
        document.execCommand('foreColor', false, event.target.value);
    });
    document.getElementById('rte-size')?.addEventListener('change', (event) => {
        wrapSelection((text) => `<span style="font-size:${event.target.value}">${text}</span>`);
    });
    const editor = document.getElementById('rte-editor');
    editor?.addEventListener('input', onEditorInput);
    editor?.addEventListener('keydown', onEditorKeydown);
    document.getElementById('rte-mentions')?.addEventListener('mousedown', (event) => {
        const item = event.target.closest('[data-mention]');
        if (!item) return;
        event.preventDefault();
        insertMention(item.dataset.mention);
    });
    document.addEventListener('click', (event) => {
        const btn = event.target.closest('[data-add-comment]');
        if (!btn) return;
        event.preventDefault();
        const id = parseInt(btn.dataset.addComment, 10);
        if (!Number.isNaN(id)) openCommentComposer([id]);
    });
}

export function openResolveComposer(ids, successMessage) {
    openComposer({
        title: 'Mark as solved',
        hint: 'Add a note for the team. Type @ to mention someone. The note is optional.',
        submit: 'Mark solved',
        required: false,
        ids,
        successMessage
    });
}

export function openCommentComposer(ids) {
    openComposer({
        title: 'Comment',
        hint: 'Type @ to mention a teammate. Bold, italic, underline, color, size, and code work like a ticket.',
        submit: 'Post comment',
        required: true,
        ids,
        successMessage: null
    });
}

function openComposer(action) {
    initCommentUI();
    pendingAction = action;
    const title = document.getElementById('comment-modal-title');
    const hint = document.getElementById('comment-modal-hint');
    const submit = document.getElementById('comment-modal-submit');
    const editor = document.getElementById('rte-editor');
    if (title) title.textContent = action.title;
    if (hint) hint.textContent = action.hint;
    if (submit) submit.textContent = action.submit;
    if (editor) editor.innerHTML = '';
    hideMentions();
    const modal = document.getElementById('comment-modal');
    if (modal) modal.hidden = false;
    ensureDirectory();
    editor?.focus();
}

function closeCommentModal() {
    const modal = document.getElementById('comment-modal');
    if (modal) modal.hidden = true;
    pendingAction = null;
    hideMentions();
}

async function submitCommentModal() {
    const action = pendingAction;
    if (!action) return;
    const blocks = serializeEditor(document.getElementById('rte-editor'));
    const hasText = blocks.some((block) => (block.t && block.t.trim()) || block.mention || (block.pre && block.pre.trim()));
    if (action.required && !hasText) {
        notify('Write a comment first', 'error');
        return;
    }
    const submit = document.getElementById('comment-modal-submit');
    if (submit) submit.disabled = true;
    try {
        if (hasText) {
            for (const id of action.ids) {
                const res = await fetchWithAuth(`/api/incidents/${id}/comments`, {
                    method: 'POST',
                    body: JSON.stringify({ blocks })
                });
                if (!res.ok) {
                    const message = await res.text();
                    throw new Error(message || 'Could not save the comment');
                }
            }
        }
        closeCommentModal();
        if (action.successMessage && typeof window.resolveIncidentsByIds === 'function') {
            action.ids.forEach((id) => window.selectedIncidents?.delete(String(id)));
            if (typeof window.updateBulkActionUI === 'function') window.updateBulkActionUI();
            await window.resolveIncidentsByIds(action.ids, action.successMessage);
        } else {
            notify('Comment posted', 'success');
            await hydrateOpenComments(true);
        }
    } catch (err) {
        notify(err.message || 'Could not save the comment', 'error');
    } finally {
        if (submit) submit.disabled = false;
    }
}

export async function hydrateOpenComments(force) {
    const nodes = [...document.querySelectorAll('.incident-comments[data-incident-id]')].filter((el) => {
        const body = el.closest('.pipeline-exec-body');
        return body && body.style.display !== 'none';
    });
    if (nodes.length === 0) return;
    const ids = [...new Set(nodes.map((el) => el.dataset.incidentId))];
    try {
        const res = await fetchWithAuth(`/api/incident-comments?ids=${ids.join(',')}`);
        if (!res.ok) return;
        const data = await res.json();
        const grouped = data.comments || {};
        nodes.forEach((el) => renderThread(el, grouped[el.dataset.incidentId] || [], force));
    } catch (_) { /* dashboard still usable without the thread */ }
}

function renderThread(el, comments) {
    const id = el.dataset.incidentId;
    const items = (comments || []).map((comment) => `
        <article class="comment-item">
            ${avatarMarkup({
                avatar_kind: comment.avatar_kind,
                avatar_value: comment.avatar_value,
                fullname: comment.fullname || comment.author,
                username: comment.author
            }, 32)}
            <div class="comment-item__body">
                <header class="comment-item__meta">
                    <strong>${escapeHtml(comment.fullname || comment.author)}</strong>
                    <span>@${escapeHtml(comment.author)}</span>
                    <time>${escapeHtml(comment.created_at || '')}</time>
                </header>
                <div class="comment-item__text">${comment.body_html || ''}</div>
            </div>
        </article>`).join('');
    el.innerHTML = `
        <div class="comment-thread">${items || '<p class="comment-empty">No notes yet.</p>'}</div>
        <button type="button" class="btn btn-secondary btn-sm" data-add-comment="${id}">Comment</button>`;
}

function applyCommand(cmd) {
    focusEditor();
    if (cmd === 'bold') document.execCommand('bold');
    if (cmd === 'italic') document.execCommand('italic');
    if (cmd === 'underline') document.execCommand('underline');
    if (cmd === 'code') wrapSelection((text) => `<code>${text}</code>`);
    if (cmd === 'pre') {
        const editor = document.getElementById('rte-editor');
        const selected = window.getSelection()?.toString() || 'code';
        document.execCommand('insertHTML', false, `<pre><code>${escapeHtml(selected)}</code></pre>`);
        editor?.focus();
    }
}

function wrapSelection(htmlForText) {
    focusEditor();
    const text = window.getSelection()?.toString() || '';
    if (!text) return;
    document.execCommand('insertHTML', false, htmlForText(escapeHtml(text)));
}

function focusEditor() {
    document.getElementById('rte-editor')?.focus();
}

function onEditorInput() {
    const query = mentionQuery();
    if (query == null) {
        hideMentions();
        return;
    }
    showMentions(query);
}

function onEditorKeydown(event) {
    const menu = document.getElementById('rte-mentions');
    if (!menu || menu.hidden) return;
    const items = [...menu.querySelectorAll('[data-mention]')];
    if (items.length === 0) return;
    const current = menu.querySelector('.is-active') || items[0];
    let index = items.indexOf(current);
    if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
        event.preventDefault();
        items.forEach((item) => item.classList.remove('is-active'));
        index = event.key === 'ArrowDown' ? (index + 1) % items.length : (index - 1 + items.length) % items.length;
        items[index].classList.add('is-active');
    }
    if (event.key === 'Enter' || event.key === 'Tab') {
        event.preventDefault();
        insertMention(current.dataset.mention);
    }
    if (event.key === 'Escape') hideMentions();
}

function mentionQuery() {
    const sel = window.getSelection();
    if (!sel || sel.rangeCount === 0 || !sel.isCollapsed) return null;
    const node = sel.anchorNode;
    if (!node || node.nodeType !== Node.TEXT_NODE) return null;
    const upto = node.textContent.slice(0, sel.anchorOffset);
    const match = upto.match(/@([a-zA-Z0-9._@+\-]{0,40})$/);
    return match ? match[1] : null;
}

async function ensureDirectory() {
    if (directoryLoaded) return;
    const res = await fetchWithAuth('/api/users/directory');
    if (!res.ok) return;
    const data = await res.json();
    directory = data.users || [];
    directoryLoaded = true;
}

async function showMentions(query) {
    await ensureDirectory();
    const q = query.toLowerCase();
    const matches = directory.filter((user) => {
        const name = `${user.username} ${user.fullname || ''}`.toLowerCase();
        return name.includes(q);
    }).slice(0, 8);
    const menu = document.getElementById('rte-mentions');
    if (!menu) return;
    if (matches.length === 0) {
        hideMentions();
        return;
    }
    menu.hidden = false;
    menu.innerHTML = matches.map((user, index) => `
        <button type="button" class="rte-mention${index === 0 ? ' is-active' : ''}" data-mention="${escapeHtml(user.username)}">
            ${avatarMarkup(user, 24)}
            <span><strong>${escapeHtml(user.fullname || user.username)}</strong><small>@${escapeHtml(user.username)}</small></span>
        </button>`).join('');
}

function hideMentions() {
    const menu = document.getElementById('rte-mentions');
    if (menu) {
        menu.hidden = true;
        menu.innerHTML = '';
    }
}

function insertMention(username) {
    const sel = window.getSelection();
    if (!sel || sel.rangeCount === 0) return;
    const node = sel.anchorNode;
    if (!node || node.nodeType !== Node.TEXT_NODE) return;
    const upto = node.textContent.slice(0, sel.anchorOffset);
    const match = upto.match(/@([a-zA-Z0-9._@+\-]{0,40})$/);
    if (!match) return;
    const start = sel.anchorOffset - match[0].length;
    const range = document.createRange();
    range.setStart(node, start);
    range.setEnd(node, sel.anchorOffset);
    range.deleteContents();
    const span = document.createElement('span');
    span.className = 'mention';
    span.dataset.user = username;
    span.setAttribute('contenteditable', 'false');
    span.textContent = `@${username}`;
    range.insertNode(span);
    const space = document.createTextNode(' ');
    span.after(space);
    const after = document.createRange();
    after.setStart(space, 1);
    after.collapse(true);
    sel.removeAllRanges();
    sel.addRange(after);
    hideMentions();
}

function serializeEditor(root) {
    if (!root) return [];
    const blocks = [];
    const walk = (node, marks) => {
        if (node.nodeType === Node.TEXT_NODE) {
            if (node.textContent) blocks.push({ t: node.textContent, ...marks });
            return;
        }
        if (node.nodeType !== Node.ELEMENT_NODE) return;
        if (node.nodeName === 'BR') {
            blocks.push({ br: true });
            return;
        }
        if (node.classList?.contains('mention')) {
            blocks.push({ mention: node.dataset.user || '' });
            return;
        }
        if (node.nodeName === 'PRE') {
            blocks.push({ pre: node.innerText || '' });
            return;
        }
        if (node.nodeName === 'DIV' || node.nodeName === 'P') {
            if (blocks.length > 0) blocks.push({ br: true });
            const blockMarks = { ...marks };
            node.childNodes.forEach((child) => walk(child, blockMarks));
            return;
        }
        const next = { ...marks };
        const tag = node.nodeName;
        if (tag === 'B' || tag === 'STRONG') next.b = true;
        if (tag === 'I' || tag === 'EM') next.i = true;
        if (tag === 'U') next.u = true;
        if (tag === 'CODE') next.code = true;
        if (tag === 'FONT' && node.color) {
            const color = normalizeColor(node.color);
            if (color) next.color = color;
        }
        if (tag === 'SPAN' && node.style) {
            const color = normalizeColor(node.style.color);
            if (color) next.color = color;
            if (SIZE_OK.test(node.style.fontSize)) next.size = node.style.fontSize;
        }
        node.childNodes.forEach((child) => walk(child, next));
    };
    root.childNodes.forEach((child) => walk(child, {}));
    return blocks.map((block) => {
        const clean = {};
        if (block.t != null && block.mention == null && block.pre == null && !block.br) clean.t = block.t;
        if (block.b) clean.b = true;
        if (block.i) clean.i = true;
        if (block.u) clean.u = true;
        if (block.code) clean.code = true;
        if (block.color) clean.color = block.color;
        if (block.size) clean.size = block.size;
        if (block.pre) clean.pre = block.pre;
        if (block.mention) clean.mention = block.mention;
        if (block.br) clean.br = true;
        return clean;
    }).filter((block) => block.t != null || block.pre || block.mention || block.br);
}

function normalizeColor(value) {
    if (!value) return '';
    const raw = String(value).trim();
    if (/^#[0-9a-fA-F]{6}$/.test(raw)) return raw.toLowerCase();
    const rgb = raw.match(/^rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)$/i);
    if (!rgb) return '';
    return '#' + rgb.slice(1, 4).map((n) => Number(n).toString(16).padStart(2, '0')).join('');
}

function escapeHtml(value) {
    return String(value ?? '')
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');
}
