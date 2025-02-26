package api

import (
	g "github.com/serpapi/google-search-results-golang"
)

func (a *Api) getTrend() {
	parameter := map[string]string{
		"engine": "google_trends_trending_now",
		"geo":    "JP",
	}

	search := g.NewGoogleSearch(parameter, "a315f225081bc2a6a47570925e9fc45ce22f3fa1fe083c826bdc41219613a893")
	results, err := search.GetJSON()
	if err != nil {
		panic(err)
	}
	trending_searches := results["trending_searches"]

	for i := 0; i < len(trending_searches.([]interface{})); i++ {
		println(trending_searches.([]interface{})[i].(string))
	}
}
