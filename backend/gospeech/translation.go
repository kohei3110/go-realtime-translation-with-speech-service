// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license.

package gospeech

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// SpeechTranslationConfig contains configuration settings for speech translation services
type SpeechTranslationConfig struct {
	*SpeechConfig
	targetLanguages []string
	voiceName       string
}

// NewSpeechTranslationConfig creates a new speech translation configuration
func NewSpeechTranslationConfig() *SpeechTranslationConfig {
	return &SpeechTranslationConfig{
		SpeechConfig:    NewSpeechConfig(),
		targetLanguages: []string{},
		voiceName:       "",
	}
}

// SpeechTranslationConfigFromSubscription creates a speech translation config from subscription information
func SpeechTranslationConfigFromSubscription(subscriptionKey, region string) (*SpeechTranslationConfig, error) {
	if subscriptionKey == "" {
		return nil, errors.New("subscription key cannot be empty")
	}
	if region == "" {
		return nil, errors.New("region cannot be empty")
	}

	config := NewSpeechTranslationConfig()
	config.SetProperty(SpeechServiceConnectionKey, subscriptionKey)
	config.SetProperty(SpeechServiceConnectionRegion, region)

	return config, nil
}

// SpeechTranslationConfigFromEndpoint creates a speech translation config from an endpoint
func SpeechTranslationConfigFromEndpoint(endpoint, subscriptionKey string) (*SpeechTranslationConfig, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint cannot be empty")
	}

	config := NewSpeechTranslationConfig()
	config.SetProperty(SpeechServiceConnectionEndpoint, endpoint)

	if subscriptionKey != "" {
		config.SetProperty(SpeechServiceConnectionKey, subscriptionKey)
	}

	return config, nil
}

// SpeechTranslationConfigFromHost creates a speech translation config from a host address
func SpeechTranslationConfigFromHost(host, subscriptionKey string) (*SpeechTranslationConfig, error) {
	if host == "" {
		return nil, errors.New("host cannot be empty")
	}

	config := NewSpeechTranslationConfig()
	config.SetProperty(SpeechServiceConnectionHost, host)

	if subscriptionKey != "" {
		config.SetProperty(SpeechServiceConnectionKey, subscriptionKey)
	}

	return config, nil
}

// SpeechTranslationConfigFromAuthToken creates a speech translation config from an authorization token
func SpeechTranslationConfigFromAuthToken(authToken, region string) (*SpeechTranslationConfig, error) {
	if authToken == "" {
		return nil, errors.New("authorization token cannot be empty")
	}
	if region == "" {
		return nil, errors.New("region cannot be empty")
	}

	config := NewSpeechTranslationConfig()
	config.SetProperty(SpeechServiceAuthorizationToken, authToken)
	config.SetProperty(SpeechServiceConnectionRegion, region)

	return config, nil
}

// AddTargetLanguage adds a language to the list of target languages for translation
func (c *SpeechTranslationConfig) AddTargetLanguage(language string) {
	// Check if language already exists in target languages
	for _, lang := range c.targetLanguages {
		if lang == language {
			return // Language already exists, no need to add it again
		}
	}

	c.targetLanguages = append(c.targetLanguages, language)
	c.SetProperty(SpeechServiceConnectionTranslationToLanguages, strings.Join(c.targetLanguages, ","))
}

// RemoveTargetLanguage removes a language from the list of target languages for translation
func (c *SpeechTranslationConfig) RemoveTargetLanguage(language string) {
	var newTargetLanguages []string
	for _, lang := range c.targetLanguages {
		if lang != language {
			newTargetLanguages = append(newTargetLanguages, lang)
		}
	}

	c.targetLanguages = newTargetLanguages
	c.SetProperty(SpeechServiceConnectionTranslationToLanguages, strings.Join(c.targetLanguages, ","))
}

// GetTargetLanguages returns the list of target languages for translation
func (c *SpeechTranslationConfig) GetTargetLanguages() []string {
	return c.targetLanguages
}

// SetVoiceName sets the voice to use for synthesized output
func (c *SpeechTranslationConfig) SetVoiceName(voiceName string) {
	c.voiceName = voiceName
	c.SetProperty(SpeechServiceConnectionTranslationVoice, voiceName)
}

// GetVoiceName returns the voice to use for synthesized output
func (c *SpeechTranslationConfig) GetVoiceName() string {
	return c.voiceName
}

// SetCustomModelCategoryID sets a Category ID that will be passed to the service
// Category ID is used to find the custom model
func (c *SpeechTranslationConfig) SetCustomModelCategoryID(categoryID string) {
	c.SetPropertyByName("CUSTOM_MODEL_CATEGORY_ID", categoryID)
}

