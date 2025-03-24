package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sync"

	"go-realtime-translation-with-speech-service/backend/gospeech"
	translatortext "go-realtime-translation-with-speech-service/backend/translatortext"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// translatorClient はアプリケーション全体で使用する翻訳クライアント
var translatorClient *translatortext.TranslatorClient

// speechSubscriptionKey はAzure Speech Serviceのサブスクリプションキー
var speechSubscriptionKey string

// speechRegion はAzure Speech Serviceのリージョン
var speechRegion string

// SetSpeechCredentials はSpeech Serviceの認証情報をセットします
func SetSpeechCredentials(subscriptionKey, region string) {
	speechSubscriptionKey = subscriptionKey
	speechRegion = region
}

// セッション情報を保持する構造体
type StreamingSession struct {
	ID             string
	SourceLanguage string
	TargetLanguage string
	AudioFormat    string
	Recognizer     *gospeech.TranslationRecognizer
	WSConnection   *websocket.Conn
	Context        context.Context
	CancelFunc     context.CancelFunc
}

// WebSocketアップグレードの設定
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // すべてのオリジンを許可（本番環境では注意）
	},
}

// アクティブなセッションを保持するマップとそのロック
var (
	activeSessionsMutex sync.RWMutex
	activeSessions      = make(map[string]*StreamingSession)
)

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
	OriginalText   string `json:"originalText"`
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

	// WebSocketへのアップグレードを待機するエンドポイントのURLを返す
	c.JSON(http.StatusOK, gin.H{
		"sessionId":      sessionID,
		"webSocketURL":   fmt.Sprintf("/api/v1/streaming/ws/%s", sessionID),
		"sourceLanguage": req.SourceLanguage,
		"targetLanguage": req.TargetLanguage,
	})
}

// ProcessAudioChunkHandler は音声チャンクを処理するハンドラー
func ProcessAudioChunkHandler(c *gin.Context) {
	var req AudioChunkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// セッションの存在確認
	activeSessionsMutex.RLock()
	_, exists := activeSessions[req.SessionID]
	activeSessionsMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なセッションIDです"})
		return
	}

	// Base64エンコードされた音声データをデコード
	_, err := base64.StdEncoding.DecodeString(req.AudioChunk)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "音声データのデコードに失敗しました"})
		return
	}

	// このエンドポイントは主にRESTfulなアプローチの場合に使用されます
	// WebSocketを使用する場合は、WebSocketハンドラー内で音声処理を行います
	// ここでは、シンプルなレスポンスを返します
	c.JSON(http.StatusOK, gin.H{"status": "音声チャンクを受信しました"})
}

