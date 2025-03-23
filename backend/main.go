package main

import (
	"log"
	"os"

	"go-realtime-translation-with-speech-service/backend/api/handlers"
	translatortext "go-realtime-translation-with-speech-service/backend/translatortext"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 認証情報の取得
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("認証情報の取得に失敗しました: %v", err)
	}

	// 2. Speech Serviceのエンドポイント設定
	endpoint := "https://api.cognitive.microsofttranslator.com/"

	// 3. TranslatorClientの作成
	client, err := translatortext.NewTranslatorClient(endpoint, cred, nil)
	if err != nil {
		log.Fatalf("TranslatorClientの作成に失敗しました: %v", err)
	}

	// ハンドラーに翻訳クライアントをセット
	handlers.SetTranslatorClient(client)

	// Ginルーターの設定
	router := gin.Default()

	// CORSミドルウェアの設定
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// APIグループの設定
	api := router.Group("/api/v1")
	{
		// ヘルスチェックエンドポイント
		api.GET("/health", handlers.HealthCheckHandler)

		// 翻訳エンドポイント
		api.POST("/translate", handlers.TranslateHandler)

		// ストリーミング翻訳関連エンドポイント
		streaming := api.Group("/streaming")
		{
			streaming.POST("/start", handlers.StartStreamingSessionHandler)
			streaming.POST("/process", handlers.ProcessAudioChunkHandler)
			streaming.POST("/close", handlers.CloseStreamingSessionHandler)
		}
	}

	// ポート番号の取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// サーバーの起動
	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
