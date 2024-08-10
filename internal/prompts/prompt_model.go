package prompts

import (
	"bytes"
	"chroma-db/internal/constants"
	"chroma-db/pkg/logger"
	"os"
	"text/template"
)

var log = logger.Log

// PromptData holds the context and prompt to be injected into the template
type PromptData struct {
	SystemPrompt string
	Content      string
	Prompt       string
}

// GetTemplate generates a prompt using the provided system prompt, prompt, and content
func GetTemplate(sytemPromptFile string, prompt string, content string) (string, error) {

	// Load the template from the file
	tmpl, err := template.ParseFiles(constants.TemplateFile)
	if err != nil {
		log.Fatal().Msgf("Failed to parse template file: %v", err)
		return "", err
	}

	systemPrompt, err := os.ReadFile(sytemPromptFile)
	if err != nil {
		log.Fatal().Msgf("Failed to read system prompt file: %v", err)
		return "", err
	}

	// Provide the context and prompt data
	data := PromptData{
		SystemPrompt: string(systemPrompt),
		Content:      content,
		Prompt:       prompt,
	}

	// Execute the template with the provided data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal().Msgf("Failed to execute template: %v", err)
		return "", err
	}

	// Get the final prompt string
	finalPrompt := buf.String()

	// Send the final prompt to Ollama for processing
	return finalPrompt, err
}