// TranslationRecognitionResult defines the translation result
type TranslationRecognitionResult struct {
	// Common recognition result properties
	ResultID string
	Text     string
	Reason   ResultReason
	Offset   int64
	Duration time.Duration

	// Translation-specific properties
	Translations map[string]string // Maps target language to translated text
}

// TranslationSynthesisResult represents the voice output in the target language
type TranslationSynthesisResult struct {
	Audio  []byte       // The audio data of the synthesized output
	Reason ResultReason // The reason for the result
}

// EventArgs is a base interface for event arguments
type EventArgs interface {
	GetSessionID() string
}

// SessionEventArgs is the base type for session events
type SessionEventArgs struct {
	SessionID string
}

// GetSessionID returns the session ID
func (e *SessionEventArgs) GetSessionID() string {
	return e.SessionID
}

// RecognitionEventArgs contains data for recognition events
type RecognitionEventArgs struct {
	SessionEventArgs
	Offset int64
}

// TranslationRecognitionEventArgs contains data for translation recognition events
type TranslationRecognitionEventArgs struct {
	RecognitionEventArgs
	Result *TranslationRecognitionResult
}

// TranslationSynthesisEventArgs contains data for translation synthesis events
type TranslationSynthesisEventArgs struct {
	SessionEventArgs
	Result *TranslationSynthesisResult
}

// CancellationDetails contains details about why a result was canceled
type CancellationDetails struct {
	Reason       CancellationReason
	ErrorCode    CancellationErrorCode
	ErrorDetails string
}

// TranslationRecognitionCanceledEventArgs contains data for translation recognition canceled events
type TranslationRecognitionCanceledEventArgs struct {
	TranslationRecognitionEventArgs
	CancellationDetails *CancellationDetails
}

// EventCallback is a type for handling event callbacks
type EventCallback func(interface{})

// EventSignal handles connections to events
type EventSignal struct {
	callbacks []EventCallback
	mu        sync.RWMutex
}

// NewEventSignal creates a new event signal
func NewEventSignal() *EventSignal {
	return &EventSignal{
		callbacks: make([]EventCallback, 0),
	}
}

// Connect connects a callback to the event signal
func (s *EventSignal) Connect(callback EventCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callbacks = append(s.callbacks, callback)
}

// Disconnect disconnects all callbacks
func (s *EventSignal) Disconnect() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callbacks = make([]EventCallback, 0)
}

// Signal signals all connected callbacks
func (s *EventSignal) Signal(args interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, callback := range s.callbacks {
		callback(args)
	}
}

// TranslationRecognizer performs translation on speech input
type TranslationRecognizer struct {
	config              *SpeechTranslationConfig
	audioConfig         *AudioConfig
	properties          *PropertyCollection
	recognizing         *EventSignal
	recognized          *EventSignal
	canceled            *EventSignal
	synthesizing        *EventSignal
	sessionStarted      *EventSignal
	sessionStopped      *EventSignal
	speechStartDetected *EventSignal
	speechEndDetected   *EventSignal
	isContinuous        bool
	continuousRunning   bool
	continuousMutex     sync.Mutex
	stopCh              chan struct{}
}

// NewTranslationRecognizer creates a new translation recognizer
func NewTranslationRecognizer(translationConfig *SpeechTranslationConfig, audioConfig *AudioConfig) (*TranslationRecognizer, error) {
	if translationConfig == nil {
		return nil, fmt.Errorf("translation config cannot be nil")
	}

	// Use default audio config (default microphone) if none provided
	var err error
	if audioConfig == nil {
		audioConfig, err = NewAudioConfigFromDefaultMicrophone()
		if err != nil {
			return nil, fmt.Errorf("failed to create default audio config: %v", err)
		}
	}

	// Validate audio source
	if audioConfig.Source() == nil {
		return nil, fmt.Errorf("audio source cannot be nil")
	}

	// Validate that the source implements io.Reader
	if _, ok := audioConfig.Source().(io.Reader); !ok {
		return nil, fmt.Errorf("audio source must implement io.Reader")
	}

	recognizer := &TranslationRecognizer{
		config:              translationConfig,
		audioConfig:         audioConfig,
		properties:          NewPropertyCollection(),
		recognizing:         NewEventSignal(),
		recognized:          NewEventSignal(),
		canceled:            NewEventSignal(),
		synthesizing:        NewEventSignal(),
		sessionStarted:      NewEventSignal(),
		sessionStopped:      NewEventSignal(),
		speechStartDetected: NewEventSignal(),
		speechEndDetected:   NewEventSignal(),
		isContinuous:        false,
		continuousRunning:   false,
		stopCh:              make(chan struct{}),
	}

	// Copy properties from translation config
	propIDs := []PropertyID{
		SpeechServiceConnectionKey,
		SpeechServiceConnectionRegion,
		SpeechServiceConnectionEndpoint,
		SpeechServiceConnectionEndpointID,
		SpeechServiceAuthorizationToken,
		SpeechServiceConnectionTranslationToLanguages,
		SpeechServiceConnectionTranslationVoice,
		SpeechServiceConnectionRecoLanguage,
	}

	for _, id := range propIDs {
		val := translationConfig.GetProperty(id)
		if val != "" {
			recognizer.properties.SetProperty(id, val)
		}
	}

	return recognizer, nil
}

