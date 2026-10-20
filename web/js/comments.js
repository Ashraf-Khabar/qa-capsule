/**
 * Incident comments: rich text, @mentions, and the resolve note.
 */
import { fetchWithAuth, parseJwt } from './api.js';
import { notify, showConfirmModal } from './ui.js';
import { avatarMarkup } from './avatars.js';

let directory = [];
let directoryLoaded = false;
let pendingAction = null;
const commentHTML = new Map();
const expandedThreads = new Set();

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
        const swatch = event.target.closest('[data-color]');
        if (swatch) {
            applyColor(swatch.dataset.color);
            closePops();
            return;
        }
        const emoji = event.target.closest('[data-emoji]');
        if (emoji) {
            insertEmoji(emoji.dataset.emoji);
            closePops();
            return;
        }
        const popBtn = event.target.closest('[data-pop]');
        if (popBtn) {
            togglePop(popBtn.dataset.pop);
            return;
        }
        const btn = event.target.closest('[data-cmd]');
        if (!btn) return;
        closePops();
        applyCommand(btn.dataset.cmd);
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
        if (!event.target.closest('.rte-popwrap')) closePops();
        const more = event.target.closest('[data-thread-more]');
        if (more) {
            event.preventDefault();
            event.stopPropagation();
            toggleThread(more.dataset.threadMore);
            return;
        }
        const edit = event.target.closest('[data-edit-comment]');
        if (edit) {
            event.preventDefault();
            event.stopPropagation();
            openEditComposer(edit.dataset.incidentId, edit.dataset.editComment);
            return;
        }
        const remove = event.target.closest('[data-delete-comment]');
        if (remove) {
            event.preventDefault();
            event.stopPropagation();
            confirmDeleteComment(remove);
            return;
        }
        const btn = event.target.closest('[data-add-comment]');
        if (!btn) return;
        event.preventDefault();
        event.stopPropagation();
        const id = parseInt(btn.dataset.addComment, 10);
        const parentId = parseInt(btn.dataset.replyTo || '0', 10);
        if (!Number.isNaN(id)) openCommentComposer([id], { parentId: Number.isNaN(parentId) ? 0 : parentId });
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

export function openCommentComposer(ids, options = {}) {
    const parentId = options.parentId || 0;
    openComposer({
        title: parentId ? 'Reply' : 'New comment',
        hint: parentId
            ? 'Your reply stays under this comment. Type @ to mention someone.'
            : 'This starts a new comment. It does not change an existing one. Type @ to mention a teammate.',
        submit: parentId ? 'Post reply' : 'Post comment',
        required: true,
        ids,
        parentId,
        successMessage: null
    });
}

