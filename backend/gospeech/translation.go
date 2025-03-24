// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license.

package gospeech

import (
	"context"
	"encoding/base64"
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

	if r.continuousRunning {
		return errors.New("continuous recognition is already running")
	}

	r.continuousRunning = true
	r.stopCh = make(chan struct{})

	go r.continuousRecognitionWorker(ctx)

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
		return
	}
	defer conn.close()

	// Audio source setup
	var audioSource io.Reader
	switch r.audioConfig.SourceType() {
	case "Microphone":
		audioSource = r.audioConfig.Source().(io.Reader)
	case "Stream":
		audioSource = r.audioConfig.Source().(io.Reader)
	case "File":
		audioSource = r.audioConfig.Source().(io.Reader)
	default:
		r.raiseCanceled(&CancellationDetails{
			Reason:       CancellationReasonError,
			ErrorCode:    CancellationErrorConnectionFailure,
			ErrorDetails: "Unsupported audio source type",
		})
		return
	}

	// オーディオデータを読み込むバッファ
	buffer := make([]byte, 8192) // 8KBのバッファ

	// エラー処理用のチャネル
	errCh := make(chan error, 1)

	// 結果受信用のゴルーチン
	go func() {
		for {
			result, err := conn.receiveResults()
			if err != nil {
				errCh <- err
				return
			}

			if result != nil {
				// イベントを発火
				r.raiseRecognizing(result)
				r.raiseRecognized(result)
			}
		}
	}()

	// Continuous recognition loop
	for {
		select {
		case <-r.stopCh:
			// Stop requested
			r.raiseSessionStopped()
			return
		case <-ctx.Done():
			// Context canceled or timed out
			r.raiseSessionStopped()
			return
		case err := <-errCh:
			// エラーが発生した場合
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
					r.raiseSessionStopped()
					return
				}
				// その他のエラー
				r.raiseCanceled(&CancellationDetails{
					Reason:       CancellationReasonError,
					ErrorCode:    CancellationErrorConnectionFailure,
					ErrorDetails: fmt.Sprintf("Error reading audio data: %v", err),
				})
				return
			}

			if n > 0 {
				// オーディオデータの送信
				if err := conn.sendAudioData(buffer[:n]); err != nil {
					r.raiseCanceled(&CancellationDetails{
						Reason:       CancellationReasonError,
						ErrorCode:    CancellationErrorConnectionFailure,
						ErrorDetails: fmt.Sprintf("Error sending audio data: %v", err),
					})
					return
				}
			}

			// 短い遅延を入れて CPU 使用率を抑える
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// StartContinuousRecognition starts continuous recognition synchronously
func (r *TranslationRecognizer) StartContinuousRecognition(ctx context.Context) error {
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
	// Speech ServiceのWebSocket APIで必要なヘッダー情報を含むメッセージを作成
	message := map[string]interface{}{
		"audio": map[string]interface{}{
			"data": base64.StdEncoding.EncodeToString(data),
		},
		"context": map[string]interface{}{
			"sourceLanguage":  sc.sourceLanguage,
			"targetLanguages": sc.languages,
		},
	}

	return sc.conn.WriteJSON(message)
}

// receiveResults は認識結果を受信します
func (sc *speechServiceConnection) receiveResults() (*TranslationRecognitionResult, error) {
	_, message, err := sc.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// メッセージをパースして結果を生成
	var response map[string]interface{}
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, err
	}

	// レスポンスから認識結果を構築
	result := &TranslationRecognitionResult{
		ResultID:     fmt.Sprintf("result_%d", time.Now().UnixNano()),
		Text:         response["recognizedText"].(string),
		Reason:       ResultReasonTranslatedSpeech,
		Offset:       time.Now().UnixNano(),
		Duration:     1 * time.Second,
		Translations: make(map[string]string),
	}

	// 翻訳結果があれば追加
	if translations, ok := response["translations"].(map[string]interface{}); ok {
		for lang, text := range translations {
			result.Translations[lang] = text.(string)
		}
	}

	return result, nil
}

// close はWebSocket接続を閉じます
func (sc *speechServiceConnection) close() error {
	return sc.conn.Close()
}
