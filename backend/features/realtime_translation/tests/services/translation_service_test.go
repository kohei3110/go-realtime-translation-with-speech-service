package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"go-realtime-translation-with-speech-service/backend/features/realtime_translation/models"
	"go-realtime-translation-with-speech-service/backend/features/realtime_translation/services"
)

// Mock for speech service
type MockSpeechService struct {
	mock.Mock
}

func (m *MockSpeechService) TranslateText(ctx context.Context, sourceLanguage, targetLanguage, text string) (string, float64, error) {
	args := m.Called(ctx, sourceLanguage, targetLanguage, text)
	return args.String(0), args.Get(1).(float64), args.Error(2)
}

func (m *MockSpeechService) StartStreamingSession(ctx context.Context, sourceLanguage, targetLanguage, audioFormat string) (string, error) {
	args := m.Called(ctx, sourceLanguage, targetLanguage, audioFormat)
	return args.String(0), args.Error(1)
}

func (m *MockSpeechService) ProcessAudioChunk(ctx context.Context, sessionID string, audioChunk []byte) ([]models.StreamingTranslationResponse, error) {
	args := m.Called(ctx, sessionID, audioChunk)
	return args.Get(0).([]models.StreamingTranslationResponse), args.Error(1)
}

func (m *MockSpeechService) CloseStreamingSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func TestTranslateText(t *testing.T) {
	// テストケース1: 正常系 - 翻訳が成功する場合
	t.Run("SuccessfulTranslation", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// モックの振る舞いを設定
		mockSpeech.On("TranslateText",
			mock.Anything,  // Context
			"ja",           // SourceLanguage
			"en",           // TargetLanguage
			"こんにちは、元気ですか？", // Text
		).Return("Hello, how are you?", 0.95, nil)

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト対象のメソッドを実行
		req := models.TranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			Text:           "こんにちは、元気ですか？",
		}

		// リクエストを実行
		resp, err := translationService.TranslateText(context.Background(), &req)

		// 検証
		require.NoError(t, err)
		assert.Equal(t, "ja", resp.SourceLanguage)
		assert.Equal(t, "en", resp.TargetLanguage)
		assert.Equal(t, "こんにちは、元気ですか？", resp.OriginalText)
		assert.Equal(t, "Hello, how are you?", resp.TranslatedText)
		assert.Equal(t, 0.95, resp.ConfidenceScore)

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})

	// テストケース2: 異常系 - 翻訳サービスがエラーを返す場合
	t.Run("TranslationServiceError", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// モックの振る舞いを設定（エラーを返す）
		mockSpeech.On("TranslateText",
			mock.Anything,  // Context
			"ja",           // SourceLanguage
			"en",           // TargetLanguage
			"こんにちは、元気ですか？", // Text
		).Return("", 0.0, errors.New("translation service error"))

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト対象のメソッドを実行
		req := models.TranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			Text:           "こんにちは、元気ですか？",
		}

		// リクエストを実行
		resp, err := translationService.TranslateText(context.Background(), &req)

		// 検証
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "translation service error")

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})

	// テストケース3: 異常系 - リクエストのバリデーションエラー
	t.Run("ValidationError", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)
		// このテストではモックは呼び出されないので振る舞いの設定は不要

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト対象のメソッドを実行（不正なリクエスト）
		req := models.TranslationRequest{
			SourceLanguage: "", // 空の値でバリデーションエラーを発生させる
			TargetLanguage: "en",
			Text:           "こんにちは、元気ですか？",
		}

		// リクエストを実行
		resp, err := translationService.TranslateText(context.Background(), &req)

		// 検証
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "source language")

		// モックが呼び出されていないことを確認
		mockSpeech.AssertNotCalled(t, "TranslateText")
	})
}

func TestStartStreamingSession(t *testing.T) {
	// テストケース1: 正常系 - ストリーミングセッションの開始が成功する場合
	t.Run("SuccessfulStreamingStart", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// モックの振る舞いを設定
		mockSpeech.On("StartStreamingSession",
			mock.Anything, // Context
			"ja",          // SourceLanguage
			"en",          // TargetLanguage
			"wav",         // AudioFormat
		).Return("session-123", nil)

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト対象のメソッドを実行
		req := models.StreamingTranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			AudioFormat:    "wav",
		}

		// リクエストを実行
		sessionID, err := translationService.StartStreamingSession(context.Background(), &req)

		// 検証
		require.NoError(t, err)
		assert.Equal(t, "session-123", sessionID)

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})

	// テストケース2: 異常系 - セッション開始時にエラーが発生する場合
	t.Run("StreamingStartError", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// モックの振る舞いを設定（エラーを返す）
		mockSpeech.On("StartStreamingSession",
			mock.Anything, // Context
			"ja",          // SourceLanguage
			"en",          // TargetLanguage
			"wav",         // AudioFormat
		).Return("", errors.New("streaming session start error"))

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト対象のメソッドを実行
		req := models.StreamingTranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			AudioFormat:    "wav",
		}

		// リクエストを実行
		sessionID, err := translationService.StartStreamingSession(context.Background(), &req)

		// 検証
		require.Error(t, err)
		assert.Empty(t, sessionID)
		assert.Contains(t, err.Error(), "streaming session start error")

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})

	// テストケース3: 異常系 - リクエストのバリデーションエラー
	t.Run("ValidationError", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)
		// このテストではモックは呼び出されないので振る舞いの設定は不要

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト対象のメソッドを実行（不正なリクエスト）
		req := models.StreamingTranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			AudioFormat:    "unknown", // 未対応のフォーマットでバリデーションエラーを発生させる
		}

		// リクエストを実行
		sessionID, err := translationService.StartStreamingSession(context.Background(), &req)

		// 検証
		require.Error(t, err)
		assert.Empty(t, sessionID)
		assert.Contains(t, err.Error(), "audio format")

		// モックが呼び出されていないことを確認
		mockSpeech.AssertNotCalled(t, "StartStreamingSession")
	})
}

