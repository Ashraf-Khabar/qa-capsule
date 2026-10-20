/**
 * AI provider metadata and brand logos for the configuration picker.
 * Marks are PNG files in /img/llm (LobeHub icon set, MIT), shown on a white chip.
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

const PNG_LOGOS = new Set([
    'openai', 'anthropic', 'gemini', 'mistral', 'groq', 'openrouter', 'azure', 'ollama',
    'deepseek', 'xai', 'cohere', 'together', 'fireworks', 'perplexity', 'cerebras',
    'sambanova', 'nvidia', 'huggingface', 'github', 'qwen', 'moonshot', 'zhipu',
    'minimax', 'lmstudio', 'cloudflare'
]);

export function logoForProvider(providerId) {
    if (PNG_LOGOS.has(providerId)) {
        return `<img class="ai-provider-logo" src="/img/llm/${providerId}.png" alt="" width="32" height="32" draggable="false">`;
    }
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
    custom: logo('<rect x="3" y="6" width="18" height="12" rx="2" fill="none" stroke="#334155" stroke-width="1.6"/><path stroke="#334155" stroke-width="1.6" d="M7 12h4M8 10v4M15 10v4"/>')
};
