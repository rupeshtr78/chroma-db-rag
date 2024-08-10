package hfmodels

import (
	"context"
	"fmt"

	huggingface "github.com/amikos-tech/chroma-go/hf"
)

func HfEmbeddingModel() {
	ef, err := huggingface.NewHuggingFaceEmbeddingInferenceFunction("http://localhost:8001/embed") //set this to the URL of the HuggingFace Embedding Inference Server
	if err != nil {
		fmt.Printf("Error creating HuggingFace embedding function: %s \n", err)
	}
	documents := []string{
		"Document 1 content here",
	}
	resp, reqErr := ef.EmbedDocuments(context.Background(), documents)
	if reqErr != nil {
		fmt.Printf("Error embedding documents: %s \n", reqErr)
	}
	fmt.Printf("Embedding response: %v \n", resp)
}