// RecognizeOnce performs a single recognition operation
func (r *TranslationRecognizer) RecognizeOnce(ctx context.Context) (*TranslationRecognitionResult, error) {
	if r.config.GetSubscriptionKey() == "" {
		return nil, errors.New("subscription key is not set")
	}

	// Signal session start
	r.raiseSessionStarted()

	// WebSocket接続を確立
	conn, err := r.connectToSpeechService()
	if err != nil {
		r.raiseCanceled(&CancellationDetails{
			Reason:       CancellationReasonError,
			ErrorCode:    CancellationErrorConnectionFailure,
			ErrorDetails: fmt.Sprintf("Failed to connect to Speech Service: %v", err),
		})
		return nil, err
	}
	defer conn.close()

	// Signal speech start detected
	r.raiseSpeechStartDetected()

	// オーディオデータの読み取り
	buffer := make([]byte, 8192)
	n, err := r.audioConfig.Source().(io.Reader).Read(buffer)
	if err != nil {
		if err != io.EOF {
			r.raiseCanceled(&CancellationDetails{
				Reason:       CancellationReasonError,
				ErrorCode:    CancellationErrorConnectionFailure,
				ErrorDetails: fmt.Sprintf("Error reading audio data: %v", err),
			})
			return nil, err
		}
	}

	if n > 0 {
		// オーディオデータの送信
		if err := conn.sendAudioData(buffer[:n]); err != nil {
			r.raiseCanceled(&CancellationDetails{
				Reason:       CancellationReasonError,
				ErrorCode:    CancellationErrorConnectionFailure,
				ErrorDetails: fmt.Sprintf("Error sending audio data: %v", err),
			})
			return nil, err
		}

		// 結果の受信
		result, err := conn.receiveResults()
		if err != nil {
			r.raiseCanceled(&CancellationDetails{
				Reason:       CancellationReasonError,
				ErrorCode:    CancellationErrorConnectionFailure,
				ErrorDetails: fmt.Sprintf("Error receiving results: %v", err),
			})
			return nil, err
		}

		// Signal speech end detected
		r.raiseSpeechEndDetected()

		// Signal the appropriate events
		r.raiseRecognizing(result)
		r.raiseRecognized(result)

		// Signal session stop
		r.raiseSessionStopped()

		return result, nil
	}

	return nil, errors.New("no audio data available")
}

// StartContinuousRecognitionAsync starts continuous recognition
func (r *TranslationRecognizer) StartContinuousRecognitionAsync(ctx context.Context) error {
	r.continuousMutex.Lock()
	defer r.continuousMutex.Unlock()

	log.Printf("[DEBUG] StartContinuousRecognitionAsync が呼び出されました")

	if r.continuousRunning {
		log.Printf("[DEBUG] 連続認識が既に実行中です")
		return errors.New("continuous recognition is already running")
	}

	r.continuousRunning = true
	r.stopCh = make(chan struct{})

	log.Printf("[DEBUG] continuousRecognitionWorker を起動します")
	go r.continuousRecognitionWorker(ctx)

	log.Printf("[DEBUG] StartContinuousRecognitionAsync が正常に完了しました")
	return nil
}

// StopContinuousRecognitionAsync stops continuous recognition
func (r *TranslationRecognizer) StopContinuousRecognitionAsync() error {
	r.continuousMutex.Lock()
	defer r.continuousMutex.Unlock()

	if !r.continuousRunning {
		return errors.New("continuous recognition is not running")
	}

	close(r.stopCh)
	r.continuousRunning = false

	return nil
}

