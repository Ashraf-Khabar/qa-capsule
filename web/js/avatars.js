/**
 * Twelve original profile avatars plus a photo fallback.
 * Preset ids match the server allow-list (a01–a12).
 */

const PRESETS = [
    { id: 'a01', skin: '#F3C7A5', hair: '#2B2118', shirt: '#0F766E', hairKind: 'short' },
    { id: 'a02', skin: '#E7B48A', hair: '#1F2937', shirt: '#BE123C', hairKind: 'bun' },
    { id: 'a03', skin: '#C68642', hair: '#111827', shirt: '#1D4ED8', hairKind: 'curls', glasses: true },
    { id: 'a04', skin: '#8D5524', hair: '#1C1917', shirt: '#B45309', hairKind: 'beard' },
    { id: 'a05', skin: '#F6D3B3', hair: '#7C3AED', shirt: '#6D28D9', hairKind: 'bob' },
    { id: 'a06', skin: '#3F2A1D', hair: '#111827', shirt: '#15803D', hairKind: 'cap' },
    { id: 'a07', skin: '#F1C6A8', hair: '#9A3412', shirt: '#E11D48', hairKind: 'long' },
    { id: 'a08', skin: '#D9A066', hair: '#0F172A', shirt: '#0369A1', hairKind: 'headphones' },
    { id: 'a09', skin: '#A86B3C', hair: '#292524', shirt: '#334155', hairKind: 'fade' },
    { id: 'a10', skin: '#F8D7C4', hair: '#CA8A04', shirt: '#65A30D', hairKind: 'ponytail' },
    { id: 'a11', skin: '#6B3A22', hair: '#0C0A09', shirt: '#1E3A8A', hairKind: 'hoodie' },
    { id: 'a12', skin: '#EDC3A1', hair: '#44403C', shirt: '#C2410C', hairKind: 'waves' }
];

function hairMarkup(kind, color) {
    switch (kind) {
        case 'bun':
            return `<circle cx="32" cy="16" r="6" fill="${color}"/><path d="M18 30c2-12 26-12 28 0v4H18z" fill="${color}"/>`;
        case 'curls':
            return `<circle cx="20" cy="24" r="6" fill="${color}"/><circle cx="32" cy="16" r="7" fill="${color}"/><circle cx="44" cy="24" r="6" fill="${color}"/>`;
        case 'beard':
            return `<path d="M18 28c1-10 27-10 28 0v3H18z" fill="${color}"/><path d="M22 42c2 8 18 8 20 0v-4H22z" fill="${color}"/>`;
        case 'bob':
            return `<path d="M16 34c1-16 31-16 32 0v8H16z" fill="${color}"/>`;
        case 'cap':
            return `<path d="M16 30c2-12 30-14 36-2 2 4-2 6-6 5H18c-3 0-4-1-2-3z" fill="${color}"/><path d="M14 32h28v3H14z" fill="${color}"/>`;
        case 'long':
            return `<path d="M17 28c2-14 28-14 30 0l2 22h-8l-1-16h-16l-1 16h-8z" fill="${color}"/>`;
        case 'headphones':
            return `<path d="M18 28c2-12 26-12 28 0v2H18z" fill="${color}"/><path d="M14 30h6v12h-6a4 4 0 0 1-4-4v-4a4 4 0 0 1 4-4z" fill="#0f172a"/><path d="M44 30h6a4 4 0 0 1 4 4v4a4 4 0 0 1-4 4h-6z" fill="#0f172a"/>`;
        case 'fade':
            return `<path d="M20 30c1-8 23-10 26-1 1 2-1 3-3 3H22c-2 0-3-1-2-2z" fill="${color}"/>`;
        case 'ponytail':
            return `<path d="M18 30c2-12 26-12 28 0v3H18z" fill="${color}"/><path d="M44 24c8 2 10 14 4 20-2-6-2-12-6-16z" fill="${color}"/>`;
        case 'hoodie':
            return `<path d="M16 36c0-14 32-16 36 2v16H16z" fill="${color}" opacity=".35"/><path d="M20 28c2-10 22-10 24 0v2H20z" fill="${color}"/>`;
        case 'waves':
            return `<path d="M16 32c4-14 28-16 34 0-4 2-6 2-8-2-2 6-8 8-14 4-4 4-8 4-12-2z" fill="${color}"/>`;
        default:
            return `<path d="M18 30c2-12 26-12 28 0v4H18z" fill="${color}"/>`;
    }
}

