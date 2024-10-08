package constants

import (
	"github.com/amikos-tech/chroma-go/types"
)

var (
	ChromaUrl              string                 = "http://0.0.0.0:8070"
	TenantName             string                 = "ollama_tenant-01"
	Database               string                 = "ollama_database-01"
	Namespace              string                 = "chroma-ollama-01"
	Collection             string                 = "ollama_collection-01"
	ScoreThreshold         float32                = 0.65
	OllamaUrl              string                 = "http://10.0.0.213:11434"
	OllamaEmbdedModel      string                 = "nomic-embed-text" //nomic-embed-text" //"mxbai-embed-large"
	OllamaChatModel        string                 = "llama3.1:8b"
	OpenAIModel            string                 = "gpt-3.5-turbo"
	OpenAIApiKey           string                 = "OPENAI_API_KEY"
	OPenAIEmbedModel       string                 = "gpt-3.5-turbo"
	DistanceFn             types.DistanceFunction = types.L2
	TemplateFile           string                 = "internal/prompts/prompt_template.tmpl"
	SystemPromptFile       string                 = "internal/prompts/system_prompt_explain.tmpl" // "internal/prompts/system_prompt_explain.tmpl"
	LogLevel               string                 = "debug"
	HuggingFaceTeiUrl      string                 = "http://10.0.0.213:50080/embed"
	HuggingFaceEmbedModel  string                 = "BAAI/bge-large-en-v1.5"
	HuggingFaceRerankUrl   string                 = "http://10.0.0.213:50081/rerank"
	HuggingFaceRerankModel string                 = "BAAI/bge-reranker-large"
	GrpcEmbeddingEnabled   bool                   = false
	GrpcEmbedServer        string                 = "10.0.0.213:50083"
	GrpcRerankEnabled      bool                   = false
	GrpcRerankServer       string                 = "10.0.0.213:50083"
)

type Metadata map[string]interface{}

type DocType string

// Supported Document types
const (
	PDF DocType = "pdf"
	TXT DocType = "txt"
)

type Embedder int

// Supported Embedders
const (
	HuggingFace Embedder = iota
	Ollama
	OpenAI
)

type Protocol int

// Supported Protocols
const (
	HTTP Protocol = iota
	HTTPS
	GRPC
)

type LLMProvider int

// Supported Chat Models
const (
	OpenAIChat LLMProvider = iota
	OllamaChat
	HuggingFaceChat
)
