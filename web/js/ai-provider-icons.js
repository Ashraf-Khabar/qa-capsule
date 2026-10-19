/**
 * AI provider metadata and brand logos for the configuration picker.
 * Logos use each vendor's real mark and brand colors, drawn on a white chip.
 */

export const AI_PROVIDERS = [
    { id: 'disabled', label: 'Disabled', defaultModel: '', defaultBaseUrl: '', defaultKeyEnv: '', models: [] },
    provider('openai', 'OpenAI', 'gpt-4o-mini', 'https://api.openai.com/v1', 'OPENAI_API_KEY', [
        ['gpt-4o', 'GPT-4o'], ['gpt-4o-mini', 'GPT-4o mini'], ['gpt-4.1', 'GPT-4.1'], ['o3-mini', 'o3-mini'], ['o1', 'o1']
    ]),
    provider('anthropic', 'Anthropic', 'claude-3-5-haiku-20241022', 'https://api.anthropic.com', 'ANTHROPIC_API_KEY', [
        ['claude-opus-4-5', 'Claude Opus 4.5'], ['claude-sonnet-4-5', 'Claude Sonnet 4.5'],
        ['claude-3-5-sonnet-20241022', 'Claude 3.5 Sonnet'], ['claude-3-5-haiku-20241022', 'Claude 3.5 Haiku']
    ]),
    provider('gemini', 'Google Gemini', 'gemini-2.0-flash', 'https://generativelanguage.googleapis.com/v1beta', 'GEMINI_API_KEY', [
        ['gemini-2.5-pro', 'Gemini 2.5 Pro'], ['gemini-2.0-flash', 'Gemini 2.0 Flash'], ['gemini-1.5-pro', 'Gemini 1.5 Pro']
    ]),
    provider('mistral', 'Mistral', 'mistral-small-latest', 'https://api.mistral.ai/v1', 'MISTRAL_API_KEY', [
        ['mistral-large-latest', 'Mistral Large'], ['mistral-small-latest', 'Mistral Small'], ['codestral-latest', 'Codestral']
    ]),
    provider('groq', 'Groq', 'llama-3.1-8b-instant', 'https://api.groq.com/openai/v1', 'GROQ_API_KEY', [
        ['llama-3.3-70b-versatile', 'Llama 3.3 70B'], ['llama-3.1-8b-instant', 'Llama 3.1 8B Instant'], ['gemma2-9b-it', 'Gemma 2 9B']
    ]),
    provider('openrouter', 'OpenRouter', 'openai/gpt-4o-mini', 'https://openrouter.ai/api/v1', 'OPENROUTER_API_KEY', [
        ['openai/gpt-4o-mini', 'OpenAI GPT-4o mini'], ['anthropic/claude-3.5-sonnet', 'Claude 3.5 Sonnet'], ['google/gemini-2.0-flash-001', 'Gemini 2.0 Flash']
    ]),
    provider('azure', 'Azure OpenAI', 'gpt-4o-mini', 'https://YOUR-RESOURCE.openai.azure.com/openai/deployments/YOUR-DEPLOYMENT', 'AZURE_OPENAI_API_KEY', [
        ['gpt-4o', 'GPT-4o'], ['gpt-4o-mini', 'GPT-4o mini'], ['gpt-4-turbo', 'GPT-4 Turbo']
    ]),
    provider('ollama', 'Ollama', 'llama3.2', 'http://localhost:11434', '', [
        ['llama3.2', 'Llama 3.2'], ['llama3.1', 'Llama 3.1'], ['mistral', 'Mistral'], ['qwen2.5', 'Qwen 2.5']
    ], true),
    provider('deepseek', 'DeepSeek', 'deepseek-chat', 'https://api.deepseek.com/v1', 'DEEPSEEK_API_KEY', [
        ['deepseek-chat', 'DeepSeek Chat'], ['deepseek-reasoner', 'DeepSeek Reasoner']
    ]),
    provider('xai', 'xAI', 'grok-3', 'https://api.x.ai/v1', 'XAI_API_KEY', [
        ['grok-3', 'Grok 3'], ['grok-3-mini', 'Grok 3 Mini'], ['grok-2', 'Grok 2']
    ]),
    provider('cohere', 'Cohere', 'command-r-plus', 'https://api.cohere.ai/compatibility/v1', 'COHERE_API_KEY', [
        ['command-a-03-2025', 'Command A'], ['command-r-plus', 'Command R+'], ['command-r', 'Command R']
    ]),
    provider('together', 'Together AI', 'meta-llama/Meta-Llama-3.1-8B-Instruct-Turbo', 'https://api.together.xyz/v1', 'TOGETHER_API_KEY', [
        ['meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo', 'Llama 3.1 70B'], ['meta-llama/Meta-Llama-3.1-8B-Instruct-Turbo', 'Llama 3.1 8B']
    ]),
    provider('fireworks', 'Fireworks', 'accounts/fireworks/models/llama-v3p1-8b-instruct', 'https://api.fireworks.ai/inference/v1', 'FIREWORKS_API_KEY', [
        ['accounts/fireworks/models/llama-v3p1-70b-instruct', 'Llama 3.1 70B'], ['accounts/fireworks/models/llama-v3p1-8b-instruct', 'Llama 3.1 8B']
    ]),
    provider('perplexity', 'Perplexity', 'sonar', 'https://api.perplexity.ai/chat/completions', 'PERPLEXITY_API_KEY', [
        ['sonar-pro', 'Sonar Pro'], ['sonar', 'Sonar'], ['sonar-reasoning-pro', 'Sonar Reasoning Pro']
    ]),
    provider('cerebras', 'Cerebras', 'llama3.1-8b', 'https://api.cerebras.ai/v1', 'CEREBRAS_API_KEY', [
        ['llama-3.3-70b', 'Llama 3.3 70B'], ['llama3.1-8b', 'Llama 3.1 8B']
    ]),
    provider('sambanova', 'SambaNova', 'Meta-Llama-3.1-8B-Instruct', 'https://api.sambanova.ai/v1', 'SAMBANOVA_API_KEY', [
        ['Meta-Llama-3.1-405B-Instruct', 'Llama 3.1 405B'], ['Meta-Llama-3.1-8B-Instruct', 'Llama 3.1 8B']
    ]),
    provider('nvidia', 'NVIDIA', 'meta/llama-3.1-8b-instruct', 'https://integrate.api.nvidia.com/v1', 'NVIDIA_API_KEY', [
        ['meta/llama-3.1-70b-instruct', 'Llama 3.1 70B'], ['meta/llama-3.1-8b-instruct', 'Llama 3.1 8B']
    ]),
    provider('huggingface', 'Hugging Face', 'meta-llama/Llama-3.1-8B-Instruct', 'https://router.huggingface.co/v1', 'HF_TOKEN', [
        ['meta-llama/Llama-3.1-8B-Instruct', 'Llama 3.1 8B'], ['Qwen/Qwen2.5-7B-Instruct', 'Qwen 2.5 7B']
    ]),
    provider('github', 'GitHub Models', 'gpt-4o-mini', 'https://models.inference.ai.azure.com', 'GITHUB_TOKEN', [
        ['gpt-4o', 'GPT-4o'], ['gpt-4o-mini', 'GPT-4o mini'], ['Meta-Llama-3.1-70B-Instruct', 'Llama 3.1 70B']
    ]),
    provider('qwen', 'Qwen', 'qwen-plus', 'https://dashscope.aliyuncs.com/compatible-mode/v1', 'DASHSCOPE_API_KEY', [
        ['qwen-max', 'Qwen Max'], ['qwen-plus', 'Qwen Plus'], ['qwen-turbo', 'Qwen Turbo']
    ]),
    provider('moonshot', 'Moonshot', 'moonshot-v1-8k', 'https://api.moonshot.ai/v1', 'MOONSHOT_API_KEY', [
        ['moonshot-v1-128k', 'Kimi 128K'], ['moonshot-v1-32k', 'Kimi 32K'], ['moonshot-v1-8k', 'Kimi 8K']
    ]),
    provider('zhipu', 'Zhipu', 'glm-4-flash', 'https://open.bigmodel.cn/api/paas/v4', 'ZHIPU_API_KEY', [
        ['glm-4-plus', 'GLM-4 Plus'], ['glm-4-flash', 'GLM-4 Flash']
    ]),
    provider('minimax', 'MiniMax', 'MiniMax-Text-01', 'https://api.minimax.io/v1', 'MINIMAX_API_KEY', [
        ['MiniMax-Text-01', 'MiniMax Text'], ['abab6.5s-chat', 'ABAB 6.5s']
    ]),
    provider('lmstudio', 'LM Studio', 'local-model', 'http://localhost:1234/v1', '', [
        ['local-model', 'Local model']
    ], true),
    provider('cloudflare', 'Cloudflare', '@cf/meta/llama-3.1-8b-instruct', 'https://api.cloudflare.com/client/v4/accounts/ACCOUNT_ID/ai/v1', 'CLOUDFLARE_API_TOKEN', [
        ['@cf/meta/llama-3.1-8b-instruct', 'Llama 3.1 8B'], ['@cf/meta/llama-3.3-70b-instruct-fp8-fast', 'Llama 3.3 70B']
    ]),
    provider('custom', 'Custom endpoint', 'gpt-4o-mini', '', 'OPENAI_API_KEY', [
        ['gpt-4o-mini', 'gpt-4o-mini']
    ])
];

