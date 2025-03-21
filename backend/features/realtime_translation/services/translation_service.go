package services

import (
	"context"
	"errors"
	"fmt"

	"go-realtime-translation-with-speech-service/backend/features/realtime_translation/models"
)

// SpeechService はAzure Speech Serviceとの連携を行うインターフェース
type SpeechService interface {
	// TranslateText はテキストを翻訳するメソッド
	TranslateText(ctx context.Context, sourceLanguage, targetLanguage, text string) (string, float64, error)

	// StartStreamingSession はストリーミング翻訳セッションを開始するメソッド
	StartStreamingSession(ctx context.Context, sourceLanguage, targetLanguage, audioFormat string) (string, error)

	// ProcessAudioChunk は音声チャンクを処理するメソッド
	ProcessAudioChunk(ctx context.Context, sessionID string, audioChunk []byte) ([]models.StreamingTranslationResponse, error)

	// CloseStreamingSession はストリーミングセッションを終了するメソッド
	CloseStreamingSession(ctx context.Context, sessionID string) error
}

// TranslationService はリアルタイム翻訳サービスの実装
type TranslationService struct {
	speechService SpeechService
}

// NewTranslationService は新しいTranslationServiceのインスタンスを作成する
func NewTranslationService(speechService SpeechService) *TranslationService {
	return &TranslationService{
		speechService: speechService,
	}
}

// TranslateText はテキストを翻訳するメソッド
func (s *TranslationService) TranslateText(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
	// リクエストのバリデーション
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Azure Speech Serviceを呼び出してテキスト翻訳を実行
	translatedText, confidenceScore, err := s.speechService.TranslateText(
		ctx,
		req.SourceLanguage,
		req.TargetLanguage,
		req.Text,
	)

	if err != nil {
		return nil, fmt.Errorf("translation failed: %w", err)
	}

	// レスポンスを作成
	response := &models.TranslationResponse{
		SourceLanguage:  req.SourceLanguage,
		TargetLanguage:  req.TargetLanguage,
		OriginalText:    req.Text,
		TranslatedText:  translatedText,
		ConfidenceScore: confidenceScore,
	}

	return response, nil
}

// StartStreamingSession はストリーミング翻訳セッションを開始するメソッド
func (s *TranslationService) StartStreamingSession(ctx context.Context, req *models.StreamingTranslationRequest) (string, error) {
	// リクエストのバリデーション
	if err := req.Validate(); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	// Azure Speech Serviceを呼び出してストリーミングセッションを開始
	sessionID, err := s.speechService.StartStreamingSession(
		ctx,
		req.SourceLanguage,
		req.TargetLanguage,
		req.AudioFormat,
	)

	if err != nil {
		return "", fmt.Errorf("failed to start streaming session: %w", err)
	}

	return sessionID, nil
}

// ProcessAudioChunk は音声チャンクを処理するメソッド
func (s *TranslationService) ProcessAudioChunk(ctx context.Context, sessionID string, audioChunk []byte) ([]models.StreamingTranslationResponse, error) {
	// セッションIDのバリデーション
	if sessionID == "" {
		return nil, errors.New("session ID is required")
	}

	// 音声チャンクのバリデーション
	if len(audioChunk) == 0 {
		return nil, errors.New("audio chunk is empty")
	}

	// Azure Speech Serviceを呼び出して音声チャンクを処理
	responses, err := s.speechService.ProcessAudioChunk(ctx, sessionID, audioChunk)
	if err != nil {
		return nil, fmt.Errorf("failed to process audio chunk: %w", err)
	}

	return responses, nil
}

// CloseStreamingSession はストリーミングセッションを終了するメソッド
func (s *TranslationService) CloseStreamingSession(ctx context.Context, sessionID string) error {
	// セッションIDのバリデーション
	if sessionID == "" {
		return errors.New("session ID is required")
	}

	// Azure Speech Serviceを呼び出してストリーミングセッションを終了
	if err := s.speechService.CloseStreamingSession(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to close streaming session: %w", err)
	}

	return nil
}
