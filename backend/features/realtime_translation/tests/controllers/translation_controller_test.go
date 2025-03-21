package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kohei3110/go-realtime-translation-with-speech-service/backend/features/realtime_translation/controllers"
	"github.com/kohei3110/go-realtime-translation-with-speech-service/backend/features/realtime_translation/models"
)

// MockTranslationService はTranslationServiceのモック
type MockTranslationService struct {
	mock.Mock
}

// TranslateText はテキスト翻訳のモックメソッド
func (m *MockTranslationService) TranslateText(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TranslationResponse), args.Error(1)
}

// StartStreamingSession はストリーミングセッション開始のモックメソッド
func (m *MockTranslationService) StartStreamingSession(ctx context.Context, req *models.StreamingTranslationRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

// ProcessAudioChunk は音声チャンク処理のモックメソッド
func (m *MockTranslationService) ProcessAudioChunk(ctx context.Context, sessionID string, audioChunk []byte) ([]models.StreamingTranslationResponse, error) {
	args := m.Called(ctx, sessionID, audioChunk)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StreamingTranslationResponse), args.Error(1)
}

// CloseStreamingSession はストリーミングセッション終了のモックメソッド
func (m *MockTranslationService) CloseStreamingSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

// テスト用のGinルーターを作成するヘルパー関数
func setupRouter(translationService *MockTranslationService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// コントローラーを作成して登録
	controller := controllers.NewTranslationController(translationService)

	// APIルートを設定
	apiGroup := router.Group("/api/v1")
	controller.RegisterRoutes(apiGroup)

	return router
}

func TestTranslateText(t *testing.T) {
	// テストケース1: 正常系 - 翻訳成功
	t.Run("SuccessfulTranslation", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// リクエストを作成
		req := models.TranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			Text:           "こんにちは、元気ですか？",
		}

		// 期待されるレスポンスを作成
		expectedResp := &models.TranslationResponse{
			SourceLanguage:  "ja",
			TargetLanguage:  "en",
			OriginalText:    "こんにちは、元気ですか？",
			TranslatedText:  "Hello, how are you?",
			ConfidenceScore: 0.95,
		}

		// モックの振る舞いを設定
		mockService.On("TranslateText", mock.Anything, &req).Return(expectedResp, nil)

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/translate", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, w.Code)

		var resp models.TranslationResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, expectedResp.SourceLanguage, resp.SourceLanguage)
		assert.Equal(t, expectedResp.TargetLanguage, resp.TargetLanguage)
		assert.Equal(t, expectedResp.OriginalText, resp.OriginalText)
		assert.Equal(t, expectedResp.TranslatedText, resp.TranslatedText)
		assert.Equal(t, expectedResp.ConfidenceScore, resp.ConfidenceScore)

		// モックが期待通り呼び出されたか確認
		mockService.AssertExpectations(t)
	})

	// テストケース2: 異常系 - バリデーションエラー
	t.Run("ValidationError", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// 不正なリクエストを作成（テキストが空）
		req := models.TranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			Text:           "", // 空のテキスト
		}

		// モックの振る舞いを設定（サービスは呼び出されないので設定不要）

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/translate", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証（400 Bad Requestが期待される）
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// エラーメッセージにテキスト関連の記述があることを確認
		var errResp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Contains(t, errResp["error"].(string), "text")

		// モックが呼び出されていないことを確認
		mockService.AssertNotCalled(t, "TranslateText")
	})

	// テストケース3: 異常系 - サービスエラー
	t.Run("ServiceError", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// リクエストを作成
		req := models.TranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			Text:           "こんにちは、元気ですか？",
		}

		// モックの振る舞いを設定（エラーを返す）
		mockService.On("TranslateText", mock.Anything, &req).Return(nil, errors.New("translation service error"))

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/translate", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証（500 Internal Server Errorが期待される）
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// エラーメッセージにサービスエラーの記述があることを確認
		var errResp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Contains(t, errResp["error"].(string), "translation")

		// モックが期待通り呼び出されたか確認
		mockService.AssertExpectations(t)
	})
}

func TestStartStreamingSession(t *testing.T) {
	// テストケース1: 正常系 - セッション開始成功
	t.Run("SuccessfulSessionStart", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// リクエストを作成
		req := models.StreamingTranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			AudioFormat:    "wav",
		}

		// モックの振る舞いを設定
		mockService.On("StartStreamingSession", mock.Anything, &req).Return("session-123", nil)

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/streaming/start", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "session-123", resp["sessionId"])

		// モックが期待通り呼び出されたか確認
		mockService.AssertExpectations(t)
	})

	// テストケース2: 異常系 - バリデーションエラー
	t.Run("ValidationError", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// 不正なリクエストを作成（未対応の音声フォーマット）
		req := models.StreamingTranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			AudioFormat:    "unknown",
		}

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/streaming/start", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証（400 Bad Requestが期待される）
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// モックが呼び出されていないことを確認
		mockService.AssertNotCalled(t, "StartStreamingSession")
	})

	// テストケース3: 異常系 - サービスエラー
	t.Run("ServiceError", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// リクエストを作成
		req := models.StreamingTranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			AudioFormat:    "wav",
		}

		// モックの振る舞いを設定（エラーを返す）
		mockService.On("StartStreamingSession", mock.Anything, &req).Return("", errors.New("session start error"))

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/streaming/start", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証（500 Internal Server Errorが期待される）
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// モックが期待通り呼び出されたか確認
		mockService.AssertExpectations(t)
	})
}