function presetSVG(preset) {
    const glasses = preset.glasses
        ? `<circle cx="26" cy="36" r="4.2" fill="none" stroke="#0f172a" stroke-width="1.4"/><circle cx="38" cy="36" r="4.2" fill="none" stroke="#0f172a" stroke-width="1.4"/><path d="M30.2 36h3.6" stroke="#0f172a" stroke-width="1.4"/>`
        : '';
    return `<svg viewBox="0 0 64 64" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
        <circle cx="32" cy="32" r="32" fill="${preset.shirt}"/>
        <circle cx="32" cy="40" r="18" fill="${preset.skin}"/>
        ${hairMarkup(preset.hairKind, preset.hair)}
        <circle cx="26" cy="36" r="1.7" fill="#1c1917"/>
        <circle cx="38" cy="36" r="1.7" fill="#1c1917"/>
        <path d="M28 43c1.5 2.2 6.5 2.2 8 0" fill="none" stroke="#1c1917" stroke-width="1.4" stroke-linecap="round"/>
        ${glasses}
    </svg>`;
}

const PRESET_MARKUP = Object.fromEntries(PRESETS.map((preset) => [preset.id, presetSVG(preset)]));

export function listAvatarPresets() {
    return PRESETS.map((preset) => ({ id: preset.id, svg: PRESET_MARKUP[preset.id] }));
}

export function avatarMarkup(person, size) {
    const px = size || 32;
    const kind = person?.avatar_kind || person?.kind || '';
    const value = person?.avatar_value || person?.value || '';
    const name = person?.fullname || person?.name || person?.username || '';
    if (kind === 'photo' && typeof value === 'string' && value.startsWith('data:image/')) {
        return `<img class="avatar-ball" width="${px}" height="${px}" alt="" src="${value}">`;
    }
    if (kind === 'preset' && PRESET_MARKUP[value]) {
        return `<span class="avatar-ball avatar-ball--svg" style="width:${px}px;height:${px}px">${PRESET_MARKUP[value]}</span>`;
    }
    const letters = initials(name);
    return `<span class="avatar-ball avatar-ball--initials" style="width:${px}px;height:${px}px;font-size:${Math.max(10, Math.round(px * 0.36))}px">${letters}</span>`;
}

function initials(name) {
    const cleaned = String(name || '').trim();
    if (!cleaned) return '?';
    const parts = cleaned.split(/\s+/).filter(Boolean);
    if (parts.length >= 2) return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
    return cleaned.slice(0, 2).toUpperCase();
}

export function renderAvatarPicker(container, selectedId, onPick) {
    if (!container) return;
    container.innerHTML = PRESETS.map((preset) => `
        <button type="button" class="avatar-choice${preset.id === selectedId ? ' is-selected' : ''}" data-avatar="${preset.id}" aria-label="Avatar ${preset.id}">
            ${PRESET_MARKUP[preset.id]}
        </button>`).join('');
    container.querySelectorAll('.avatar-choice').forEach((btn) => {
        btn.addEventListener('click', () => onPick(btn.dataset.avatar));
    });
}

export function resizeProfilePhoto(file) {
    return new Promise((resolve, reject) => {
        if (!file || !/^image\/(png|jpeg|webp)$/.test(file.type)) {
            reject(new Error('Use a PNG, JPEG, or WebP image'));
            return;
        }
        const reader = new FileReader();
        reader.onerror = () => reject(new Error('Could not read the image'));
        reader.onload = () => {
            const img = new Image();
            img.onload = () => {
                const canvas = document.createElement('canvas');
                const size = 128;
                canvas.width = size;
                canvas.height = size;
                const ctx = canvas.getContext('2d');
                const scale = Math.max(size / img.width, size / img.height);
                const w = img.width * scale;
                const h = img.height * scale;
                ctx.drawImage(img, (size - w) / 2, (size - h) / 2, w, h);
                resolve(canvas.toDataURL('image/jpeg', 0.85));
            };
            img.onerror = () => reject(new Error('Could not read the image'));
            img.src = String(reader.result);
        };
        reader.readAsDataURL(file);
    });
}
