package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hirenami/TrendSpotter/api"
	"github.com/hirenami/TrendSpotter/dao"
	"github.com/hirenami/TrendSpotter/handler"
	"github.com/hirenami/TrendSpotter/sqlc"
	"github.com/hirenami/TrendSpotter/usecase"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("db/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	if mysqlUser == "" || mysqlUserPwd == "" || mysqlDatabase == "" {
		log.Fatal("fail :Getenv")
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3308)/%s?parseTime=true", mysqlUser, mysqlUserPwd, mysqlDatabase))

	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("fail: db.Ping, %v\n", err)
	}
	defer db.Close()

	_db := sqlc.New(db)
	Dao := dao.NewDao(db, _db)
	Usecase := usecase.NewUsecase(Dao)
	Handler := handler.Newhandler(Usecase)

	r := handler.SetupRoutes(Handler)

	a := api.NewApi()
	//a.GetTrend()

	items := []string{
		"しゃぶ葉 牛タン", "ジョナサン", "築地 銀だこ", "コナズ珈琲", "青いじゃがりこ",
		"ヒロシ", "ファミリーマート", "ドミノ・ピザ", "キットカット", "バーガーキング",
		"霞ヶ浦", "行田市", "アサイー", "河合郁人", "ガストフィットメニュー",
		"ロイヤルホスト", "dish", "z世代", "不二家", "しさんうどん",
		"シャカシャカポテト", "マックスバリュ",
	}

	for _, item := range items {
		// Perplexity APIで分類を確認
		perplexityResponse, err := a.CallPerplexityAPI(item)
		if err != nil {
			log.Printf("Perplexity APIエラー: %v\n", err)
			continue
		}

		// 分類結果を出力
		log.Printf(perplexityResponse.Choices[0].Message.Content)
	}

	log.Println("Listening...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
