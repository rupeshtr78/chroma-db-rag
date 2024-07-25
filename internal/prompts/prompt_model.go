package prompts

import (
	"bytes"
	"chroma-db/internal/constants"
	"fmt"
	"log"
	"text/template"
)

// PromptData holds the context and prompt to be injected into the template
type PromptData struct {
	System  string
	Content string
	Prompt  string
}

func GetTemplate(prompt string, content string) (string, error) {
	// Define the path to the template file
	// templateFile := "internal/templates/prompt_template.tmpl"

	// Load the template from the file
	tmpl, err := template.ParseFiles(constants.TemplateFile)
	if err != nil {
		log.Fatalf("Failed to parse template file: %v", err)
		return "", err
	}

	// Provide the context and prompt data
	data := PromptData{
		Content: content,
		Prompt:  prompt,
	}

	// Execute the template with the provided data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
		return "", err
	}

	// Get the final prompt string
	finalPrompt := buf.String()

	// Send the final prompt to Ollama for processing
	fmt.Println("Final Prompt:", finalPrompt)
	// Here you would send `finalPrompt` to Ollama.

	return finalPrompt, err
}
