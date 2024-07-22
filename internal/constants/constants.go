package constants

var (
	ChromaUrl      string  = "http://0.0.0.0:8070"
	TenantName     string  = "ollama_tenant-01"
	Database       string  = "ollama_database-01"
	Namespace      string  = "chroma-ollama-01"
	ScoreThreshold float32 = 0.65
	OllamaModel    string  = "nomic-embed-text" //nomic-embed-text" //"mxbai-embed-large"
	OllamaUrl      string  = "http://10.0.0.213:11434"
	DistanceFn     string  = "cosine"
)
