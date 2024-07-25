package gitquery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

const githubAPI = "https://api.github.com/search/repositories"

type GitHubResponse struct {
	TotalCount        int  `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
	Items             []struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		FullName        string `json:"full_name"`
		HTMLURL         string `json:"html_url"`
		Description     string `json:"description"`
		StargazersCount int    `json:"stargazers_count"`
	} `json:"items"`
}

func GitQuery() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <search-query>", os.Args[0])
	}
	query := os.Args[1]

	searchURL := fmt.Sprintf("%s?q=%s", githubAPI, url.QueryEscape(query))
	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("GitHub API request failed with status: %s", resp.Status)
	}

	var result GitHubResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}

	fmt.Printf("Total Count: %d\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("Name: %s\nFull Name: %s\nURL: %s\nStars: %d\nDescription: %s\n\n",
			item.Name, item.FullName, item.HTMLURL, item.StargazersCount, item.Description)
	}
}
