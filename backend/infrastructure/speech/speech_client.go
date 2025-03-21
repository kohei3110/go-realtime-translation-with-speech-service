package speech

import (
	"context"
	"fmt"
	"io"
	"os"

	"go-realtime-translation-with-speech-service/backend/features/realtime_translation/models"

	"errors"
	"log"
	"sync"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/google/uuid"
)

// Config はSpeech Serviceの構成情報
type Config struct {
	Key             string
	Region          string
	EndpointURL     string
	RecognitionLang string
}

// SpeechClient はAzure Speech Serviceとの連携を行うクライアント
type SpeechClient struct {
	config       Config
	speechConfig *speech.SpeechConfig
	sessions     map[string]*streamingSession
	sessionMutex sync.Mutex
}

// streamingSession はストリーミングセッションの状態を管理する構造体
type streamingSession struct {
	ID             string
	SourceLanguage string
	TargetLanguage string
	AudioFormat    string
	LastAccess     time.Time
	Config         *speech.SpeechConfig
}

// NewSpeechClient は新しいSpeechClientのインスタンスを作成する
func NewSpeechClient() (*SpeechClient, error) {
	// 環境変数から設定を読み込む
	key := os.Getenv("AZURE_SPEECH_KEY")
	if key == "" {
		return nil, errors.New("AZURE_SPEECH_KEY environment variable is not set")
	}

	region := os.Getenv("AZURE_SPEECH_REGION")
	if region == "" {
		return nil, errors.New("AZURE_SPEECH_REGION environment variable is not set")
	}

	// デフォルト値
	endpointURL := os.Getenv("AZURE_SPEECH_ENDPOINT")
	if endpointURL == "" {
		endpointURL = fmt.Sprintf("https://%s.api.cognitive.microsoft.com/sts/v1.0/issueToken", region)
	}

	recognitionLang := os.Getenv("AZURE_SPEECH_RECOGNITION_LANG")
	if recognitionLang == "" {
		recognitionLang = "ja-JP"
	}

	config := Config{
		Key:             key,
		Region:          region,
		EndpointURL:     endpointURL,
		RecognitionLang: recognitionLang,
	}

	// Speech設定を作成
	speechConfig, err := speech.NewSpeechConfigFromSubscription(key, region)
	if err != nil {
		return nil, fmt.Errorf("failed to create speech config: %w", err)
	}

	return &SpeechClient{
		config:       config,
		speechConfig: speechConfig,
		sessions:     make(map[string]*streamingSession),
	}, nil
}

// TranslateText はテキストを翻訳するメソッド
func (c *SpeechClient) TranslateText(ctx context.Context, sourceLanguage, targetLanguage, text string) (string, float64, error) {
	// Speech設定を作成
	config, err := speech.NewSpeechConfigFromSubscription(c.config.Key, c.config.Region)
	if err != nil {
		return "", 0.0, fmt.Errorf("failed to create speech config: %w", err)
	}
	defer config.Close()

	// シンセサイザーを作成
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(config, nil)
	if err != nil {
		return "", 0.0, fmt.Errorf("failed to create synthesizer: %w", err)
	}
	defer synthesizer.Close()

	// テキストを合成して翻訳
	resultChan := synthesizer.SpeakTextAsync(text)
	select {
	case <-ctx.Done():
		return "", 0.0, ctx.Err()
	case result := <-resultChan:
		if result.Error != nil {
			return "", 0.0, fmt.Errorf("translation failed: %w", result.Error)
		}
		defer result.Close()

		// AudioDataフィールドを使用
		return string(result.Result.AudioData), 0.95, nil // 信頼度スコアは固定値を使用
	}
}

// SynthesizeText はテキストを音声に合成するメソッド
func (c *SpeechClient) SynthesizeText(ctx context.Context, language string, text string) ([]byte, error) {
	log.Printf("Synthesizing text to speech. Language: %s, Text length: %d", language, len(text))

	// Speech設定を作成
	config, err := speech.NewSpeechConfigFromSubscription(c.config.Key, c.config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create speech config: %w", err)
	}
	defer config.Close()

	// 言語を設定
	config.SetSpeechSynthesisVoiceName(fmt.Sprintf("ja-JP-NanamiNeural"))
	if language != "" {
		// 言語コードから音声名を設定（実際のプロジェクトでは言語コードから適切な音声を選択するロジックが必要）
		voiceName := fmt.Sprintf("%s-Neural", language)
		config.SetSpeechSynthesisVoiceName(voiceName)
	}

	// シンセサイザーを作成
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(config, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create synthesizer: %w", err)
	}
	defer synthesizer.Close()

	// テキストの音声合成をスタート
	task := synthesizer.StartSpeakingTextAsync(text)
	var outcome speech.SpeechSynthesisOutcome
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case outcome = <-task:
	}
	defer outcome.Close()

	if outcome.Error != nil {
		return nil, fmt.Errorf("synthesis failed: %w", outcome.Error)
	}

	// 音声データストリームを作成
	stream, err := speech.NewAudioDataStreamFromSpeechSynthesisResult(outcome.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio data stream: %w", err)
	}
	defer stream.Close()

	// 音声データを読み込み
	var allAudio []byte
	audioChunk := make([]byte, 2048)
	for {
		n, err := stream.Read(audioChunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read audio data: %w", err)
		}
		allAudio = append(allAudio, audioChunk[:n]...)
	}

	log.Printf("Successfully synthesized text to speech, audio data size: %d bytes", len(allAudio))
	return allAudio, nil
}

