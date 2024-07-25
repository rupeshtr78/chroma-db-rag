package gitquery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
)

const githubCodeAPI = "https://api.github.com/search/code"

type GitHubCodeResponse struct {
	TotalCount        int  `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
	Items             []struct {
		Name       string `json:"name"`
		Path       string `json:"path"`
		Sha        string `json:"sha"`
		URL        string `json:"url"`
		HTMLURL    string `json:"html_url"`
		Repository struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			FullName        string `json:"full_name"`
			Description     string `json:"description"`
			StargazersCount int    `json:"stargazers_count"`
			HTMLURL         string `json:"html_url"`
		} `json:"repository"`
	} `json:"items"`
}

func GitCodeQuery() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <search-query> <language>", os.Args[0])
	}
	query, language := os.Args[1], os.Args[2]

	// Get GitHub Token from Environment Variable (you should set it beforehand)
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	// GitHub API for searching code
	// To get a substantial number of results for sorting
	searchURL := fmt.Sprintf("%s?q=%s+language:%s&per_page=100", githubCodeAPI, url.QueryEscape(query), url.QueryEscape(language))

	// Create a new request with the Authorization header
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "token "+githubToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("GitHub API request failed with status: %s", resp.Status)
	}

	var result GitHubCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}

	// Sort the items by stargazers count in descending order
	sort.Slice(result.Items, func(i, j int) bool {
		return result.Items[i].Repository.StargazersCount > result.Items[j].Repository.StargazersCount
	})

	fmt.Printf("Total Count: %d\n", result.TotalCount)
	topCount := 10
	if len(result.Items) < topCount {
		topCount = len(result.Items)
	}

	// Print only the top 10 results
	for i := 0; i < topCount; i++ {
		item := result.Items[i]
		fmt.Printf("File Name: %s\nPath: %s\nURL: %s\nRepository Full Name: %s\nRepository URL: %s\nStars: %d\nDescription: %s\n\n",
			item.Name, item.Path, item.HTMLURL, item.Repository.FullName, item.Repository.HTMLURL, item.Repository.StargazersCount, item.Repository.Description)
	}
}