function provider(id, label, defaultModel, defaultBaseUrl, defaultKeyEnv, models, keyOptional) {
    return {
        id, label, defaultModel, defaultBaseUrl, defaultKeyEnv, keyOptional: !!keyOptional,
        models: models.map(([modelId, modelLabel]) => ({ id: modelId, label: modelLabel }))
    };
}

export function iconAiInsight() {
    return `<svg class="ai-svg ai-svg--insight" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <path d="M12 2a4 4 0 0 1 4 4v1a3 3 0 0 1 3 3v1a2 2 0 0 1-2 2h-1.2a1 1 0 0 0-.8.4l-.6.8a1 1 0 0 1-1.6 0l-.6-.8a1 1 0 0 0-.8-.4H7a2 2 0 0 1-2-2v-1a3 3 0 0 1 3-3V6a4 4 0 0 1 4-4z"/>
        <path d="M9 18h6"/><path d="M10 22h4"/>
        <circle cx="9" cy="8" r="1" fill="currentColor" stroke="none"/><circle cx="15" cy="8" r="1" fill="currentColor" stroke="none"/>
    </svg>`;
}

export function logoForProvider(providerId) {
    return LOGOS[providerId] || LOGOS.disabled;
}

export function getProviderMeta(providerId) {
    return AI_PROVIDERS.find((item) => item.id === providerId) || AI_PROVIDERS[0];
}

