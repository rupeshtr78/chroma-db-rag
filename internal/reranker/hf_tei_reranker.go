package reranker

import (
	"bytes"
	"chroma-db/internal/constants"
	"chroma-db/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HttpRerankClient struct {
	Client         *http.Client
	BaseURL        string
	Model          string
	apiKey         string
	DefaultHeaders map[string]string
}

// Added constructor function
func NewHttpRerankClient(client *http.Client, baseURL, model, apiKey string, defaultHeaders map[string]string) *HttpRerankClient {
	return &HttpRerankClient{
		Client:         client,
		BaseURL:        baseURL,
		Model:          model,
		apiKey:         apiKey,
		DefaultHeaders: defaultHeaders,
	}
}

func GetReRankClient() *HttpRerankClient {
	return reRankClient
}

var (
	reRankClient *HttpRerankClient
)

func GetHttpRerankClient(client *http.Client, baseURL, model, apiKey string, defaultHeaders map[string]string) *HttpRerankClient {
	once.Do(func() {
		reRankClient = &HttpRerankClient{
			Client:         client,
			BaseURL:        baseURL,
			Model:          model,
			apiKey:         apiKey,
			DefaultHeaders: defaultHeaders,
		}
	})

	return reRankClient
}

type HttpRerankRequest struct {
	Query       string   `json:"query"`
	Texts       []string `json:"texts"`
	RawScores   bool     `json:"raw_scores"`
	ReturnTexts bool     `json:"return_text"`
}

func (c *HttpRerankRequest) JSON() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// [{"index":1,"score":0.9987814},{"index":0,"score":0.022949383}]%
type HttpRerankResponse struct {
	Index int     `json:"index"`
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

// CreateRerankingRequest creates a reranking request to the Hugging Face Reranker
// Returns a list of reranked responses or an error
// Example:
// q := "What is Deep Learning?"
// texts := []string{"Tomatos are fruits...", "Deep Learning is not...", "Deep learning is..."}
// Response: [{"index":2,"score":0.9987814},{"index":1,"score":0.022949383},{"index":0,"score":0.000076250595}]
func (c *HttpRerankClient) CreateRerankingRequest(ctx context.Context, req *HttpRerankRequest) (*[]HttpRerankResponse, error) {
	reqJSON, err := req.JSON()
	if err != nil {
		return nil, err
	}
	var url = c.BaseURL

	// + c.Model
	// if !strings.HasSuffix(c.BaseURL, "/") && c.Model != "" {
	// 	url = c.BaseURL + "/" + c.Model
	// }
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(reqJSON))
	if err != nil {
		return nil, err
	}
	for k, v := range c.DefaultHeaders {
		httpReq.Header.Set(k, v)
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected code [%v] while making a request to %v", resp.Status, url)
	}

	// read the response body
	respBuff := new(bytes.Buffer)
	_, err = io.Copy(respBuff, resp.Body)
	if err != nil {
		return nil, err
	}

	resData := respBuff.Bytes()
	logger.Log.Debug().Msgf("Response: %v\n", string(resData))

	var rerankResponses []HttpRerankResponse
	if err := json.Unmarshal(resData, &rerankResponses); err != nil {
		return nil, err
	}

	// logger.Log.Debug().Msgf("Rerank Responses: %v\n", rerankResponses)

	return &rerankResponses, nil
}

// RerankQueryResult reranks the query results using the HuggingFace reranker
// TODO: Use GRPC https://github.com/huggingface/text-embeddings-inference?tab=readme-ov-file#grpc
func (c *HttpRerankClient) RerankQueryResult(ctx context.Context, queryTexts []string, queryResults []string) (string, error) {

	queryString := strings.Builder{}
	for _, text := range queryTexts {
		queryString.WriteString(text)
	}
	request := &HttpRerankRequest{
		Query:       queryString.String(),
		Texts:       queryResults,
		RawScores:   false,
		ReturnTexts: true,
	}

	client := &HttpRerankClient{
		Client:  &http.Client{},
		BaseURL: constants.HuggingFaceRerankUrl,
		Model:   constants.HuggingFaceRerankModel,
	}

	res, err := client.CreateRerankingRequest(ctx, request)
	if err != nil {
		log.Error().Msgf("Error reranking query results: %v\n", err)
	}

	log.Info().Msgf("Reranked Results: %v\n", res)
	// For now return the first result
	firstResult := (*res)[0]
	reRankedResult := strings.TrimSpace(firstResult.Text)
	return reRankedResult, nil
}
