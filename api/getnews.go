package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"source"`
		Author      string `json:"author"`
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PublishedAt string `json:"publishedAt"`
	} `json:"articles"`
}

func (a *Api) getNews(query string) ([]string, error) {
	// Google News APIのURL
	apiKey := "YOUR_NEWS_API_KEY"
	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&apiKey=%s", query, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var newsResponse NewsAPIResponse
	if err := json.Unmarshal(body, &newsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse news response: %v", err)
	}

	var newsTitles []string
	for _, article := range newsResponse.Articles {
		newsTitles = append(newsTitles, article.Title)
	}
	return newsTitles, nil
}
