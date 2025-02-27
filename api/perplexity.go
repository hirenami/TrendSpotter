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

func (a *Api) CallPerplexityAPI(query string) (*PerplexityResponse, error) {
	// Perplexity APIのURL
	url := "https://api.perplexity.ai/chat/completions" // Perplexity APIエンドポイント

	// プロンプトを作成
	// 食品・食材か飲食店かなどを判定するプロンプト
	prompt := fmt.Sprintf("次のアイテムが飲食店名、食品名、食材名、その他のカテゴリにどれに該当するかを一つだけ分類してください。必ずカテゴリのみ出力するようして、ほかは一切出力させないでください: %s", query)

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

	// .envファイルを読み込む
	env := godotenv.Load()
	if env != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")

	// ヘッダーを設定
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// HTTPクライアントでリクエストを送信
	client := &http.Client{}
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

	// 食品名や食材名なら売られている場所を尋ねる
	if response.Choices[0].Message.Content == "食品名" || response.Choices[0].Message.Content == "食材名" {
		// 売られている場所を尋ねる
		locationPrompt := fmt.Sprintf("次のアイテムはどこで売られていますか？最も売られている可能性が高い場所の単語を一つ出力してください。その単語以外は絶対に出力しないでください。: %s", query)
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

		// 売られている場所を問い合わせ
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

		// 売られている場所の情報を回答としてセット
		response.Choices[0].Message.Content = locationResponse.Choices[0].Message.Content
	}

	return &response, nil
}