func TestProcessAudioChunk(t *testing.T) {
	// テストケース1: 正常系 - 音声チャンク処理成功
	t.Run("SuccessfulAudioProcessing", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// セッションIDとテスト用音声データ
		sessionID := "session-123"
		audioChunk := []byte{0x01, 0x02, 0x03, 0x04}

		// 期待されるレスポンスを作成
		expectedResponses := []models.StreamingTranslationResponse{
			{
				SourceLanguage: "ja",
				TargetLanguage: "en",
				TranslatedText: "Hello",
				IsFinal:        false,
				SegmentID:      "segment-1",
			},
		}

		// モックの振る舞いを設定
		mockService.On("ProcessAudioChunk", mock.Anything, sessionID, audioChunk).Return(expectedResponses, nil)

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// リクエストデータを作成
		requestData := map[string]interface{}{
			"sessionId":  sessionID,
			"audioChunk": audioChunk,
		}

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(requestData)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/streaming/process", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, w.Code)

		var resp []models.StreamingTranslationResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, 1, len(resp))
		assert.Equal(t, expectedResponses[0].SourceLanguage, resp[0].SourceLanguage)
		assert.Equal(t, expectedResponses[0].TargetLanguage, resp[0].TargetLanguage)
		assert.Equal(t, expectedResponses[0].TranslatedText, resp[0].TranslatedText)
		assert.Equal(t, expectedResponses[0].IsFinal, resp[0].IsFinal)
		assert.Equal(t, expectedResponses[0].SegmentID, resp[0].SegmentID)

		// モックが期待通り呼び出されたか確認
		mockService.AssertExpectations(t)
	})

	// テストケース2: 異常系 - セッションIDなし
	t.Run("MissingSessionID", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// セッションIDなし
		audioChunk := []byte{0x01, 0x02, 0x03, 0x04}

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// リクエストデータを作成（セッションIDなし）
		requestData := map[string]interface{}{
			"audioChunk": audioChunk,
		}

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(requestData)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/streaming/process", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証（400 Bad Requestが期待される）
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// モックが呼び出されていないことを確認
		mockService.AssertNotCalled(t, "ProcessAudioChunk")
	})

	// テストケース3: 異常系 - サービスエラー
	t.Run("ServiceError", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// セッションIDとテスト用音声データ
		sessionID := "session-123"
		audioChunk := []byte{0x01, 0x02, 0x03, 0x04}

		// モックの振る舞いを設定（エラーを返す）
		mockService.On("ProcessAudioChunk", mock.Anything, sessionID, audioChunk).Return(nil, errors.New("audio processing error"))

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// リクエストデータを作成
		requestData := map[string]interface{}{
			"sessionId":  sessionID,
			"audioChunk": audioChunk,
		}

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(requestData)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/streaming/process", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証（500 Internal Server Errorが期待される）
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// モックが期待通り呼び出されたか確認
		mockService.AssertExpectations(t)
	})
}

func TestCloseStreamingSession(t *testing.T) {
	// テストケース1: 正常系 - セッション終了成功
	t.Run("SuccessfulSessionClose", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// セッションID
		sessionID := "session-123"

		// モックの振る舞いを設定
		mockService.On("CloseStreamingSession", mock.Anything, sessionID).Return(nil)

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// リクエストデータを作成
		requestData := map[string]interface{}{
			"sessionId": sessionID,
		}

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(requestData)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/streaming/close", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, w.Code)

		// モックが期待通り呼び出されたか確認
		mockService.AssertExpectations(t)
	})

	// テストケース2: 異常系 - セッションIDなし
	t.Run("MissingSessionID", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// リクエストデータを作成（セッションIDなし）
		requestData := map[string]interface{}{}

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(requestData)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/streaming/close", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証（400 Bad Requestが期待される）
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// モックが呼び出されていないことを確認
		mockService.AssertNotCalled(t, "CloseStreamingSession")
	})

	// テストケース3: 異常系 - サービスエラー
	t.Run("ServiceError", func(t *testing.T) {
		// モックサービスを作成
		mockService := new(MockTranslationService)

		// セッションID
		sessionID := "session-123"

		// モックの振る舞いを設定（エラーを返す）
		mockService.On("CloseStreamingSession", mock.Anything, sessionID).Return(errors.New("session close error"))

		// ルーターをセットアップ
		router := setupRouter(mockService)

		// リクエストデータを作成
		requestData := map[string]interface{}{
			"sessionId": sessionID,
		}

		// HTTPリクエストを作成
		jsonBody, _ := json.Marshal(requestData)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/api/v1/streaming/close", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// リクエストを実行
		router.ServeHTTP(w, httpReq)

		// レスポンスを検証（500 Internal Server Errorが期待される）
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// モックが期待通り呼び出されたか確認
		mockService.AssertExpectations(t)
	})
}