// continuousRecognitionWorker handles the continuous recognition process
func (r *TranslationRecognizer) continuousRecognitionWorker(ctx context.Context) {
	log.Printf("[DEBUG] continuousRecognitionWorker が開始されました")

	// Signal session start
	r.raiseSessionStarted()

	// WebSocket接続を確立
	log.Printf("[DEBUG] Speech Serviceへの接続を試みています")
	conn, err := r.connectToSpeechService()
	if err != nil {
		log.Printf("[ERROR] Speech Serviceへの接続に失敗しました: %v", err)
		r.raiseCanceled(&CancellationDetails{
			Reason:       CancellationReasonError,
			ErrorCode:    CancellationErrorConnectionFailure,
			ErrorDetails: fmt.Sprintf("Failed to connect to Speech Service: %v", err),
		})
		return
	}
	defer conn.close()
	log.Printf("[DEBUG] Speech Serviceへの接続が確立されました: sourceLanguage=%s, targetLanguages=%v",
		r.config.GetSpeechRecognitionLanguage(), r.GetTargetLanguages())

	// Audio source setup
	log.Printf("[DEBUG] オーディオソースの設定: SourceType=%s", r.audioConfig.SourceType())
	var audioSource io.Reader
	switch r.audioConfig.SourceType() {
	case "Microphone":
		audioSource = r.audioConfig.Source().(io.Reader)
		log.Printf("[DEBUG] マイク入力をオーディオソースとして設定しました")
	case "Stream":
		audioSource = r.audioConfig.Source().(io.Reader)
		log.Printf("[DEBUG] ストリームをオーディオソースとして設定しました")
	case "File":
		audioSource = r.audioConfig.Source().(io.Reader)
		log.Printf("[DEBUG] ファイルをオーディオソースとして設定しました")
	case "PushStream":
		audioSource = r.audioConfig.Source().(io.Reader)
		log.Printf("[DEBUG] PushStreamをオーディオソースとして設定しました: %T", r.audioConfig.Source())
	default:
		log.Printf("[ERROR] サポートされていないオーディオソースタイプ: %s", r.audioConfig.SourceType())
		r.raiseCanceled(&CancellationDetails{
			Reason:       CancellationReasonError,
			ErrorCode:    CancellationErrorConnectionFailure,
			ErrorDetails: "Unsupported audio source type",
		})
		return
	}

	// オーディオデータを読み込むバッファ
	buffer := make([]byte, 8192) // 8KBのバッファ
	log.Printf("[DEBUG] 8KBのオーディオバッファを作成しました")

	// 音声レベルのログ出力用の変数
	lastLogTime := time.Now()
	logInterval := 500 * time.Millisecond // 500ミリ秒ごとにログを出力
	log.Printf("[DEBUG] 音声レベルログの間隔を %v に設定しました", logInterval)

	// データ読み取り統計情報
	var totalBytesRead int
	var readAttempts int
	var successfulReads int
	var logStats time.Time = time.Now()
	statsLogInterval := 5 * time.Second // 5秒ごとに統計情報をログ出力

	// エラー処理用のチャネル
	errCh := make(chan error, 1)
	log.Printf("[DEBUG] エラー処理用チャネルを作成しました")

	// 結果受信用のゴルーチン
	log.Printf("[DEBUG] 結果受信用ゴルーチンを開始します")
	go func() {
		for {
			log.Printf("[DEBUG] WebSocketから結果を待機中...")
			result, err := conn.receiveResults()
			if err != nil {
				log.Printf("[ERROR] 結果の受信中にエラーが発生: %v", err)
				errCh <- err
				return
			}

			if result != nil {
				log.Printf("[DEBUG] 認識結果を受信: Text=%s", result.Text)
				// イベントを発火
				r.raiseRecognizing(result)
				r.raiseRecognized(result)
			}
		}
	}()

	log.Printf("[DEBUG] 連続認識ループを開始します")
	// Continuous recognition loop
	for {
		select {
		case <-r.stopCh:
			// Stop requested
			log.Printf("[DEBUG] 停止リクエストを受信しました")
			r.raiseSessionStopped()
			return
		case <-ctx.Done():
			// Context canceled or timed out
			log.Printf("[DEBUG] コンテキストがキャンセルまたはタイムアウトしました")
			r.raiseSessionStopped()
			return
		case err := <-errCh:
			// エラーが発生した場合
			log.Printf("[ERROR] 連続認識中にエラーが発生: %v", err)
			r.raiseCanceled(&CancellationDetails{
				Reason:       CancellationReasonError,
				ErrorCode:    CancellationErrorConnectionFailure,
				ErrorDetails: fmt.Sprintf("Error in continuous recognition: %v", err),
			})
			return
		default:
			// オーディオデータの読み込み
			n, err := audioSource.Read(buffer)
			if err != nil {
				if err == io.EOF {
					// ファイル終端に達した場合
					log.Printf("[DEBUG] ファイル終端に達しました")
					r.raiseSessionStopped()
					return
				}
				// その他のエラー
				log.Printf("[ERROR] オーディオデータの読み込み中にエラーが発生: %v", err)
				r.raiseCanceled(&CancellationDetails{
					Reason:       CancellationReasonError,
					ErrorCode:    CancellationErrorConnectionFailure,
					ErrorDetails: fmt.Sprintf("Error reading audio data: %v", err),
				})
				return
			}

			if n > 0 {
				// データ読み取り統計情報を更新
				readAttempts++
				successfulReads++
				totalBytesRead += n

				// 定期的に統計情報をログ出力
				if time.Since(logStats) >= statsLogInterval {
					log.Printf("[STATS] オーディオ読み取り統計: 試行=%d, 成功=%d, 総バイト数=%d, 平均バイト数=%.2f/読み取り",
						readAttempts, successfulReads, totalBytesRead, float64(totalBytesRead)/float64(successfulReads))
					logStats = time.Now()
				}

				log.Printf("[DEBUG] オーディオデータを %d バイト読み込みました", n)

				// 音声レベルの計算と定期的なログ出力
				if time.Since(lastLogTime) >= logInterval {
					level := calculateAudioLevel(buffer[:n], n)
					log.Printf("マイク音声レベル: %d/100", level)
					lastLogTime = time.Now()
				}

				// オーディオデータの送信
				if err := conn.sendAudioData(buffer[:n]); err != nil {
					log.Printf("[ERROR] オーディオデータの送信中にエラーが発生: %v", err)
					r.raiseCanceled(&CancellationDetails{
						Reason:       CancellationReasonError,
						ErrorCode:    CancellationErrorConnectionFailure,
						ErrorDetails: fmt.Sprintf("Error sending audio data: %v", err),
					})
					return
				}
				log.Printf("[DEBUG] オーディオデータを送信しました")
			} else {
				log.Printf("[DEBUG] 読み込まれたオーディオデータがありません (n=0)")
			}

			// 短い遅延を入れて CPU 使用率を抑える
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// StartContinuousRecognition starts continuous recognition synchronously
func (r *TranslationRecognizer) StartContinuousRecognition(ctx context.Context) error {
	log.Printf("[DEBUG] StartContinuousRecognition が呼び出されました")
	return r.StartContinuousRecognitionAsync(ctx)
}

// StopContinuousRecognition stops continuous recognition synchronously
func (r *TranslationRecognizer) StopContinuousRecognition() error {
	return r.StopContinuousRecognitionAsync()
}

// Event properties

// Recognizing returns the event signal for recognizing events
func (r *TranslationRecognizer) Recognizing() *EventSignal {
	return r.recognizing
}

// Recognized returns the event signal for recognized events
func (r *TranslationRecognizer) Recognized() *EventSignal {
	return r.recognized
}

// Canceled returns the event signal for canceled events
func (r *TranslationRecognizer) Canceled() *EventSignal {
	return r.canceled
}

// Synthesizing returns the event signal for synthesizing events
func (r *TranslationRecognizer) Synthesizing() *EventSignal {
	return r.synthesizing
}

// SessionStarted returns the event signal for session started events
func (r *TranslationRecognizer) SessionStarted() *EventSignal {
	return r.sessionStarted
}

// SessionStopped returns the event signal for session stopped events
func (r *TranslationRecognizer) SessionStopped() *EventSignal {
	return r.sessionStopped
}

// SpeechStartDetected returns the event signal for speech start detected events
func (r *TranslationRecognizer) SpeechStartDetected() *EventSignal {
	return r.speechStartDetected
}

// SpeechEndDetected returns the event signal for speech end detected events
func (r *TranslationRecognizer) SpeechEndDetected() *EventSignal {
	return r.speechEndDetected
}

// Event raisers

func (r *TranslationRecognizer) raiseSessionStarted() {
	args := &SessionEventArgs{
		SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
	}
	r.sessionStarted.Signal(args)
}

func (r *TranslationRecognizer) raiseSessionStopped() {
	args := &SessionEventArgs{
		SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
	}
	r.sessionStopped.Signal(args)
}

func (r *TranslationRecognizer) raiseSpeechStartDetected() {
	args := &RecognitionEventArgs{
		SessionEventArgs: SessionEventArgs{
			SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
		},
		Offset: time.Now().UnixNano(),
	}
	r.speechStartDetected.Signal(args)
}

func (r *TranslationRecognizer) raiseSpeechEndDetected() {
	args := &RecognitionEventArgs{
		SessionEventArgs: SessionEventArgs{
			SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
		},
		Offset: time.Now().UnixNano(),
	}
	r.speechEndDetected.Signal(args)
}

func (r *TranslationRecognizer) raiseRecognizing(result *TranslationRecognitionResult) {
	args := &TranslationRecognitionEventArgs{
		RecognitionEventArgs: RecognitionEventArgs{
			SessionEventArgs: SessionEventArgs{
				SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
			},
			Offset: result.Offset,
		},
		Result: result,
	}
	r.recognizing.Signal(args)
}

func (r *TranslationRecognizer) raiseRecognized(result *TranslationRecognitionResult) {
	args := &TranslationRecognitionEventArgs{
		RecognitionEventArgs: RecognitionEventArgs{
			SessionEventArgs: SessionEventArgs{
				SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
			},
			Offset: result.Offset,
		},
		Result: result,
	}
	r.recognized.Signal(args)
}

func (r *TranslationRecognizer) raiseCanceled(details *CancellationDetails) {
	result := &TranslationRecognitionResult{
		ResultID: fmt.Sprintf("canceled_%d", time.Now().UnixNano()),
		Reason:   ResultReasonCanceled,
		Offset:   time.Now().UnixNano(),
	}

	args := &TranslationRecognitionCanceledEventArgs{
		TranslationRecognitionEventArgs: TranslationRecognitionEventArgs{
			RecognitionEventArgs: RecognitionEventArgs{
				SessionEventArgs: SessionEventArgs{
					SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
				},
				Offset: result.Offset,
			},
			Result: result,
		},
		CancellationDetails: details,
	}
	r.canceled.Signal(args)
}

func (r *TranslationRecognizer) raiseSynthesizing(audio []byte) {
	result := &TranslationSynthesisResult{
		Audio:  audio,
		Reason: ResultReasonTranslatedSpeech,
	}

	args := &TranslationSynthesisEventArgs{
		SessionEventArgs: SessionEventArgs{
			SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
		},
		Result: result,
	}

	r.synthesizing.Signal(args)
}

// GetTargetLanguages returns the list of target languages for translation
func (r *TranslationRecognizer) GetTargetLanguages() []string {
	languagesStr := r.properties.GetProperty(SpeechServiceConnectionTranslationToLanguages)
	if languagesStr == "" {
		return []string{}
	}

	return strings.Split(languagesStr, ",")
}

// AddTargetLanguage adds a language to the list of target languages for translation
func (r *TranslationRecognizer) AddTargetLanguage(language string) {
	languages := r.GetTargetLanguages()

	// Check if language already exists
	for _, lang := range languages {
		if lang == language {
			return // Already exists
		}
	}

	// Add the language
	languages = append(languages, language)
	r.properties.SetProperty(
		SpeechServiceConnectionTranslationToLanguages,
		strings.Join(languages, ","),
	)
}

// RemoveTargetLanguage removes a language from the list of target languages for translation
func (r *TranslationRecognizer) RemoveTargetLanguage(language string) {
	languages := r.GetTargetLanguages()
	var newLanguages []string

	for _, lang := range languages {
		if lang != language {
			newLanguages = append(newLanguages, lang)
		}
	}

	r.properties.SetProperty(
		SpeechServiceConnectionTranslationToLanguages,
		strings.Join(newLanguages, ","),
	)
}

// Close cleans up resources
func (r *TranslationRecognizer) Close() error {
	// Stop continuous recognition if running
	if r.continuousRunning {
		r.StopContinuousRecognition()
	}

	// Clean up event signals
	r.recognizing.Disconnect()
	r.recognized.Disconnect()
	r.canceled.Disconnect()
	r.synthesizing.Disconnect()
	r.sessionStarted.Disconnect()
	r.sessionStopped.Disconnect()
	r.speechStartDetected.Disconnect()
	r.speechEndDetected.Disconnect()

	// Close audio config
	if r.audioConfig != nil {
		return r.audioConfig.Close()
	}

	return nil
}

// speechServiceConnection はAzure Speech ServiceのWebSocket接続を管理します
type speechServiceConnection struct {
	conn           *websocket.Conn
	authToken      string
	region         string
	languages      []string
	sourceLanguage string
}

// connectToSpeechService はAzure Speech ServiceのWebSocket APIに接続します
func (r *TranslationRecognizer) connectToSpeechService() (*speechServiceConnection, error) {
	log.Printf("Speech Serviceへの接続開始: region=%s", r.config.GetRegion())

	dialer := websocket.Dialer{
		EnableCompression: true,
	}

	// ヘッダーの準備
	header := http.Header{}
	authToken := r.config.GetAuthorizationToken()
	if authToken == "" {
		authToken = r.config.GetSubscriptionKey()
	}

	if authToken == "" {
		return nil, fmt.Errorf("認証情報が設定されていません")
	}

	header.Add("Authorization", "Bearer "+authToken)
	header.Add("Ocp-Apim-Subscription-Key", os.Getenv("SPEECH_SERVICE_KEY"))
	header.Add("X-ConnectionId", uuid.New().String())

	// WebSocket URLの構築
	wsURL := fmt.Sprintf("wss://%s.stt.speech.microsoft.com/speech/universal/v2", r.config.GetRegion())
	log.Printf("Speech Service WebSocket URL: %s", wsURL)

	// WebSocket接続の確立
	log.Printf("WebSocket接続を試行中...")
	conn, resp, err := dialer.Dial(wsURL, header)
	if err != nil {
		if resp != nil {
			log.Printf("接続エラー - Status: %d, Headers: %v", resp.StatusCode, resp.Header)
		}
		return nil, fmt.Errorf("failed to connect to Speech Service: %v", err)
	}
	log.Printf("Speech ServiceのWebSocket接続が確立されました")

	return &speechServiceConnection{
		conn:           conn,
		authToken:      authToken,
		region:         r.config.GetRegion(),
		languages:      r.GetTargetLanguages(),
		sourceLanguage: r.config.GetSpeechRecognitionLanguage(),
	}, nil
}

// sendAudioData はオーディオデータをWebSocket経由で送信します
func (sc *speechServiceConnection) sendAudioData(data []byte) error {
	log.Printf("[DEBUG] Speech Serviceに送信するオーディオデータ: %d バイト", len(data))

	requestID := uuid.New().String()

	// 言語コードの正規化と検証
	normalizedSourceLang := normalizeLanguageCode(sc.sourceLanguage)
	if normalizedSourceLang == "" {
		return fmt.Errorf("invalid source language code: %s", sc.sourceLanguage)
	}
	log.Printf("[DEBUG] 正規化されたソース言語: %s (元の値: %s)", normalizedSourceLang, sc.sourceLanguage)

	// ターゲット言語の正規化と検証
	normalizedTargetLangs := make([]string, 0, len(sc.languages))
	for _, lang := range sc.languages {
		normalized := normalizeLanguageCode(lang)
		if normalized == "" {
			return fmt.Errorf("invalid target language code: %s", lang)
		}
		normalizedTargetLangs = append(normalizedTargetLangs, normalized)
	}
	log.Printf("[DEBUG] 正規化されたターゲット言語: %v", normalizedTargetLangs)

	// WebSocket設定メッセージの構築
	configMsg := map[string]interface{}{
		"context": map[string]interface{}{
			"system": map[string]interface{}{
				"name":    "SpeechSDK",
				"version": "1.28.0",
				"build":   "JavaScript",
			},
		},
		"config": map[string]interface{}{
			"speechConfig": map[string]interface{}{
				"speechRecognitionLanguage":    normalizedSourceLang,
				"translationLanguages":         normalizedTargetLangs,
				"sourceLanguageForTranslation": normalizedSourceLang,
				"features": map[string]interface{}{
					"enableTranslation": true,
				},
			},
			"input": map[string]interface{}{
				"format": "audio/x-wav",
			},
		},
	}

	// 設定メッセージをJSONに変換
	configBytes, err := json.Marshal(configMsg)
	if err != nil {
		log.Printf("[ERROR] 設定メッセージのJSONエンコードに失敗: %v", err)
		return err
	}

	log.Printf("[DEBUG] Speech Service設定: %s", string(configBytes))

	// Speech Service用のヘッダー形式でメッセージを構築
	configHeader := fmt.Sprintf("Path: speech.config\r\nX-RequestId: %s\r\nX-Timestamp: %s\r\nContent-Type: application/json\r\n\r\n%s",
		requestID,
		time.Now().UTC().Format(time.RFC3339),
		configBytes)

	// 設定メッセージを送信
	if err := sc.conn.WriteMessage(websocket.TextMessage, []byte(configHeader)); err != nil {
		log.Printf("[ERROR] 設定メッセージの送信に失敗: %v", err)
		return err
	}

	// オーディオメッセージのヘッダーを構築
	audioHeader := fmt.Sprintf("Path: audio\r\nX-RequestId: %s\r\nX-Timestamp: %s\r\nContent-Type: audio/x-wav\r\n\r\n",
		requestID,
		time.Now().UTC().Format(time.RFC3339))

	// オーディオヘッダーを送信
	if err := sc.conn.WriteMessage(websocket.TextMessage, []byte(audioHeader)); err != nil {
		log.Printf("[ERROR] オーディオヘッダーの送信に失敗: %v", err)
		return err
	}

	// オーディオデータを送信
	if err := sc.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Printf("[ERROR] オーディオデータの送信に失敗: %v", err)
		return err
	}

	log.Printf("[DEBUG] メッセージを送信完了 - RequestID: %s, DataSize: %d bytes", requestID, len(data))
	return nil
}

// receiveResults は認識結果を受信します
func (sc *speechServiceConnection) receiveResults() (*TranslationRecognitionResult, error) {
	messageType, message, err := sc.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] クライアントからメッセージを受信: type=%d, dataSize=%d bytes", messageType, len(message))

	// テキストメッセージの場合（ヘッダーとJSONボディ）
	if messageType == websocket.TextMessage {
		// メッセージをヘッダーとボディに分割
		parts := strings.Split(string(message), "\r\n\r\n")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid message format: expected header and body, got %d parts", len(parts))
		}

		headers := parts[0]
		body := parts[1]

		log.Printf("[DEBUG] 受信したヘッダー:\n%s", headers)
		log.Printf("[DEBUG] 受信したボディ:\n%s", body)

		// JSONをパース
		var response map[string]interface{}
		if err := json.Unmarshal([]byte(body), &response); err != nil {
			return nil, fmt.Errorf("JSON parse error: %v", err)
		}

		// レスポンスタイプをチェック - Pathヘッダーを確認
		var messagePath string
		headerLines := strings.Split(headers, "\r\n")
		for _, line := range headerLines {
			if strings.HasPrefix(line, "Path:") {
				messagePath = strings.TrimPrefix(line, "Path:")
				messagePath = strings.TrimSpace(messagePath)
				break
			}
		}

		log.Printf("[DEBUG] Message path: %s", messagePath)

		// 異なるメッセージタイプを処理
		switch messagePath {
		case "turn.start":
			// ターンスタートの処理 - 必要に応じてログを出力
			log.Printf("[DEBUG] Turn started with context: %s", body)
			return nil, nil
		case "speech.phrase":
			// 音声認識結果の処理
			if response["type"] == "final" {
				result := &TranslationRecognitionResult{
					ResultID:     fmt.Sprintf("result_%d", time.Now().UnixNano()),
					Reason:       ResultReasonTranslatedSpeech,
					Offset:       time.Now().UnixNano(),
					Duration:     1 * time.Second,
					Translations: make(map[string]string),
				}

				// 認識テキストの取得
				if nbest, ok := response["NBest"].([]interface{}); ok && len(nbest) > 0 {
					if firstResult, ok := nbest[0].(map[string]interface{}); ok {
						if display, ok := firstResult["Display"].(string); ok {
							result.Text = display
						}
					}
				}

				// 翻訳結果の取得
				if translations, ok := response["Translations"].(map[string]interface{}); ok {
					for lang, text := range translations {
						if textStr, ok := text.(string); ok {
							result.Translations[lang] = textStr
						}
					}
				}

				return result, nil
			}
		}
	}

	// 他のメッセージタイプやレスポンスタイプの場合はnilを返す
	return nil, nil
}

