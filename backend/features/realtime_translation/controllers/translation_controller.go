package controllers

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"

	"go-realtime-translation-with-speech-service/backend/features/realtime_translation/models"

	"github.com/gin-gonic/gin"
)

// TranslationService はリアルタイム翻訳サービスのインターフェース
type TranslationService interface {
	TranslateText(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error)
	StartStreamingSession(ctx context.Context, req *models.StreamingTranslationRequest) (string, error)
	ProcessAudioChunk(ctx context.Context, sessionID string, audioChunk []byte) ([]models.StreamingTranslationResponse, error)
	CloseStreamingSession(ctx context.Context, sessionID string) error
}

// TranslationController はリアルタイム翻訳に関するAPIエンドポイントを処理するコントローラー
type TranslationController struct {
	translationService TranslationService
}

// NewTranslationController は新しいTranslationControllerのインスタンスを作成する
func NewTranslationController(translationService TranslationService) *TranslationController {
	return &TranslationController{
		translationService: translationService,
	}
}

// RegisterRoutes はルーターにエンドポイントを登録する
func (c *TranslationController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/translate", c.TranslateText)

	// ストリーミング翻訳関連のエンドポイント
	streaming := router.Group("/streaming")
	{
		streaming.POST("/start", c.StartStreamingSession)
		streaming.POST("/process", c.ProcessAudioChunk)
		streaming.POST("/close", c.CloseStreamingSession)
	}
}

// TranslateText はテキスト翻訳を処理するエンドポイント
// POST /api/v1/translate
func (c *TranslationController) TranslateText(ctx *gin.Context) {
	var req models.TranslationRequest

	// リクエストボディをパース
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// リクエストのバリデーション
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 翻訳サービスを呼び出す
	resp, err := c.translationService.TranslateText(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Translation failed: " + err.Error()})
		return
	}

	// 成功レスポンスを返す
	ctx.JSON(http.StatusOK, resp)
}

// StartStreamingSession はストリーミング翻訳セッションを開始するエンドポイント
// POST /api/v1/streaming/start
func (c *TranslationController) StartStreamingSession(ctx *gin.Context) {
	var req models.StreamingTranslationRequest

	// リクエストボディをパース
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// リクエストのバリデーション
	if err := req.Validate(); err != nil {
		log.Printf("Invalid streaming request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Starting streaming session with source=%s, target=%s, format=%s",
		req.SourceLanguage, req.TargetLanguage, req.AudioFormat)

	// セッションを開始
	sessionID, err := c.translationService.StartStreamingSession(ctx, &req)
	if err != nil {
		log.Printf("Failed to start streaming session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start streaming session: " + err.Error()})
		return
	}

	log.Printf("Successfully started streaming session: %s", sessionID)

	// 成功レスポンスを返す
	ctx.JSON(http.StatusOK, gin.H{"sessionId": sessionID})
}

// AudioChunkRequest は音声チャンク処理リクエストの構造体
type AudioChunkRequest struct {
	SessionID  string `json:"sessionId"`
	AudioChunk string `json:"audioChunk"` // Base64エンコードされた音声データ
}

// ProcessAudioChunk は音声チャンクを処理するエンドポイント
// POST /api/v1/streaming/process
func (c *TranslationController) ProcessAudioChunk(ctx *gin.Context) {
	var req AudioChunkRequest

	// リクエストボディをパース
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// セッションIDのバリデーション
	if req.SessionID == "" {
		log.Print("Session ID is missing in request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// 音声データのバリデーション
	if req.AudioChunk == "" {
		log.Printf("Audio chunk is missing for session: %s", req.SessionID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audio chunk is required"})
		return
	}

	log.Printf("Processing audio chunk request for session: %s, audio data length: %d", req.SessionID, len(req.AudioChunk))

	// Base64デコード
	audioData, err := base64.StdEncoding.DecodeString(req.AudioChunk)
	if err != nil {
		log.Printf("Failed to decode base64 audio data for session %s: %v", req.SessionID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio data encoding: " + err.Error()})
		return
	}

	log.Printf("Decoded audio data for session %s: %d bytes", req.SessionID, len(audioData))

	// 音声チャンクを処理
	responses, err := c.translationService.ProcessAudioChunk(ctx, req.SessionID, audioData)
	if err != nil {
		log.Printf("Failed to process audio chunk for session %s: %v", req.SessionID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process audio chunk: " + err.Error()})
		return
	}

	log.Printf("Successfully processed audio chunk for session %s, got %d responses", req.SessionID, len(responses))

	// 成功レスポンスを返す
	ctx.JSON(http.StatusOK, responses)
}

// SessionRequest はセッションIDのみを含むリクエストの構造体
type SessionRequest struct {
	SessionID string `json:"sessionId"`
}

// CloseStreamingSession はストリーミングセッションを終了するエンドポイント
// POST /api/v1/streaming/close
func (c *TranslationController) CloseStreamingSession(ctx *gin.Context) {
	var req SessionRequest

	// リクエストボディをパース
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// セッションIDのバリデーション
	if req.SessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// セッションを終了
	if err := c.translationService.CloseStreamingSession(ctx, req.SessionID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close streaming session: " + err.Error()})
		return
	}

	// 成功レスポンスを返す
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
