package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	translatortext "go-realtime-translation-with-speech-service/backend/translatortext"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// translatorClient はアプリケーション全体で使用する翻訳クライアント
var translatorClient *translatortext.TranslatorClient

// セッション情報を保持する構造体
type StreamingSession struct {
	ID             string
	SourceLanguage string
	TargetLanguage string
	AudioFormat    string
}

// アクティブなセッションを保持するマップ
var activeSessions = make(map[string]*StreamingSession)

// SetTranslatorClient は翻訳クライアントをセットします
func SetTranslatorClient(client *translatortext.TranslatorClient) {
	translatorClient = client
}

// TranslationRequest は翻訳リクエストの構造体
type TranslationRequest struct {
	Text           string `json:"text" binding:"required"`
	TargetLanguage string `json:"targetLanguage" binding:"required"`
	SourceLanguage string `json:"sourceLanguage"`
}

// TranslationResponse は翻訳レスポンスの構造体
type TranslationResponse struct {
	OriginalText   string  `json:"originalText"`
	TranslatedText string  `json:"translatedText"`
	SourceLanguage string  `json:"sourceLanguage"`
	TargetLanguage string  `json:"targetLanguage"`
	Confidence     float64 `json:"confidence,omitempty"`
}

// StreamingTranslationRequest はストリーミング翻訳開始リクエストの構造体
type StreamingTranslationRequest struct {
	SourceLanguage string `json:"sourceLanguage" binding:"required"`
	TargetLanguage string `json:"targetLanguage" binding:"required"`
	AudioFormat    string `json:"audioFormat" binding:"required"`
}

// AudioChunkRequest は音声チャンクリクエストの構造体
type AudioChunkRequest struct {
	SessionID  string `json:"sessionId" binding:"required"`
	AudioChunk string `json:"audioChunk" binding:"required"` // Base64エンコードされた音声データ
}

// StreamingTranslationResponse はストリーミング翻訳レスポンスの構造体
type StreamingTranslationResponse struct {
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage"`
	TranslatedText string `json:"translatedText"`
	IsFinal        bool   `json:"isFinal"`
	SegmentID      string `json:"segmentId"`
}

// SessionCloseRequest はセッション終了リクエストの構造体
type SessionCloseRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}

// TranslateHandler はテキスト翻訳のハンドラー
func TranslateHandler(c *gin.Context) {
	var req TranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ターゲット言語の設定
	targetLanguages := []string{req.TargetLanguage}

	// 翻訳リクエストの作成
	textParam := []*translatortext.TranslateTextInput{
		{
			Text: &req.Text,
		},
	}

	// 翻訳の実行
	log.Printf("翻訳リクエスト: %s", req.Text)
	log.Printf("ターゲット言語: %s", req.TargetLanguage)
	result, err := translatorClient.Translate(context.Background(), targetLanguages, textParam, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("翻訳の実行に失敗しました: %v", err)})
		return
	}

	// レスポンスの作成
	if result.TranslateResultAllItemArray != nil && len(result.TranslateResultAllItemArray) > 0 {
		item := result.TranslateResultAllItemArray[0]

		response := TranslationResponse{
			OriginalText:   req.Text,
			TargetLanguage: req.TargetLanguage,
		}

		// 検出された言語情報
		if item.DetectedLanguage != nil {
			response.SourceLanguage = *item.DetectedLanguage.Language
			response.Confidence = *item.DetectedLanguage.Score
		} else if req.SourceLanguage != "" {
			response.SourceLanguage = req.SourceLanguage
		}

		// 翻訳テキスト
		if item.Translations != nil && len(item.Translations) > 0 {
			response.TranslatedText = *item.Translations[0].Text
		}

		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "翻訳結果がありません"})
	}
}

// HealthCheckHandler はヘルスチェックのハンドラー
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// StartStreamingSessionHandler はストリーミング翻訳セッションを開始するハンドラー
func StartStreamingSessionHandler(c *gin.Context) {
	var req StreamingTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 新しいセッションIDを生成
	sessionID := uuid.New().String()

	// セッション情報を保存
	activeSessions[sessionID] = &StreamingSession{
		ID:             sessionID,
		SourceLanguage: req.SourceLanguage,
		TargetLanguage: req.TargetLanguage,
		AudioFormat:    req.AudioFormat,
	}

	c.JSON(http.StatusOK, gin.H{"sessionId": sessionID})
}

// ProcessAudioChunkHandler は音声チャンクを処理するハンドラー
func ProcessAudioChunkHandler(c *gin.Context) {
	var req AudioChunkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// セッションの存在確認
	session, exists := activeSessions[req.SessionID]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なセッションIDです"})
		return
	}

	// 実際のプロジェクトでは、ここで音声データをAzure Speech Serviceに送信して
	// 実際の翻訳を行う処理を実装します。
	// 1. 認証情報の取得
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("認証情報の取得に失敗しました: %v", err)
	}

	// 2. Speech Serviceのエンドポイント設定
	endpoint := "https://api.cognitive.microsofttranslator.com/"

	// 3. NewTranslatorClientの作成
	client, err := translatortext.NewTranslatorClient(endpoint, cred, nil)
	if err != nil {
		log.Fatalf("NewTranslatorClientの作成に失敗しました: %v", err)
	}

	// 4. 翻訳リクエストの作成
	// FIXME: 文字起こし後のテキストを使用する必要があります
	// ここではダミーのテキストを使用しています
	text := "Hello, how are you?"
	textParam := []*translatortext.TranslateTextInput{
		{
			Text: &text,
		},
	}

	// 5. 翻訳の実行
	ctx := context.Background()
	result, err := client.Translate(ctx, []string{"ja"}, textParam, nil)

	if result.TranslateResultAllItemArray != nil {
		for _, item := range result.TranslateResultAllItemArray {
			if item.Translations != nil {
				for j, translation := range item.Translations {
					log.Printf("Translation %d: %s\n", j, *translation.Text)
					log.Printf("Detected Language: %s\n", *item.DetectedLanguage.Language)
					log.Printf("Confidence: %f\n", *item.DetectedLanguage.Score)
					response := []StreamingTranslationResponse{
						{
							SourceLanguage: session.SourceLanguage,
							TargetLanguage: session.TargetLanguage,
							TranslatedText: *translation.Text,
							IsFinal:        true,
							SegmentID:      uuid.New().String(),
						},
					}
					c.JSON(http.StatusOK, response)
				}
			}
		}
	}
}

// CloseStreamingSessionHandler はストリーミングセッションを終了するハンドラー
func CloseStreamingSessionHandler(c *gin.Context) {
	var req SessionCloseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// セッションの削除
	delete(activeSessions, req.SessionID)

	c.JSON(http.StatusOK, gin.H{"status": "セッションを終了しました"})
}