// close はWebSocket接続を閉じます
func (sc *speechServiceConnection) close() error {
	return sc.conn.Close()
}

// calculateAudioLevel は音声バッファから平均音声レベル（0-100の範囲）を計算します
func calculateAudioLevel(buffer []byte, n int) int {
	if n == 0 {
		log.Printf("DEBUG: 音声バッファが空です（サイズ=0）")
		return 0
	}

	// 音声データをint16に変換
	var sum int64
	for i := 0; i < n-1; i += 2 {
		// リトルエンディアンでint16に変換
		value := int16(buffer[i]) | (int16(buffer[i+1]) << 8)
		if value < 0 {
			value = -value // 絶対値を取る
		}
		sum += int64(value)
	}

	// 平均振幅を計算し、0-100の範囲にスケーリング
	avgAmplitude := sum / int64(n/2)
	// int16の最大値は32767なので、その値で割って0-100のスケールに変換
	level := int((avgAmplitude * 100) / 32767)

	// バッファ内の最初の数バイトをデバッグのために表示
	var bytesStr string
	maxBytes := 16
	if n < maxBytes {
		maxBytes = n
	}
	for i := 0; i < maxBytes; i++ {
		bytesStr += fmt.Sprintf("%02x ", buffer[i])
	}
	log.Printf("DEBUG: 音声バッファ先頭バイト: %s, バッファサイズ: %d, 平均振幅: %d, レベル: %d/100",
		bytesStr, n, avgAmplitude, level)

	return level
}

