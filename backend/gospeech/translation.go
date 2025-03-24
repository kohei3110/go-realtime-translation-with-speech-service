// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license.

package gospeech

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
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
		return nil, errors.New("translation config cannot be nil")
	}

	// Use default audio config (default microphone) if none provided
	if audioConfig == nil {
		var err error
		audioConfig, err = NewAudioConfigFromDefaultMicrophone()
		if err != nil {
			return nil, fmt.Errorf("failed to create default microphone config: %v", err)
		}
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
	// Create a context with timeout if not already done
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	// Signal session start
	r.raiseSessionStarted()

	// In a real implementation, this would connect to the Speech Service
	// Here we're simulating basic functionality

	// Signal speech start detected
	r.raiseSpeechStartDetected()

	// Simulate processing time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(1 * time.Second):
		// Continue processing
	}

	// Generate a result based on the configured languages
	result := &TranslationRecognitionResult{
		ResultID:     "SimulatedResultID",
		Text:         "Hello, how are you?",
		Reason:       ResultReasonTranslatedSpeech,
		Offset:       0,
		Duration:     1 * time.Second,
		Translations: make(map[string]string),
	}

	// Add a translation for each target language
	targetLangs := r.config.GetTargetLanguages()
	for _, lang := range targetLangs {
		switch lang {
		case "ja":
			result.Translations[lang] = "こんにちは、お元気ですか？"
		case "es":
			result.Translations[lang] = "Hola, ¿cómo estás?"
		case "fr":
			result.Translations[lang] = "Bonjour, comment allez-vous?"
		case "de":
			result.Translations[lang] = "Hallo, wie geht es Ihnen?"
		default:
			result.Translations[lang] = "Hello, how are you?"
		}
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
		default:
			// Perform recognition
			result := &TranslationRecognitionResult{
				ResultID:     fmt.Sprintf("ContinuousResult_%d", time.Now().UnixNano()),
				Text:         "This is a continuous recognition result.",
				Reason:       ResultReasonTranslatedSpeech,
				Offset:       time.Now().UnixNano(),
				Duration:     500 * time.Millisecond,
				Translations: make(map[string]string),
			}

			// Add a translation for each target language
			targetLangs := r.config.GetTargetLanguages()
			for _, lang := range targetLangs {
				switch lang {
				case "ja":
					result.Translations[lang] = "これは継続的な認識結果です。"
				case "es":
					result.Translations[lang] = "Este es un resultado de reconocimiento continuo."
				case "fr":
					result.Translations[lang] = "C'est un résultat de reconnaissance continue."
				case "de":
					result.Translations[lang] = "Dies ist ein kontinuierliches Erkennungsergebnis."
				default:
					result.Translations[lang] = "This is a continuous recognition result."
				}
			}

			// Signal recognizing and recognized events
			r.raiseRecognizing(result)
			r.raiseRecognized(result)

			// Simulate processing time
			time.Sleep(2 * time.Second)
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
