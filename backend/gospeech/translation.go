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

	log.Printf("[DEBUG] StartContinuousRecognitionAsync called")

	if r.continuousRunning {
		log.Printf("[DEBUG] Continuous recognition is already running")
		return errors.New("continuous recognition is already running")
	}

	r.continuousRunning = true
	r.stopCh = make(chan struct{})

	log.Printf("[DEBUG] Launching continuousRecognitionWorker")
	go r.continuousRecognitionWorker(ctx)

	log.Printf("[DEBUG] StartContinuousRecognitionAsync completed successfully")
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
	log.Printf("[DEBUG] continuousRecognitionWorker started")

	// Signal session start
	r.raiseSessionStarted()

	// WebSocket接続を確立
	log.Printf("[DEBUG] Attempting to connect to Speech Service")
	conn, err := r.connectToSpeechService()
	if err != nil {
		log.Printf("[ERROR] Failed to connect to Speech Service: %v", err)
		r.raiseCanceled(&CancellationDetails{
			Reason:       CancellationReasonError,
			ErrorCode:    CancellationErrorConnectionFailure,
			ErrorDetails: fmt.Sprintf("Failed to connect to Speech Service: %v", err),
		})
		return
	}
	defer conn.close()
	log.Printf("[DEBUG] Connection to Speech Service established: sourceLanguage=%s, targetLanguages=%v",
		r.config.GetSpeechRecognitionLanguage(), r.GetTargetLanguages())

	// Audio source setup
	log.Printf("[DEBUG] Audio source configuration: SourceType=%s", r.audioConfig.SourceType())
	var audioSource io.Reader
	switch r.audioConfig.SourceType() {
	case "Microphone":
		audioSource = r.audioConfig.Source().(io.Reader)
		log.Printf("[DEBUG] Microphone input set as audio source")
	case "Stream":
		audioSource = r.audioConfig.Source().(io.Reader)
		log.Printf("[DEBUG] Stream set as audio source")
	case "File":
		audioSource = r.audioConfig.Source().(io.Reader)
		log.Printf("[DEBUG] File set as audio source")
	case "PushStream":
		audioSource = r.audioConfig.Source().(io.Reader)
		log.Printf("[DEBUG] PushStream set as audio source: %T", r.audioConfig.Source())
	default:
		log.Printf("[ERROR] Unsupported audio source type: %s", r.audioConfig.SourceType())
		r.raiseCanceled(&CancellationDetails{
			Reason:       CancellationReasonError,
			ErrorCode:    CancellationErrorConnectionFailure,
			ErrorDetails: "Unsupported audio source type",
		})
		return
	}

	// オーディオデータを読み込むバッファ
	buffer := make([]byte, 8192) // 8KBのバッファ
	log.Printf("[DEBUG] Created 8KB audio buffer")

	// 音声レベルのログ出力用の変数
	lastLogTime := time.Now()
	logInterval := 500 * time.Millisecond // 500ミリ秒ごとにログを出力
	log.Printf("[DEBUG] Set voice level log interval to %v", logInterval)

	// データ読み取り統計情報
	var totalBytesRead int
	var readAttempts int
	var successfulReads int
	var logStats time.Time = time.Now()
	statsLogInterval := 5 * time.Second // 5秒ごとに統計情報をログ出力

	// エラー処理用のチャネル
	errCh := make(chan error, 1)
	log.Printf("[DEBUG] Created channel for error handling")

	// 結果受信用のゴルーチン
	log.Printf("[DEBUG] Starting goroutine for receiving results")
	go func() {
		for {
			log.Printf("[DEBUG] Waiting for results from WebSocket...")
			result, err := conn.receiveResults()
			if err != nil {
				log.Printf("[ERROR] Error occurred while receiving results: %v", err)
				errCh <- err
				return
			}

			if result != nil {
				log.Printf("[DEBUG] Received recognition result: Text=%s", result.Text)
				// イベントを発火
				r.raiseRecognizing(result)
				r.raiseRecognized(result)
			}
		}
	}()

	log.Printf("[DEBUG] Starting continuous recognition loop")
	// Continuous recognition loop
	for {
		select {
		case <-r.stopCh:
			// Stop requested
			log.Printf("[DEBUG] Stop request received")
			r.raiseSessionStopped()
			return
		case <-ctx.Done():
			// Context canceled or timed out
			log.Printf("[DEBUG] Context was canceled or timed out")
			r.raiseSessionStopped()
			return
		case err := <-errCh:
			// エラーが発生した場合
			log.Printf("[ERROR] Error occurred during continuous recognition: %v", err)
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
					log.Printf("[DEBUG] Reached end of file")
					r.raiseSessionStopped()
					return
				}
				// その他のエラー
				log.Printf("[ERROR] Error while reading audio data: %v", err)
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
					log.Printf("[STATS] Audio reading statistics: attempts=%d, successful=%d, totalBytes=%d, avgBytes=%.2f/read",
						readAttempts, successfulReads, totalBytesRead, float64(totalBytesRead)/float64(successfulReads))
					logStats = time.Now()
				}

				log.Printf("[DEBUG] Read %d bytes of audio data", n)

				// 音声レベルの計算と定期的なログ出力
				if time.Since(lastLogTime) >= logInterval {
					level := calculateAudioLevel(buffer[:n], n)
					log.Printf("Microphone audio level: %d/100", level)
					lastLogTime = time.Now()
				}

				// オーディオデータの送信
				if err := conn.sendAudioData(buffer[:n]); err != nil {
					log.Printf("[ERROR] Error while sending audio data: %v", err)
					r.raiseCanceled(&CancellationDetails{
						Reason:       CancellationReasonError,
						ErrorCode:    CancellationErrorConnectionFailure,
						ErrorDetails: fmt.Sprintf("Error sending audio data: %v", err),
					})
					return
				}
				log.Printf("[DEBUG] Audio data sent")
			} else {
				log.Printf("[DEBUG] No audio data read (n=0)")
			}

			// 短い遅延を入れて CPU 使用率を抑える
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// StartContinuousRecognition starts continuous recognition synchronously
func (r *TranslationRecognizer) StartContinuousRecognition(ctx context.Context) error {
	log.Printf("[DEBUG] StartContinuousRecognition called")
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

// connectToSpeechService connects to the Azure Speech Service WebSocket API
func (r *TranslationRecognizer) connectToSpeechService() (*speechServiceConnection, error) {
	log.Printf("[DEBUG] Speech Service connection start: region=%s", r.config.GetRegion())

	dialer := websocket.Dialer{
		EnableCompression: true,
	}

	// Prepare headers
	header := http.Header{}
	authToken := r.config.GetAuthorizationToken()
	if authToken == "" {
		authToken = r.config.GetSubscriptionKey()
	}

	if authToken == "" {
		return nil, fmt.Errorf("authentication information is not configured")
	}

	header.Add("Authorization", "Bearer "+authToken)
	header.Add("Ocp-Apim-Subscription-Key", os.Getenv("SPEECH_SERVICE_KEY"))
	header.Add("X-ConnectionId", uuid.New().String())

	// Construct WebSocket URL
	wsURL := fmt.Sprintf("wss://%s.stt.speech.microsoft.com/speech/universal/v2", r.config.GetRegion())
	log.Printf("[DEBUG] Speech Service WebSocket URL: %s", wsURL)

	// Establish WebSocket connection
	log.Printf("[DEBUG] Attempting WebSocket connection...")
	conn, resp, err := dialer.Dial(wsURL, header)
	if err != nil {
		if resp != nil {
			log.Printf("Connection error - Status: %d, Headers: %v", resp.StatusCode, resp.Header)
		}
		return nil, fmt.Errorf("failed to connect to Speech Service: %v", err)
	}
	log.Printf("WebSocket connection to Speech Service established")

	return &speechServiceConnection{
		conn:           conn,
		authToken:      authToken,
		region:         r.config.GetRegion(),
		languages:      r.GetTargetLanguages(),
		sourceLanguage: r.config.GetSpeechRecognitionLanguage(),
	}, nil
}

// sendAudioData sends audio data via WebSocket
func (sc *speechServiceConnection) sendAudioData(data []byte) error {
	log.Printf("[DEBUG] Audio data to send to Speech Service: %d bytes", len(data))

	requestID := uuid.New().String()

	// Normalize and validate language codes
	normalizedSourceLang := normalizeLanguageCode(sc.sourceLanguage, true)
	if normalizedSourceLang == "" {
		return fmt.Errorf("invalid source language code: %s", sc.sourceLanguage)
	}
	log.Printf("[DEBUG] Normalized source language: %s (original: %s)", normalizedSourceLang, sc.sourceLanguage)

	// Normalize and validate target languages
	normalizedTargetLangs := make([]string, 0, len(sc.languages))
	for _, lang := range sc.languages {
		normalized := normalizeLanguageCode(lang, false)
		if normalized == "" {
			return fmt.Errorf("invalid target language code: %s", lang)
		}
		normalizedTargetLangs = append(normalizedTargetLangs, normalized)
	}
	log.Printf("[DEBUG] Normalized target languages: %v", normalizedTargetLangs)

	// Construct WebSocket configuration message
	configMsg := map[string]interface{}{
		"context": map[string]interface{}{
			"system": map[string]interface{}{
				"name":    "SpeechSDK",
				"version": "1.30.0",
				"build":   "Go",
			},
		},
		"config": map[string]interface{}{
			"speechConfig": map[string]interface{}{
				"speechRecognitionLanguage":    normalizedSourceLang,
				"translationLanguages":         normalizedTargetLangs,
				"sourceLanguageForTranslation": normalizedSourceLang,
				"features": map[string]interface{}{
					"enableTranslation":   true,
					"wordLevelTimestamps": true,
					"punctuation":         "explicit",
				},
				"profanity":               "masked",
				"timeToDetectEndOfSpeech": "1500",
				"scenarios":               []string{"conversation"},
			},
			"input": map[string]interface{}{
				"format": "audio/x-wav",
				"audioParameters": map[string]interface{}{
					"sampleRate": 16000,
				},
			},
		},
	}

	// Convert configuration message to JSON
	configBytes, err := json.Marshal(configMsg)
	if err != nil {
		log.Printf("[ERROR] Failed to JSON encode configuration message: %v", err)
		return err
	}

	log.Printf("[DEBUG] Speech Service configuration: %s", string(configBytes))

	// Construct message in Speech Service header format
	configHeader := fmt.Sprintf("Path: speech.config\r\nX-RequestId: %s\r\nX-Timestamp: %s\r\nContent-Type: application/json\r\n\r\n%s",
		requestID,
		time.Now().UTC().Format(time.RFC3339),
		configBytes)

	// Send configuration message
	if err := sc.conn.WriteMessage(websocket.TextMessage, []byte(configHeader)); err != nil {
		log.Printf("[ERROR] Failed to send configuration message: %v", err)
		return err
	}

	// Construct audio message header
	audioHeader := fmt.Sprintf("Path: audio\r\nX-RequestId: %s\r\nX-Timestamp: %s\r\nContent-Type: audio/x-wav\r\n\r\n",
		requestID,
		time.Now().UTC().Format(time.RFC3339))

	// Send audio header
	if err := sc.conn.WriteMessage(websocket.TextMessage, []byte(audioHeader)); err != nil {
		log.Printf("[ERROR] Failed to send audio header: %v", err)
		return err
	}

	// Send audio data
	if err := sc.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Printf("[ERROR] Failed to send audio data: %v", err)
		return err
	}

	log.Printf("[DEBUG] Message sent successfully - RequestID: %s, DataSize: %d bytes", requestID, len(data))
	return nil
}

