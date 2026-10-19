package ai

// Catalog of LLM providers. OpenAI-compatible vendors share one chat client.
// Anthropic, Gemini, and Ollama keep their own protocols.

func IsOpenAICompatible(p ProviderKind) bool {
	switch p {
	case ProviderOpenAI, ProviderMistral, ProviderGroq, ProviderOpenRouter, ProviderAzure,
		ProviderDeepSeek, ProviderXAI, ProviderCohere, ProviderTogether, ProviderFireworks,
		ProviderPerplexity, ProviderCerebras, ProviderSambaNova, ProviderNVIDIA,
		ProviderHuggingFace, ProviderGitHub, ProviderQwen, ProviderMoonshot,
		ProviderZhipu, ProviderMiniMax, ProviderLMStudio, ProviderCloudflare, ProviderCustom:
		return true
	default:
		return false
	}
}

// KeyOptional reports providers that answer without an API key (local runtimes).
func KeyOptional(p ProviderKind) bool {
	return p == ProviderOllama || p == ProviderLMStudio
}

// DefaultBaseURL is the origin the chat client appends "/chat/completions" to,
// unless the value already contains that path (Perplexity).
func DefaultBaseURL(p ProviderKind) string {
	switch p {
	case ProviderOpenAI:
		return "https://api.openai.com/v1"
	case ProviderMistral:
		return "https://api.mistral.ai/v1"
	case ProviderGroq:
		return "https://api.groq.com/openai/v1"
	case ProviderOpenRouter:
		return "https://openrouter.ai/api/v1"
	case ProviderAzure:
		return "https://YOUR-RESOURCE.openai.azure.com/openai/deployments/YOUR-DEPLOYMENT"
	case ProviderDeepSeek:
		return "https://api.deepseek.com/v1"
	case ProviderXAI:
		return "https://api.x.ai/v1"
	case ProviderCohere:
		return "https://api.cohere.ai/compatibility/v1"
	case ProviderTogether:
		return "https://api.together.xyz/v1"
	case ProviderFireworks:
		return "https://api.fireworks.ai/inference/v1"
	case ProviderPerplexity:
		return "https://api.perplexity.ai/chat/completions"
	case ProviderCerebras:
		return "https://api.cerebras.ai/v1"
	case ProviderSambaNova:
		return "https://api.sambanova.ai/v1"
	case ProviderNVIDIA:
		return "https://integrate.api.nvidia.com/v1"
	case ProviderHuggingFace:
		return "https://router.huggingface.co/v1"
	case ProviderGitHub:
		return "https://models.inference.ai.azure.com"
	case ProviderQwen:
		return "https://dashscope.aliyuncs.com/compatible-mode/v1"
	case ProviderMoonshot:
		return "https://api.moonshot.ai/v1"
	case ProviderZhipu:
		return "https://open.bigmodel.cn/api/paas/v4"
	case ProviderMiniMax:
		return "https://api.minimax.io/v1"
	case ProviderLMStudio:
		return "http://localhost:1234/v1"
	case ProviderCloudflare:
		return "https://api.cloudflare.com/client/v4/accounts/ACCOUNT_ID/ai/v1"
	case ProviderOllama:
		return "http://localhost:11434"
	case ProviderAnthropic:
		return "https://api.anthropic.com"
	case ProviderGemini:
		return "https://generativelanguage.googleapis.com/v1beta"
	default:
		return ""
	}
}

// DefaultKeyEnv is the server environment variable used when the UI key is empty.
func DefaultKeyEnv(p ProviderKind) string {
	switch p {
	case ProviderOpenAI, ProviderAzure, ProviderCustom:
		return "OPENAI_API_KEY"
	case ProviderAnthropic:
		return "ANTHROPIC_API_KEY"
	case ProviderGemini:
		return "GEMINI_API_KEY"
	case ProviderMistral:
		return "MISTRAL_API_KEY"
	case ProviderGroq:
		return "GROQ_API_KEY"
	case ProviderOpenRouter:
		return "OPENROUTER_API_KEY"
	case ProviderDeepSeek:
		return "DEEPSEEK_API_KEY"
	case ProviderXAI:
		return "XAI_API_KEY"
	case ProviderCohere:
		return "COHERE_API_KEY"
	case ProviderTogether:
		return "TOGETHER_API_KEY"
	case ProviderFireworks:
		return "FIREWORKS_API_KEY"
	case ProviderPerplexity:
		return "PERPLEXITY_API_KEY"
	case ProviderCerebras:
		return "CEREBRAS_API_KEY"
	case ProviderSambaNova:
		return "SAMBANOVA_API_KEY"
	case ProviderNVIDIA:
		return "NVIDIA_API_KEY"
	case ProviderHuggingFace:
		return "HF_TOKEN"
	case ProviderGitHub:
		return "GITHUB_TOKEN"
	case ProviderQwen:
		return "DASHSCOPE_API_KEY"
	case ProviderMoonshot:
		return "MOONSHOT_API_KEY"
	case ProviderZhipu:
		return "ZHIPU_API_KEY"
	case ProviderMiniMax:
		return "MINIMAX_API_KEY"
	case ProviderCloudflare:
		return "CLOUDFLARE_API_TOKEN"
	default:
		return ""
	}
}

// CatalogDefaultModel is the model id suggested for a provider.
func CatalogDefaultModel(p ProviderKind) string {
	switch p {
	case ProviderOllama:
		return "llama3.2"
	case ProviderOpenAI, ProviderAzure, ProviderCustom:
		return "gpt-4o-mini"
	case ProviderAnthropic:
		return "claude-3-5-haiku-20241022"
	case ProviderGemini:
		return "gemini-2.0-flash"
	case ProviderMistral:
		return "mistral-small-latest"
	case ProviderGroq:
		return "llama-3.1-8b-instant"
	case ProviderOpenRouter:
		return "openai/gpt-4o-mini"
	case ProviderDeepSeek:
		return "deepseek-chat"
	case ProviderXAI:
		return "grok-3"
	case ProviderCohere:
		return "command-r-plus"
	case ProviderTogether:
		return "meta-llama/Meta-Llama-3.1-8B-Instruct-Turbo"
	case ProviderFireworks:
		return "accounts/fireworks/models/llama-v3p1-8b-instruct"
	case ProviderPerplexity:
		return "sonar"
	case ProviderCerebras:
		return "llama3.1-8b"
	case ProviderSambaNova:
		return "Meta-Llama-3.1-8B-Instruct"
	case ProviderNVIDIA:
		return "meta/llama-3.1-8b-instruct"
	case ProviderHuggingFace:
		return "meta-llama/Llama-3.1-8B-Instruct"
	case ProviderGitHub:
		return "gpt-4o-mini"
	case ProviderQwen:
		return "qwen-plus"
	case ProviderMoonshot:
		return "moonshot-v1-8k"
	case ProviderZhipu:
		return "glm-4-flash"
	case ProviderMiniMax:
		return "MiniMax-Text-01"
	case ProviderLMStudio:
		return "local-model"
	case ProviderCloudflare:
		return "@cf/meta/llama-3.1-8b-instruct"
	default:
		return ""
	}
}
