package api

import (
	"fmt"

	g "github.com/serpapi/google-search-results-golang"
)

func (a *Api) GetTrend() {
	parameter := map[string]string{
		"engine": "google_trends_trending_now",
		"geo":    "JP",
	}

	search := g.NewGoogleSearch(parameter, "a315f225081bc2a6a47570925e9fc45ce22f3fa1fe083c826bdc41219613a893")
	results, err := search.GetJSON()
	if err != nil {
		panic(err)
	}

	trendingSearches, ok := results["trending_searches"].([]interface{})
	if !ok {
		panic("unexpected type for trending_searches")
	}

	// 各トレンド情報を走査
	for i := 0; i < len(trendingSearches); i++ {
		trend, ok := trendingSearches[i].(map[string]interface{})
		if !ok {
			continue
		}

		categories, ok := trend["categories"].([]interface{})
		if !ok {
			continue
		}

		// カテゴリIDが5かどうかをチェック
		for _, cat := range categories {
			category, ok := cat.(map[string]interface{})
			if !ok {
				continue
			}
			if idVal, exists := category["id"]; exists {
				if idFloat, ok := idVal.(float64); ok && int(idFloat) == 5 {
					// カテゴリID 5 (Food and Drink) に該当する場合、クエリを出力
					if query, ok := trend["query"].(string); ok {
						fmt.Println(query)
					}
					break // 該当するカテゴリが見つかったので、他のカテゴリはチェック不要
				}
			}
		}
	}
}