// receiveResults は認識結果を受信します
func (sc *speechServiceConnection) receiveResults() (*TranslationRecognitionResult, error) {
	messageType, message, err := sc.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Message received from client: type=%d, dataSize=%d bytes", messageType, len(message))

	// テキストメッセージの場合（ヘッダーとJSONボディ）
	if messageType == websocket.TextMessage {
		// メッセージをヘッダーとボディに分割
		parts := strings.Split(string(message), "\r\n\r\n")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid message format: expected header and body, got %d parts", len(parts))
		}

		headers := parts[0]
		body := parts[1]

		log.Printf("[DEBUG] Received headers:\n%s", headers)
		log.Printf("[DEBUG] Received body:\n%s", body)

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
		log.Printf("DEBUG: Audio buffer is empty (size=0)")
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
	log.Printf("DEBUG: Audio buffer head bytes: %s, buffer size: %d, average amplitude: %d, level: %d/100",
		bytesStr, n, avgAmplitude, level)

	return level
}

// normalizeLanguageCode normalizes language codes to BCP-47 format or simple language code
func normalizeLanguageCode(lang string, isSourceLanguage bool) string {
	// Remove spaces and convert to lowercase
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return ""
	}

	// Language code mapping
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

	if isSourceLanguage {
		// Source language requires full BCP-47 format
		if strings.Contains(lang, "-") {
			parts := strings.Split(lang, "-")
			if len(parts) != 2 {
				log.Printf("[WARNING] Invalid language code format: %s", lang)
				return ""
			}
			// Language code in lowercase, region code in uppercase
			langCode := strings.ToLower(parts[0])
			regionCode := strings.ToUpper(parts[1])
			return fmt.Sprintf("%s-%s", langCode, regionCode)
		}

		// Use mapping to convert to full format
		if normalized, ok := langMap[strings.ToLower(lang)]; ok {
			return normalized
		}
	} else {
		// Target language only needs language code
		if strings.Contains(lang, "-") {
			// If contains hyphen, get only the language part
			parts := strings.Split(lang, "-")
			return strings.ToLower(parts[0])
		}

		// If already just a language code, return it in lowercase
		return strings.ToLower(lang)
	}

	log.Printf("[WARNING] Unsupported language code: %s", lang)
	return ""
}
