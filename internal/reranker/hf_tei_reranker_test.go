package reranker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestHttpRerankRequest_JSON(t *testing.T) {
	req := HttpRerankRequest{
		Query:       "test query",
		Texts:       []string{"text1", "text2"},
		RawScores:   true,
		ReturnTexts: false,
	}

	jsonStr, err := req.JSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{"query": "test query", "texts": ["text1", "text2"], "raw_scores": true, "return_text": false}`, jsonStr)
}

func TestGetHttpRerankClient(t *testing.T) {
	client := &http.Client{}
	baseURL := "http://example.com"
	model := "test-model"
	apiKey := "test-key"
	headers := map[string]string{"Header": "Value"}

	rerankClient := GetHttpRerankClient(client, baseURL, model, apiKey, headers)
	assert.NotNil(t, rerankClient)
	assert.Equal(t, client, rerankClient.Client)
	assert.Equal(t, baseURL, rerankClient.BaseURL)
	assert.Equal(t, model, rerankClient.Model)
	assert.Equal(t, apiKey, rerankClient.apiKey)
	assert.Equal(t, headers, rerankClient.DefaultHeaders)
}

func TestHttpRerankClient_CreateRerankingRequest_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "error message", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &HttpRerankClient{
		Client:  server.Client(),
		BaseURL: server.URL,
	}
	req := &HttpRerankRequest{
		Query:       "test query",
		Texts:       []string{"text1", "text2"},
		RawScores:   true,
		ReturnTexts: false,
	}

	res, err := client.CreateRerankingRequest(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestHttpRerankClient_RerankQueryResult(t *testing.T) {
	respData := `[{"index":0,"text":"reranked result","score":0.998}]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(respData))
	}))
	defer server.Close()

	client := &HttpRerankClient{
		Client:  server.Client(),
		BaseURL: server.URL,
	}

	queryTexts := []string{"test query"}
	queryResults := []string{"result1", "result2"}

	_, err := client.RerankQueryResult(context.Background(), queryTexts, queryResults)
	assert.NoError(t, err)
}
