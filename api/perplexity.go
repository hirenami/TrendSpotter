package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type PerplexityResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Item struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rank     int32    `json:"rank"`
	EndTimestamp int32 `json:"end_timestamp"`
	IncreasePercentage int32 `json:"increase_percentage"`
}

func (a *Api) CallPerplexityAPI(query []TrendingSearch) ([]Item, error) {
	// Perplexity APIのURL
	url := "https://api.perplexity.ai/chat/completions" // Perplexity APIエンドポイント

	var items []Item // 結果を格納するスライス
	rank := 0

	// .envファイルを読み込む
	env := godotenv.Load()
	if env != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")

	// HTTPクライアントを作成
	client := &http.Client{}

	// queryを順番に処理
	for _, q := range query {
		// プロンプトを作成
		prompt := fmt.Sprintf("次のアイテムが飲食店名、食品名、食材名、その他のカテゴリにどれに該当するかを一つだけ分類してください。飲食店名、食品名、食材名、その他のいずれかを一つのみ出力するようして、絶対にほかは出力させないでください。##出力例(厳守):その他 ##アイテム %s", q.Query)

		// リクエストボディの設定
		requestBody := map[string]interface{}{
			"model": "sonar", // モデルの指定。適切なものに変更
			"messages": []map[string]string{
				{
					"role":    "system",
					"content": "You are a helpful assistant.",
				},
				{
					"role":    "user",
					"content": prompt,
				},
			},
		}

		// リクエストボディをJSONに変換
		body, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}

		// HTTPリクエストを作成
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

		// ヘッダーを設定
		req.Header.Set("accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		// HTTPリクエストを送信
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// レスポンスボディを読み取る
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// JSONレスポンスを解析
		var response PerplexityResponse
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		// 出力内容をチェック
		content := response.Choices[0].Message.Content

		// "その他"の場合、スキップ
		if content != "飲食店名" && content != "食品名" && content != "食材名" {
			continue
		}

		// Rankを付ける
		rank += 1

		var item Item
		item.Rank = int32(rank)

		// "飲食店名"の場合
		if content == "飲食店名" {
			item.Name = q.Query
			item.Location = q.Query
			item.EndTimestamp = q.EndTimestamp
			item.IncreasePercentage = q.IncreasePercentage
		} else if content == "食品名" || content == "食材名" {
			// "食品名"または"食材名"の場合、場所を生成するためのプロンプトを作成
			locationPrompt := fmt.Sprintf("次のアイテムはどこで売られていますか？最も売られている可能性が高い場所の単語を一つ出力してください。地名ではなく、店でお願いします。その単語以外は絶対に出力しないでください。: %s", q.Query)

			// 売られている場所を問い合わせ
			locationRequestBody := map[string]interface{}{
				"model": "sonar",
				"messages": []map[string]string{
					{
						"role":    "system",
						"content": "You are a helpful assistant.",
					},
					{
						"role":    "user",
						"content": locationPrompt,
					},
				},
			}

			locationBody, err := json.Marshal(locationRequestBody)
			if err != nil {
				return nil, err
			}

			// 新たなリクエストを作成
			req, err = http.NewRequest("POST", url, bytes.NewBuffer(locationBody))
			if err != nil {
				return nil, err
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+apiKey)

			resp, err = client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			// レスポンスボディを読み取る
			respBody, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			var locationResponse PerplexityResponse
			err = json.Unmarshal(respBody, &locationResponse)
			if err != nil {
				return nil, err
			}

			// 売られている場所の情報を取得
			item.Name = q.Query
			item.Location = locationResponse.Choices[0].Message.Content
			item.EndTimestamp = q.EndTimestamp
			item.IncreasePercentage = q.IncreasePercentage
		}

		// itemsスライスに追加
		items = append(items, item)
	}

	return items, nil
}