// WebSocketHandler はWebSocket接続を処理するハンドラー
func WebSocketHandler(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// WebSocketにアップグレード
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocketへのアップグレードに失敗しました: %v", err)
		return
	}

	// セッションコンテキストの作成
	ctx, cancel := context.WithCancel(context.Background())

	// Speech Translation設定
	translationConfig, err := gospeech.SpeechTranslationConfigFromSubscription(speechSubscriptionKey, speechRegion)
	if err != nil {
		log.Printf("Speech Translation設定の作成に失敗しました: %v", err)
		conn.Close()
		return
	}

	// オーディオ設定（デフォルトマイク）
	audioConfig, err := gospeech.NewAudioConfigFromDefaultMicrophone()
	if err != nil {
		log.Printf("オーディオ設定の作成に失敗しました: %v", err)
		conn.Close()
		return
	}

	// クライアントからの初期設定メッセージを待機
	var setupMsg StreamingTranslationRequest
	if err := conn.ReadJSON(&setupMsg); err != nil {
		log.Printf("初期設定メッセージの読み取りに失敗しました: %v", err)
		conn.Close()
		return
	}

	// 認識する言語の設定
	translationConfig.SetSpeechRecognitionLanguage(setupMsg.SourceLanguage)

	// 翻訳先言語の追加
	translationConfig.AddTargetLanguage(setupMsg.TargetLanguage)

	// 音声認識器の作成
	recognizer, err := gospeech.NewTranslationRecognizer(translationConfig, audioConfig)
	if err != nil {
		log.Printf("音声認識器の作成に失敗しました: %v", err)
		conn.Close()
		return
	}

	// セッション情報を保存
	session := &StreamingSession{
		ID:             sessionID,
		SourceLanguage: setupMsg.SourceLanguage,
		TargetLanguage: setupMsg.TargetLanguage,
		AudioFormat:    setupMsg.AudioFormat,
		Recognizer:     recognizer,
		WSConnection:   conn,
		Context:        ctx,
		CancelFunc:     cancel,
	}

	// セッションの保存
	activeSessionsMutex.Lock()
	activeSessions[sessionID] = session
	activeSessionsMutex.Unlock()

	// クライアントに準備完了を通知
	conn.WriteJSON(gin.H{"status": "ready", "sessionId": sessionID})

	// 認識結果のイベントハンドラーの設定
	recognizer.Recognized().Connect(func(eventArgs interface{}) {
		args, ok := eventArgs.(*gospeech.TranslationRecognitionEventArgs)
		if !ok {
			return
		}

		result := args.Result
		if result.Reason == gospeech.ResultReasonTranslatedSpeech {
			// 翻訳結果を取得
			translatedText, exists := result.Translations[setupMsg.TargetLanguage]
			if !exists {
				log.Printf("指定された言語の翻訳結果がありません: %s", setupMsg.TargetLanguage)
				return
			}

			// WebSocketを通じて結果を送信
			response := StreamingTranslationResponse{
				SourceLanguage: setupMsg.SourceLanguage,
				TargetLanguage: setupMsg.TargetLanguage,
				TranslatedText: translatedText,
				OriginalText:   result.Text,
				IsFinal:        true,
				SegmentID:      uuid.New().String(),
			}

			if err := conn.WriteJSON(response); err != nil {
				log.Printf("WebSocketへの書き込みに失敗しました: %v", err)
			}
		}
	})

	// 認識中イベントのハンドラー（途中経過）
	recognizer.Recognizing().Connect(func(eventArgs interface{}) {
		args, ok := eventArgs.(*gospeech.TranslationRecognitionEventArgs)
		if !ok {
			return
		}

		result := args.Result
		if result.Reason == gospeech.ResultReasonTranslatedSpeech {
			// 翻訳結果を取得
			translatedText, exists := result.Translations[setupMsg.TargetLanguage]
			if !exists {
				return
			}

			// WebSocketを通じて途中経過を送信
			response := StreamingTranslationResponse{
				SourceLanguage: setupMsg.SourceLanguage,
				TargetLanguage: setupMsg.TargetLanguage,
				TranslatedText: translatedText,
				OriginalText:   result.Text,
				IsFinal:        false,
				SegmentID:      uuid.New().String(),
			}

			if err := conn.WriteJSON(response); err != nil {
				log.Printf("WebSocketへの書き込みに失敗しました: %v", err)
			}
		}
	})

	// 連続認識を開始
	if err := recognizer.StartContinuousRecognition(ctx); err != nil {
		log.Printf("連続認識の開始に失敗しました: %v", err)
		conn.WriteJSON(gin.H{"error": "連続認識の開始に失敗しました"})
		conn.Close()
		return
	}

	// WebSocketのクローズを監視
	go func() {
		defer func() {
			// 連続認識を停止
			if err := recognizer.StopContinuousRecognition(); err != nil {
				log.Printf("連続認識の停止に失敗しました: %v", err)
			}

			// 認識器のクリーンアップ
			if err := recognizer.Close(); err != nil {
				log.Printf("認識器のクリーンアップに失敗しました: %v", err)
			}

			// セッションを削除
			activeSessionsMutex.Lock()
			delete(activeSessions, sessionID)
			activeSessionsMutex.Unlock()

			// キャンセル関数を呼び出し
			cancel()

			// WebSocket接続を閉じる
			conn.Close()

			log.Printf("セッション %s を終了しました", sessionID)
		}()

		for {
			// クライアントからのメッセージを待機
			_, _, err := conn.ReadMessage()
			if err != nil {
				// クライアントが切断した場合など
				log.Printf("WebSocket読み取りエラー: %v", err)
				return
			}
		}
	}()
}

// CloseStreamingSessionHandler はストリーミングセッションを終了するハンドラー
func CloseStreamingSessionHandler(c *gin.Context) {
	var req SessionCloseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// セッションの存在確認
	activeSessionsMutex.RLock()
	session, exists := activeSessions[req.SessionID]
	activeSessionsMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusOK, gin.H{"status": "セッションはすでに終了しています"})
		return
	}

	// WebSocket接続が存在する場合は閉じる
	if session.WSConnection != nil {
		session.WSConnection.Close()
	}

	// 音声認識器が存在する場合はクリーンアップ
	if session.Recognizer != nil {
		// 連続認識を停止
		session.Recognizer.StopContinuousRecognition()
		// 認識器のクリーンアップ
		session.Recognizer.Close()
	}

	// キャンセル関数を呼び出し
	if session.CancelFunc != nil {
		session.CancelFunc()
	}

	// セッションを削除
	activeSessionsMutex.Lock()
	delete(activeSessions, req.SessionID)
	activeSessionsMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "セッションを終了しました"})
}