func TestProcessAudioChunk(t *testing.T) {
	// テストケース1: 正常系 - 音声チャンクの処理が成功する場合
	t.Run("SuccessfulAudioChunkProcessing", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// 期待される応答を設定
		expectedResponse := []models.StreamingTranslationResponse{
			{
				SourceLanguage: "ja",
				TargetLanguage: "en",
				TranslatedText: "Hello",
				IsFinal:        false,
				SegmentID:      "segment-1",
			},
		}

		// モックの振る舞いを設定
		mockSpeech.On("ProcessAudioChunk",
			mock.Anything,                  // Context
			"session-123",                  // SessionID
			mock.AnythingOfType("[]uint8"), // AudioChunk
		).Return(expectedResponse, nil)

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト用の音声データを準備
		audioChunk := []byte{0x01, 0x02, 0x03, 0x04}

		// テスト対象のメソッドを実行
		responses, err := translationService.ProcessAudioChunk(context.Background(), "session-123", audioChunk)

		// 検証
		require.NoError(t, err)
		assert.Equal(t, 1, len(responses))
		assert.Equal(t, "ja", responses[0].SourceLanguage)
		assert.Equal(t, "en", responses[0].TargetLanguage)
		assert.Equal(t, "Hello", responses[0].TranslatedText)
		assert.False(t, responses[0].IsFinal)
		assert.Equal(t, "segment-1", responses[0].SegmentID)

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})

	// テストケース2: 異常系 - 音声チャンク処理時にエラーが発生する場合
	t.Run("AudioChunkProcessingError", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// モックの振る舞いを設定（エラーを返す）
		mockSpeech.On("ProcessAudioChunk",
			mock.Anything,                  // Context
			"session-123",                  // SessionID
			mock.AnythingOfType("[]uint8"), // AudioChunk
		).Return([]models.StreamingTranslationResponse{}, errors.New("audio processing error"))

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト用の音声データを準備
		audioChunk := []byte{0x01, 0x02, 0x03, 0x04}

		// テスト対象のメソッドを実行
		responses, err := translationService.ProcessAudioChunk(context.Background(), "session-123", audioChunk)

		// 検証
		require.Error(t, err)
		assert.Empty(t, responses)
		assert.Contains(t, err.Error(), "audio processing error")

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})

	// テストケース3: 異常系 - 無効なセッションID
	t.Run("InvalidSessionID", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// モックの振る舞いを設定（エラーを返す）
		mockSpeech.On("ProcessAudioChunk",
			mock.Anything,                  // Context
			"",                             // 空のSessionID
			mock.AnythingOfType("[]uint8"), // AudioChunk
		).Return([]models.StreamingTranslationResponse{}, errors.New("invalid session ID"))

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト用の音声データを準備
		audioChunk := []byte{0x01, 0x02, 0x03, 0x04}

		// テスト対象のメソッドを実行（空のセッションID）
		responses, err := translationService.ProcessAudioChunk(context.Background(), "", audioChunk)

		// 検証
		require.Error(t, err)
		assert.Empty(t, responses)
		assert.Contains(t, err.Error(), "invalid session ID")

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})
}

func TestCloseStreamingSession(t *testing.T) {
	// テストケース1: 正常系 - ストリーミングセッションの終了が成功する場合
	t.Run("SuccessfulStreamingClose", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// モックの振る舞いを設定
		mockSpeech.On("CloseStreamingSession",
			mock.Anything, // Context
			"session-123", // SessionID
		).Return(nil)

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト対象のメソッドを実行
		err := translationService.CloseStreamingSession(context.Background(), "session-123")

		// 検証
		require.NoError(t, err)

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})

	// テストケース2: 異常系 - セッション終了時にエラーが発生する場合
	t.Run("StreamingCloseError", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// モックの振る舞いを設定（エラーを返す）
		mockSpeech.On("CloseStreamingSession",
			mock.Anything, // Context
			"session-123", // SessionID
		).Return(errors.New("streaming session close error"))

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト対象のメソッドを実行
		err := translationService.CloseStreamingSession(context.Background(), "session-123")

		// 検証
		require.Error(t, err)
		assert.Contains(t, err.Error(), "streaming session close error")

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})

	// テストケース3: 異常系 - 無効なセッションID
	t.Run("InvalidSessionID", func(t *testing.T) {
		// モックの準備
		mockSpeech := new(MockSpeechService)

		// モックの振る舞いを設定（エラーを返す）
		mockSpeech.On("CloseStreamingSession",
			mock.Anything, // Context
			"",            // 空のSessionID
		).Return(errors.New("invalid session ID"))

		// テスト対象のサービスを作成
		translationService := services.NewTranslationService(mockSpeech)

		// テスト対象のメソッドを実行（空のセッションID）
		err := translationService.CloseStreamingSession(context.Background(), "")

		// 検証
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid session ID")

		// モックが期待通り呼び出されたか確認
		mockSpeech.AssertExpectations(t)
	})
}