function openEditComposer(incidentId, commentId) {
    openComposer({
        title: 'Edit comment',
        hint: 'Update your note. Type @ to mention someone.',
        submit: 'Save',
        required: true,
        mode: 'edit',
        incidentId,
        commentId,
        html: commentHTML.get(String(commentId)) || ''
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
    if (editor) editor.innerHTML = action.html || '';
    hideMentions();
    const modal = document.getElementById('comment-modal');
    if (modal) modal.hidden = false;
    closePops();
    ensureDirectory();
    editor?.focus();
}

function closeCommentModal() {
    const modal = document.getElementById('comment-modal');
    if (modal) modal.hidden = true;
    pendingAction = null;
    closePops();
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
        if (action.mode === 'edit') {
            if (!hasText) {
                notify('Write a comment first', 'error');
                return;
            }
            const res = await fetchWithAuth(`/api/incidents/${action.incidentId}/comments/${action.commentId}`, {
                method: 'PUT',
                body: JSON.stringify({ blocks })
            });
            if (!res.ok) {
                const message = await res.text();
                throw new Error(message || 'Could not update the comment');
            }
            closeCommentModal();
            notify('Comment updated', 'success');
            await hydrateOpenComments(true);
            return;
        }
        if (hasText) {
            for (const id of action.ids) {
                const res = await fetchWithAuth(`/api/incidents/${id}/comments`, {
                    method: 'POST',
                    body: JSON.stringify({ blocks, parent_id: action.parentId || 0 })
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
        if (el.closest('#unified-report-root')) return true;
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
        nodes.forEach((el) => {
            try {
                renderThread(el, grouped[el.dataset.incidentId] || []);
            } catch (_) { /* one broken thread must not hide the others */ }
        });
    } catch (_) { /* dashboard still usable without the thread */ }
}

function currentUsername() {
    return String(parseJwt(localStorage.getItem('sre-jwt'))?.username || '').toLowerCase();
}

function isMine(comment) {
    return String(comment.author || '').toLowerCase() === currentUsername();
}

function toggleThread(incidentId) {
    const key = String(incidentId || '');
    if (!key) return;
    if (expandedThreads.has(key)) expandedThreads.delete(key);
    else expandedThreads.add(key);
    document.querySelectorAll(`.incident-comments[data-incident-id="${key}"]`).forEach((node) => {
        if (node.__comments) renderThread(node, node.__comments);
    });
}

function confirmDeleteComment(btn) {
    const incidentId = btn.dataset.incidentId;
    const commentId = btn.dataset.deleteComment;
    const isReply = btn.dataset.deleteKind === 'reply';
    const hasReplies = btn.dataset.deleteReplies === '1';
    const message = isReply
        ? 'Delete this reply? This cannot be undone.'
        : hasReplies
            ? 'Delete this comment and its replies? This cannot be undone.'
            : 'Delete this comment? This cannot be undone.';
    showConfirmModal(isReply ? 'Delete reply?' : 'Delete comment?', message, 'danger', async () => {
        try {
            const res = await fetchWithAuth(`/api/incidents/${incidentId}/comments/${commentId}`, { method: 'DELETE' });
            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || 'Could not delete the comment');
            }
            notify(isReply ? 'Reply deleted' : 'Comment deleted', 'success');
            await hydrateOpenComments(true);
        } catch (err) {
            notify(err.message || 'Could not delete the comment', 'error');
        }
    });
}

function renderThread(el, comments) {
    const id = el.dataset.incidentId;
    el.__comments = comments || [];
    const byParent = new Map();
    (comments || []).forEach((comment) => {
        commentHTML.set(String(comment.id), comment.body_html || '');
        const key = comment.parent_id || 0;
        if (!byParent.has(key)) byParent.set(key, []);
        byParent.get(key).push(comment);
    });
    const roots = byParent.get(0) || [];
    const expanded = expandedThreads.has(String(id));
    const clip = !expanded && roots.length > 2;
    const renderList = (parentId, nested) => {
        const list = byParent.get(parentId) || [];
        return list.map((comment) => {
            const replies = byParent.get(comment.id) || [];
            const mine = isMine(comment);
            const actions = [];
            if (!nested && !mine) {
                actions.push(`<button type="button" class="btn btn-secondary btn-sm comment-reply-btn" data-add-comment="${id}" data-reply-to="${comment.id}">Reply</button>`);
            }
            if (mine) {
                actions.push(`<button type="button" class="btn btn-secondary btn-sm" data-edit-comment="${comment.id}" data-incident-id="${id}">Edit</button>`);
                actions.push(`<button type="button" class="btn btn-secondary btn-sm btn-danger" data-delete-comment="${comment.id}" data-incident-id="${id}" data-delete-kind="${nested ? 'reply' : 'comment'}" data-delete-replies="${replies.length ? '1' : '0'}">Delete</button>`);
            }
            return `
        <article class="comment-item${nested ? ' comment-item--reply' : ''}" data-comment-id="${escapeHtml(comment.id)}">
            ${avatarMarkup({
                avatar_kind: comment.avatar_kind,
                avatar_value: comment.avatar_value,
                fullname: comment.fullname || comment.author,
                username: comment.author
            }, nested ? 26 : 32)}
            <div class="comment-item__body">
                <header class="comment-item__meta">
                    <strong>${escapeHtml(comment.fullname || comment.author)}</strong>
                    <time>${escapeHtml(comment.created_at || '')}</time>
                </header>
                <div class="comment-item__text">${comment.body_html || ''}</div>
                ${actions.length ? `<div class="comment-item__actions">${actions.join('')}</div>` : ''}
                <div class="comment-replies">${renderList(comment.id, true)}</div>
            </div>
        </article>`;
        }).join('');
    };
    const items = renderList(0, false);
    const more = clip
        ? `<button type="button" class="comment-see-more" data-thread-more="${id}">See more</button>`
        : (expanded && roots.length > 2 ? `<button type="button" class="comment-see-more" data-thread-more="${id}">See less</button>` : '');
    const total = (comments || []).length;
    el.innerHTML = `
        <button type="button" class="btn btn-primary btn-sm comment-new-btn" data-add-comment="${id}">New comment</button>
        <div class="comment-thread${clip ? ' comment-thread--clip' : ''}">${items || '<p class="comment-empty">No notes yet.</p>'}</div>
        ${more}`;
    const badge = el.closest('.ur-test-block')?.querySelector('.ur-comment-count');
    if (badge) {
        badge.hidden = total === 0;
        badge.textContent = total === 1 ? '1 comment' : `${total} comments`;
    }
}

function applyCommand(cmd) {
    focusEditor();
    if (cmd === 'bold') document.execCommand('bold');
    if (cmd === 'italic') document.execCommand('italic');
    if (cmd === 'underline') document.execCommand('underline');
    if (cmd === 'ul') document.execCommand('insertUnorderedList');
    if (cmd === 'ol') document.execCommand('insertOrderedList');
    if (cmd === 'code') wrapSelection((text) => `<code>${text}</code>`);
    if (cmd === 'pre') {
        const editor = document.getElementById('rte-editor');
        const selected = window.getSelection()?.toString() || 'code';
        document.execCommand('insertHTML', false, `<pre><code>${escapeHtml(selected)}</code></pre>`);
        editor?.focus();
    }
}

function togglePop(name) {
    const pop = document.getElementById(name === 'emoji' ? 'rte-emoji-pop' : 'rte-color-pop');
    const opening = !!pop?.hidden;
    closePops();
    if (!opening || !pop) return;
    pop.hidden = false;
    document.querySelector(`[data-pop="${name}"]`)?.setAttribute('aria-expanded', 'true');
}

function closePops() {
    document.querySelectorAll('.rte-pop').forEach((pop) => {
        pop.hidden = true;
    });
    document.querySelectorAll('[data-pop]').forEach((btn) => {
        btn.setAttribute('aria-expanded', 'false');
    });
}

function applyColor(color) {
    focusEditor();
    document.execCommand('styleWithCSS', false, false);
    if (color === 'default') {
        clearTextColor();
    } else if (/^#[0-9a-fA-F]{6}$/.test(color || '')) {
        document.execCommand('foreColor', false, color);
    } else {
        return;
    }
    const ink = color === 'default' ? '' : color.toLowerCase();
    document.querySelectorAll('.rte-swatch').forEach((swatch) => {
        swatch.classList.toggle('is-active', swatch.dataset.color === (ink || 'default'));
    });
    const mark = document.querySelector('.rte-color-mark');
    if (mark) mark.style.setProperty('--rte-ink', ink || 'var(--text-main)');
}

function clearTextColor() {
    const editor = document.getElementById('rte-editor');
    document.execCommand('foreColor', false, '#010101');
    if (!editor) return;
    editor.querySelectorAll('font, span').forEach((node) => {
        if (!node.isConnected || paintedColor(node) !== '#010101') return;
        let parent = node.parentElement;
        while (parent && parent !== editor && (parent.nodeName === 'FONT' || parent.nodeName === 'SPAN')) {
            if (parent.classList.contains('mention') || parent.dataset.user) break;
            if (parent.nodeName === 'FONT') parent.removeAttribute('color');
            else parent.style.color = '';
            const next = parent.parentElement;
            if (!parent.getAttribute('color') && !parent.getAttribute('style')?.trim() && !parent.className) unwrapNode(parent);
            parent = next;
        }
        if (!node.isConnected) return;
        if (node.nodeName === 'FONT') node.removeAttribute('color');
        else node.style.color = '';
        if (!node.classList.contains('mention') && !node.dataset.user && !node.getAttribute('color') && !node.getAttribute('style')?.trim() && !node.className) {
            unwrapNode(node);
        }
    });
}

function paintedColor(node) {
    if (node.nodeName === 'FONT') return normalizeColor(node.color);
    return normalizeColor(node.style?.color);
}

function unwrapNode(node) {
    const parent = node.parentNode;
    if (!parent) return;
    while (node.firstChild) parent.insertBefore(node.firstChild, node);
    parent.removeChild(node);
}

function insertEmoji(emoji) {
    if (!emoji) return;
    focusEditor();
    document.execCommand('insertText', false, emoji);
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
    const me = currentUsername();
    const matches = directory.filter((user) => {
        if (me && String(user.username).toLowerCase() === me) return false;
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
            <span><strong>${escapeHtml(mentionLabel(user))}</strong></span>
        </button>`).join('');
}

function hideMentions() {
    const menu = document.getElementById('rte-mentions');
    if (menu) {
        menu.hidden = true;
        menu.innerHTML = '';
    }
}

function mentionLabel(user) {
    const name = String(user?.fullname || '').trim();
    return name || user?.username || '';
}

export function revealComment(commentId) {
    const id = String(commentId || '');
    if (!/^\d+$/.test(id)) return;
    document.querySelectorAll('.incident-comments').forEach((node) => {
        const comments = node.__comments || [];
        if (!comments.some((comment) => String(comment.id) === id)) return;
        const key = String(node.dataset.incidentId || '');
        if (key && !expandedThreads.has(key)) {
            expandedThreads.add(key);
            renderThread(node, comments);
        }
    });
    const target = document.querySelector(`.comment-item[data-comment-id="${id}"]`);
    if (!target) return;
    target.classList.add('comment-item--focus');
    target.scrollIntoView({ behavior: 'smooth', block: 'center' });
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
    const person = directory.find((user) => String(user.username).toLowerCase() === String(username).toLowerCase());
    span.textContent = `@${mentionLabel(person) || username}`;
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
    const walk = (node, marks, into) => {
        if (node.nodeType === Node.TEXT_NODE) {
            if (node.textContent) into.push({ t: node.textContent, ...marks });
            return;
        }
        if (node.nodeType !== Node.ELEMENT_NODE) return;
        if (node.nodeName === 'UL' || node.nodeName === 'OL') {
            const kind = node.nodeName === 'OL' ? 'ol' : 'ul';
            node.childNodes.forEach((child) => {
                if (child.nodeName !== 'LI') return;
                const item = [];
                child.childNodes.forEach((grand) => walk(grand, marks, item));
                if (!item.some((block) => (block.t && block.t.trim()) || block.mention || (block.pre && block.pre.trim()))) return;
                item[0].item = true;
                item.forEach((block) => {
                    block.li = kind;
                    into.push(block);
                });
            });
            return;
        }
        if (node.nodeName === 'BR') {
            into.push({ br: true });
            return;
        }
        if (node.classList?.contains('mention')) {
            into.push({ mention: node.dataset.user || '' });
            return;
        }
        if (node.nodeName === 'PRE') {
            into.push({ pre: node.innerText || '' });
            return;
        }
        if (node.nodeName === 'DIV' || node.nodeName === 'P') {
            if (into.length > 0) into.push({ br: true });
            const blockMarks = { ...marks };
            node.childNodes.forEach((child) => walk(child, blockMarks, into));
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
        node.childNodes.forEach((child) => walk(child, next, into));
    };
    root.childNodes.forEach((child) => walk(child, {}, blocks));
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
        if (block.li === 'ul' || block.li === 'ol') clean.li = block.li;
        if (block.item) clean.item = true;
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
