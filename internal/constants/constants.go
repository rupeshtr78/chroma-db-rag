package constants

import "github.com/amikos-tech/chroma-go/types"

var (
	ChromaUrl         string                 = "http://0.0.0.0:8070"
	TenantName        string                 = "ollama_tenant-01"
	Database          string                 = "ollama_database-01"
	Namespace         string                 = "chroma-ollama-01"
	Collection        string                 = "ollama_collection-01"
	ScoreThreshold    float32                = 0.65
	OllamaEmbdedModel string                 = "nomic-embed-text" //nomic-embed-text" //"mxbai-embed-large"
	OllamaChatModel   string                 = "llama3:latest"
	OllamaUrl         string                 = "http://10.0.0.213:11434"
	DistanceFn        types.DistanceFunction = types.COSINE
	TemplateFile      string                 = "internal/prompts/prompt_template.tmpl"
	LogLevel          string                 = "debug"
)

type Metadata map[string]interface{}

type DocType string

// Supported Document types
const (
	PDF DocType = "pdf"
	TXT DocType = "txt"
)