function logo(body) {
    return `<svg class="ai-svg ai-svg--logo" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">${body}</svg>`;
}

const LOGOS = {
    disabled: logo('<circle cx="12" cy="12" r="8" fill="none" stroke="#94a3b8" stroke-width="1.6"/><path d="M8 16 16 8" stroke="#94a3b8" stroke-width="1.6"/>'),
    openai: logo('<path fill="#000" d="M22.282 9.821a5.985 5.985 0 0 0-.516-4.938 6.046 6.046 0 0 0-6.51-2.9A6.065 6.065 0 0 0 4.981 4.18a5.985 5.985 0 0 0-3.998 2.9 6.046 6.046 0 0 0 .743 7.097 5.98 5.98 0 0 0 .51 4.911 6.051 6.051 0 0 0 6.516 2.9A5.985 5.985 0 0 0 13.26 24a6.056 6.056 0 0 0 5.772-4.206 5.99 5.99 0 0 0 3.997-2.9 6.056 6.056 0 0 0-.747-7.073zM13.26 22.43a4.476 4.476 0 0 1-2.876-1.04l.141-.081 4.779-2.758a.795.795 0 0 0 .392-.681v-6.737l2.02 1.168a.071.071 0 0 1 .038.052v5.583a4.504 4.504 0 0 1-4.494 4.494zM3.6 18.305a4.47 4.47 0 0 1-.535-3.014l.142.085 4.783 2.759a.771.771 0 0 0 .78 0l5.843-3.369v2.332a.08.08 0 0 1-.033.062L9.74 19.95a4.5 4.5 0 0 1-6.14-1.646zM2.34 7.896a4.485 4.485 0 0 1 2.366-1.973V11.6a.766.766 0 0 0 .388.676l5.815 3.355-2.02 1.168a.076.076 0 0 1-.071 0l-4.83-2.786A4.504 4.504 0 0 1 2.34 7.872zm16.597 3.855-5.833-3.387 2.015-1.164a.076.076 0 0 1 .071 0l4.83 2.791a4.494 4.494 0 0 1-.676 8.105v-5.678a.79.79 0 0 0-.407-.667zm2.01-3.023-.142-.085-4.774-2.782a.776.776 0 0 0-.785 0L9.409 9.229V6.897a.066.066 0 0 1 .028-.061l4.83-2.787a4.5 4.5 0 0 1 6.68 4.66zm-12.64 4.135-2.02-1.164a.08.08 0 0 1-.038-.057V6.843a4.5 4.5 0 0 1 7.375-3.453l-.142.08-4.778 2.758a.795.795 0 0 0-.393.681zm1.49-3.208 2.602-1.5 2.607 1.5v3l-2.597 1.5-2.607-1.5z"/>'),
    anthropic: logo('<path fill="#D97757" d="M16.8 2.2 13.4 12h-2.2L14.6 2.2zm-6.2 0L7.2 12H5L8.4 2.2zM19.2 2.2 22 12h-2.3L17 2.2z"/>'),
    gemini: logo('<defs><linearGradient id="sre-gemini" x1="0" y1="1" x2="1" y2="0"><stop offset="0" stop-color="#1a73e8"/><stop offset=".55" stop-color="#8ab4f8"/><stop offset="1" stop-color="#f4397c"/></linearGradient></defs><path fill="url(#sre-gemini)" d="M12 1.2c.2 4.6 2.4 7.2 6.8 8.4-4.4 1.1-6.6 3.8-6.8 8.4-.2-4.6-2.4-7.3-6.8-8.4C9.6 8.4 11.8 5.8 12 1.2z"/>'),
    mistral: logo('<g fill="#F54E00"><path d="M3 3h4v4H3zm6 0h4v4H9zm6 0h4v4h-4zM3 9h4v4H3zm12 0h4v4h-4zM9 15h4v4H9zm6 0h4v4h-4z"/></g>'),
    groq: logo('<path fill="#F55036" d="M12 2.2A9.8 9.8 0 1 0 21.8 12 9.8 9.8 0 0 0 12 2.2zm0 3.2A6.6 6.6 0 0 1 18.2 15h-2.3A4.4 4.4 0 1 0 12 16.4V18a7.6 7.6 0 0 1 0-15.2z"/><path fill="#111" d="M11 10.2h4.2v1.7H11z"/>'),
    openrouter: logo('<path fill="#6467F2" d="M4 7.5 12 3l8 4.5v9L12 21l-8-4.5z"/><path fill="#fff" d="M12 8.2 8.2 12 12 15.8 15.8 12z"/>'),
    azure: logo('<path fill="#0078D4" d="M13.2 3.2 8.4 18.6h3.2l.8-2.4h4.7L18 18.6h3.1L16.4 3.2zm.9 10.2 1.5-4.6 1.5 4.6z"/><path fill="#50E6FF" d="M3 16.8 8.6 4.2h2.2L5.2 16.8z" opacity=".9"/>'),
    ollama: logo('<path fill="#000" d="M8 4.5a2.4 2.4 0 0 0-1.6 4.1C5 9.4 4 11 4 13v6.5h2.2V13c0-1.2.8-2.2 2-2.4V20h2.1v-6.2c.5.2 1.1.3 1.7.3s1.2-.1 1.7-.3V20H16v-9.4c1.2.2 2 1.2 2 2.4v6.5H20.2V13c0-2-1-3.6-2.4-4.4A2.4 2.4 0 0 0 16 4.5a2.4 2.4 0 0 0-2 .9 2.4 2.4 0 0 0-2-.9 2.4 2.4 0 0 0-2 .9 2.4 2.4 0 0 0-2-.9z"/>'),
    deepseek: logo('<path fill="#4D6BFE" d="M4 14c2-6 6-8 8-8 4 0 6 3 8 3 2 0 2-2 2-2-3 5-6 7-9 7-3 0-6-2-9-0z"/><circle cx="16.5" cy="8.5" r="1.3" fill="#4D6BFE"/>'),
    xai: logo('<path fill="#000" d="M4 5h3.2l4.8 6.4L16.8 5H20l-6.4 8.2L20.2 19H17l-5-6.6L7 19H3.8l6.6-8.4z"/>'),
    cohere: logo('<path fill="#FF7759" d="M8 4h5.2a6.2 6.2 0 0 1 0 12.4H8z"/><path fill="#39594D" d="M8 7.2h4.4a3 3 0 0 1 0 6H8z"/>'),
    together: logo('<path fill="#000" d="M4 16.5 12 4l8 12.5h-3.2L12 9.2 7.2 16.5z"/><path fill="#000" d="M7 19h10v2H7z"/>'),
    fireworks: logo('<path fill="#6D28D9" d="M12 2 14 8l6 .4-4.6 3.8L16.8 18 12 14.8 7.2 18l1.4-5.8L4 8.4 10 8z"/>'),
    perplexity: logo('<path fill="#20808D" d="M12 2.2 13.4 8 18.8 5.6 15.2 10.2 21 12l-5.8 1.8 3.6 4.6-5.4-2.4L12 21.8 10.6 16l-5.4 2.4 3.6-4.6L3 12l5.8-1.8L5.2 5.6 10.6 8z"/>'),
    cerebras: logo('<path fill="#F15A29" d="M7 5h8.2A5.8 5.8 0 0 1 15 16.6L9.2 19V5z"/><path fill="#fff" d="M9.4 8.2h4.2a2.6 2.6 0 0 1 0 5.2H9.4z"/>'),
    sambanova: logo('<path fill="#EE3A43" d="M3 16c4-8 8-8 9-8s5 2 9 8c-4-3-6-3-9-3s-5 0-9 3z"/><path fill="#111" d="M3 18h18v1.6H3z"/>'),
    nvidia: logo('<path fill="#76B900" d="M4 12c1.2-4 4.2-6.4 8-6.4s6.8 2.4 8 6.4c-1.2 4-4.2 6.4-8 6.4S5.2 16 4 12z"/><circle cx="12" cy="12" r="3.1" fill="#fff"/><circle cx="12" cy="12" r="1.5" fill="#76B900"/>'),
    huggingface: logo('<circle cx="12" cy="12" r="9" fill="#FFD21E"/><circle cx="3.6" cy="13" r="1.5" fill="#FFD21E"/><circle cx="20.4" cy="13" r="1.5" fill="#FFD21E"/><circle cx="9" cy="11" r="1.1" fill="#111"/><circle cx="15" cy="11" r="1.1" fill="#111"/><path d="M8.6 14.4c.9 1.6 5.9 1.6 6.8 0" fill="none" stroke="#111" stroke-width="1.2" stroke-linecap="round"/>'),
    github: logo('<path fill="#181717" d="M12 .3C5.4.3 0 5.7 0 12.3c0 5.3 3.4 9.8 8.2 11.4.6.1.8-.3.8-.6v-2.1c-3.3.7-4-1.6-4-1.6-.5-1.4-1.3-1.8-1.3-1.8-1.1-.7.1-.7.1-.7 1.2.1 1.8 1.2 1.8 1.2 1.1 1.8 2.8 1.3 3.5 1 .1-.8.4-1.3.8-1.6-2.7-.3-5.5-1.3-5.5-5.9 0-1.3.5-2.4 1.2-3.2-.1-.3-.5-1.5.1-3.2 0 0 1-.3 3.3 1.2a11.5 11.5 0 0 1 6 0C17.3 4.7 18.3 5 18.3 5c.6 1.7.2 2.9.1 3.2.8.8 1.2 1.9 1.2 3.2 0 4.6-2.8 5.6-5.5 5.9.4.4.8 1.1.8 2.2v3.3c0 .3.2.7.8.6 4.8-1.6 8.2-6.1 8.2-11.4C24 5.7 18.6.3 12 .3z"/>'),
    qwen: logo('<path fill="#615CED" d="M12 2.4 14.2 8l6 .6-4.6 3.6 1.4 5.8L12 14.8 6.9 18l1.4-5.8L3.8 8.6 9.8 8z"/><circle cx="12" cy="12" r="2" fill="#fff"/>'),
    moonshot: logo('<circle cx="12" cy="12" r="8" fill="#111"/><path fill="#fff" d="M14.2 6.2a6 6 0 1 0 0 11.6 5 5 0 0 1 0-11.6z"/>'),
    zhipu: logo('<path fill="#2563EB" d="M4 7.5 12 3l8 4.5v9L12 21l-8-4.5z"/><path fill="#fff" d="M8 9.2h8l-3.2 2.2H16L9.2 16v-2.6L12.4 12H8z"/>'),
    minimax: logo('<path fill="#111" d="M3 17 7.2 6h2.2L12 13.2 14.6 6h2.2L21 17h-2.3l-2-7.2L14.2 17h-2.2l-2.5-7.2L7.3 17z"/>'),
    lmstudio: logo('<rect x="3" y="3" width="18" height="18" rx="5" fill="#7C5CFF"/><path d="M7 14.5c1.2-3 2.2-3 3.2 0 .8 2.2 1.6 2.2 2.4 0 1-3 2-3 3.2 0" fill="none" stroke="#fff" stroke-width="1.6" stroke-linecap="round"/>'),
    cloudflare: logo('<path fill="#F38020" d="M18.6 15.2H7.8a3.6 3.6 0 0 1-.4-7.2 4.8 4.8 0 0 1 9.2-1.1 3.1 3.1 0 0 1 2 8.3z"/><path fill="#FAAD3F" d="M8.2 16.8h9.2c.8 0 1.2.2 1.2.8H7.4c0-.6.4-.8.8-.8z"/>'),
    custom: logo('<rect x="3" y="6" width="18" height="12" rx="2" fill="none" stroke="#334155" stroke-width="1.6"/><path stroke="#334155" stroke-width="1.6" d="M7 12h4M8 10v4M15 10v4M16 12h.1"/>')
};
