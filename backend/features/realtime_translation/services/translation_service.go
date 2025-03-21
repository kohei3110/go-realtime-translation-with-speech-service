package services

import (
	"context"
	"errors"
	"fmt"
	"log"

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

	// SynthesizeText はテキストを音声に合成するメソッド
	SynthesizeText(ctx context.Context, language string, text string) ([]byte, error)
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
	log.Printf("Service: Starting streaming session with request: %+v", req)

	// リクエストのバリデーション
	if err := req.Validate(); err != nil {
		log.Printf("Service: Invalid streaming request: %v", err)
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
		log.Printf("Service: Failed to start streaming session: %v", err)
		return "", fmt.Errorf("failed to start streaming session: %w", err)
	}

	log.Printf("Service: Successfully started streaming session: %s", sessionID)
	return sessionID, nil
}

// ProcessAudioChunk は音声チャンクを処理するメソッド
func (s *TranslationService) ProcessAudioChunk(ctx context.Context, sessionID string, audioChunk []byte) ([]models.StreamingTranslationResponse, error) {
	log.Printf("Service: Processing audio chunk for session %s, size: %d bytes", sessionID, len(audioChunk))

	// セッションIDのバリデーション
	if sessionID == "" {
		log.Print("Service: Empty session ID provided")
		return nil, errors.New("session ID is required")
	}

	// 音声チャンクのバリデーション
	if len(audioChunk) == 0 {
		log.Printf("Service: Empty audio chunk for session %s", sessionID)
		return nil, errors.New("audio chunk is empty")
	}

	// Azure Speech Serviceを呼び出して音声チャンクを処理
	responses, err := s.speechService.ProcessAudioChunk(ctx, sessionID, audioChunk)
	if err != nil {
		log.Printf("Service: Failed to process audio chunk for session %s: %v", sessionID, err)
		return nil, fmt.Errorf("failed to process audio chunk: %w", err)
	}

	log.Printf("Service: Successfully processed audio chunk for session %s, got %d responses", sessionID, len(responses))
	return responses, nil
}

// CloseStreamingSession はストリーミングセッションを終了するメソッド
func (s *TranslationService) CloseStreamingSession(ctx context.Context, sessionID string) error {
	log.Printf("Service: Closing streaming session: %s", sessionID)

	// セッションIDのバリデーション
	if sessionID == "" {
		log.Print("Service: Empty session ID provided")
		return errors.New("session ID is required")
	}

	// Azure Speech Serviceを呼び出してストリーミングセッションを終了
	if err := s.speechService.CloseStreamingSession(ctx, sessionID); err != nil {
		log.Printf("Service: Failed to close streaming session %s: %v", sessionID, err)
		return fmt.Errorf("failed to close streaming session: %w", err)
	}

	log.Printf("Service: Successfully closed streaming session: %s", sessionID)
	return nil
}

// SynthesizeTextToSpeech はテキストを音声に合成するメソッド
func (s *TranslationService) SynthesizeTextToSpeech(ctx context.Context, req *models.SynthesisRequest) (*models.SynthesisResponse, error) {
	log.Printf("Service: Synthesizing text to speech with request: %+v", req)

	// リクエストのバリデーション
	if err := req.Validate(); err != nil {
		log.Printf("Service: Invalid synthesis request: %v", err)
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Azure Speech Serviceを呼び出してテキストを音声に合成
	audioData, err := s.speechService.SynthesizeText(ctx, req.Language, req.Text)
	if err != nil {
		log.Printf("Service: Failed to synthesize text: %v", err)
		return nil, fmt.Errorf("failed to synthesize text: %w", err)
	}

	log.Printf("Service: Successfully synthesized text to speech, audio data size: %d bytes", len(audioData))
	// レスポンスを作成
	response := &models.SynthesisResponse{
		Language:  req.Language,
		Text:      req.Text,
		AudioData: audioData,
	}

	return response, nil
}
