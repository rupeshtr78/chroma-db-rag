package reranker

import (
	"bytes"
	"chroma-db/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RerankManeger interface {
	CreateRerankingRequest(ctx context.Context, req *HfRerankRequest) ([]HfRerankResponse, error)
}

type HfRerankClient struct {
	Client         *http.Client
	BaseURL        string
	Model          string
	apiKey         string
	DefaultHeaders map[string]string
}

type HfRerankRequest struct {
	Query     string   `json:"query"`
	Texts     []string `json:"texts"`
	RawScores bool     `json:"raw_scores"`
}

func (c *HfRerankRequest) JSON() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// [{"index":1,"score":0.9987814},{"index":0,"score":0.022949383}]%
type HfRerankResponse struct {
	Index int     `json:"index"`
	Score float64 `json:"score"`
}

// CreateRerankingRequest creates a reranking request to the Hugging Face Reranker
// Returns a list of reranked responses or an error
// Example:
// q := "What is Deep Learning?"
// texts := []string{"Tomatos are fruits...", "Deep Learning is not...", "Deep learning is..."}
// Response: [{"index":2,"score":0.9987814},{"index":1,"score":0.022949383},{"index":0,"score":0.000076250595}]
func (c *HfRerankClient) CreateRerankingRequest(ctx context.Context, req *HfRerankRequest) ([]HfRerankResponse, error) {
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

	var rerankResponses []HfRerankResponse
	if err := json.Unmarshal(resData, &rerankResponses); err != nil {
		return nil, err
	}

	return rerankResponses, nil
}