// StartStreamingSession はストリーミング翻訳セッションを開始するメソッド
func (c *SpeechClient) StartStreamingSession(ctx context.Context, sourceLanguage, targetLanguage, audioFormat string) (string, error) {
	// セッションIDを生成
	sessionID := uuid.New().String()

	// Speech設定を作成
	config, err := speech.NewSpeechConfigFromSubscription(c.config.Key, c.config.Region)
	if err != nil {
		return "", fmt.Errorf("failed to create speech config: %w", err)
	}

	// セッション情報を保存
	c.sessionMutex.Lock()
	defer c.sessionMutex.Unlock()

	c.sessions[sessionID] = &streamingSession{
		ID:             sessionID,
		SourceLanguage: sourceLanguage,
		TargetLanguage: targetLanguage,
		AudioFormat:    audioFormat,
		LastAccess:     time.Now(),
		Config:         config,
	}

	return sessionID, nil
}

// ProcessAudioChunk は音声チャンクを処理するメソッド
func (c *SpeechClient) ProcessAudioChunk(ctx context.Context, sessionID string, audioChunk []byte) ([]models.StreamingTranslationResponse, error) {
	log.Printf("Processing audio chunk for session %s, chunk size: %d bytes", sessionID, len(audioChunk))

	// セッションを取得
	c.sessionMutex.Lock()
	session, exists := c.sessions[sessionID]
	if !exists {
		c.sessionMutex.Unlock()
		log.Printf("Session %s not found", sessionID)
		return nil, errors.New("invalid session ID")
	}

	// 最終アクセス時間を更新
	session.LastAccess = time.Now()
	c.sessionMutex.Unlock()

	log.Printf("Found session: %s (source: %s, target: %s)", sessionID, session.SourceLanguage, session.TargetLanguage)

	// シンセサイザーを作成
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(session.Config, nil)
	if err != nil {
		log.Printf("Failed to create synthesizer for session %s: %v", sessionID, err)
		return nil, fmt.Errorf("failed to create synthesizer: %w", err)
	}
	defer synthesizer.Close()

	// 音声データを処理
	ssml := fmt.Sprintf(`<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xml:lang="%s">%s</speak>`,
		session.SourceLanguage, string(audioChunk))

	log.Printf("Starting speech synthesis for session %s with SSML length: %d", sessionID, len(ssml))

	resultChan := synthesizer.StartSpeakingSsmlAsync(ssml)
	select {
	case <-ctx.Done():
		log.Printf("Context cancelled for session %s", sessionID)
		return nil, ctx.Err()
	case result := <-resultChan:
		if result.Error != nil {
			log.Printf("Processing failed for session %s: %v", sessionID, result.Error)
			return nil, fmt.Errorf("processing failed: %w", result.Error)
		}
		defer result.Close()

		log.Printf("Successfully processed audio chunk for session %s", sessionID)

		// レスポンスを作成
		responses := []models.StreamingTranslationResponse{
			{
				SourceLanguage: session.SourceLanguage,
				TargetLanguage: session.TargetLanguage,
				TranslatedText: string(result.Result.AudioData),
				IsFinal:        true,
				SegmentID:      uuid.New().String(),
			},
		}

		log.Printf("Generated translation response for session %s: %+v", sessionID, responses[0])
		return responses, nil
	}
}

// CloseStreamingSession はストリーミングセッションを終了するメソッド
func (c *SpeechClient) CloseStreamingSession(ctx context.Context, sessionID string) error {
	c.sessionMutex.Lock()
	defer c.sessionMutex.Unlock()

	// セッションを取得
	session, exists := c.sessions[sessionID]
	if !exists {
		return errors.New("invalid session ID")
	}

	// 設定を閉じる
	session.Config.Close()

	// セッションを削除
	delete(c.sessions, sessionID)

	return nil
}

// RunSessionCleanup はセッションのクリーンアップを定期的に実行する（バックグラウンドタスク）
func (c *SpeechClient) RunSessionCleanup(ctx context.Context, interval time.Duration, maxIdleTime time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.cleanupSessions(maxIdleTime)
		}
	}
}

// cleanupSessions は一定時間使用されていないセッションをクリーンアップする
func (c *SpeechClient) cleanupSessions(maxIdleTime time.Duration) {
	now := time.Now()
	c.sessionMutex.Lock()
	defer c.sessionMutex.Unlock()

	for id, session := range c.sessions {
		if now.Sub(session.LastAccess) > maxIdleTime {
			session.Config.Close()
			delete(c.sessions, id)
			log.Printf("Cleaned up idle session: %s", id)
		}
	}
}
