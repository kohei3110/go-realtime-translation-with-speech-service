package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kohei3110/go-realtime-translation-with-speech-service/backend/features/realtime_translation/controllers"
	"github.com/kohei3110/go-realtime-translation-with-speech-service/backend/features/realtime_translation/services"
	"github.com/kohei3110/go-realtime-translation-with-speech-service/backend/infrastructure/speech"
)

// SetupRouter はAPIルーターをセットアップする
func SetupRouter() (*gin.Engine, error) {
	// Ginルーターを作成
	router := gin.Default()

	// CORSの設定
	router.Use(corsMiddleware())

	// ヘルスチェックエンドポイント
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Azure Speech Service クライアントを初期化
	speechClient, err := speech.NewSpeechClient()
	if err != nil {
		return nil, err
	}

	// リアルタイム翻訳サービスを作成
	translationService := services.NewTranslationService(speechClient)

	// リアルタイム翻訳コントローラーを作成
	translationController := controllers.NewTranslationController(translationService)

	// APIルートグループを作成
	api := router.Group("/api/v1")

	// リアルタイム翻訳のエンドポイントを登録
	translationController.RegisterRoutes(api)

	return router, nil
}

// corsMiddleware はCORS設定を行うミドルウェア
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