// normalizeLanguageCode は言語コードをBCP-47形式に正規化します
func normalizeLanguageCode(lang string) string {
	// 空白を削除し、小文字に変換
	lang = strings.TrimSpace(lang)

	// 既にBCP-47形式（xx-XX）の場合はバリデーションを行う
	if strings.Contains(lang, "-") {
		parts := strings.Split(lang, "-")
		if len(parts) != 2 {
			log.Printf("[WARNING] 無効な言語コード形式: %s", lang)
			return ""
		}
		// 言語コードは小文字、地域コードは大文字に正規化
		langCode := strings.ToLower(parts[0])
		regionCode := strings.ToUpper(parts[1])
		return fmt.Sprintf("%s-%s", langCode, regionCode)
	}

	// 言語コードのマッピング
	langMap := map[string]string{
		"ja": "ja-JP",
		"en": "en-US",
		"zh": "zh-CN",
		"ko": "ko-KR",
		"es": "es-ES",
		"fr": "fr-FR",
		"de": "de-DE",
		"it": "it-IT",
		"pt": "pt-BR",
		"ru": "ru-RU",
		"ar": "ar-SA",
		"hi": "hi-IN",
		"th": "th-TH",
		"vi": "vi-VN",
		"id": "id-ID",
		"ms": "ms-MY",
	}

	// 小文字に変換してマッピングを検索
	if normalized, ok := langMap[strings.ToLower(lang)]; ok {
		return normalized
	}

	log.Printf("[WARNING] サポートされていない言語コード: %s", lang)
	return ""
}
