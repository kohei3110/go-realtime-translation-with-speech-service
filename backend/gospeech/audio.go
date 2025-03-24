// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license.

package gospeech

import (
	"errors"
	"io"
	"os"
)

// AudioStreamFormat represents the audio stream format
type AudioStreamFormat struct {
	samplesPerSecond int
	bitsPerSample    int
	channels         int
}

// NewAudioStreamFormat creates a new AudioStreamFormat instance
func NewAudioStreamFormat(samplesPerSecond, bitsPerSample, channels int) *AudioStreamFormat {
	return &AudioStreamFormat{
		samplesPerSecond: samplesPerSecond,
		bitsPerSample:    bitsPerSample,
		channels:         channels,
	}
}

// GetDefaultInputFormat returns the default audio input format (16kHz, 16 bits, mono)
func GetDefaultInputFormat() *AudioStreamFormat {
	return NewAudioStreamFormat(16000, 16, 1)
}

// GetWaveFormatPCM returns a PCM wave format
func GetWaveFormatPCM(samplesPerSecond, bitsPerSample, channels int) *AudioStreamFormat {
	return NewAudioStreamFormat(samplesPerSecond, bitsPerSample, channels)
}

// SamplesPerSecond returns the samples per second of the audio format
func (f *AudioStreamFormat) SamplesPerSecond() int {
	return f.samplesPerSecond
}

// BitsPerSample returns the bits per sample of the audio format
func (f *AudioStreamFormat) BitsPerSample() int {
	return f.bitsPerSample
}

// Channels returns the number of channels of the audio format
func (f *AudioStreamFormat) Channels() int {
	return f.channels
}

// AudioConfig represents audio input configuration
type AudioConfig struct {
	format     *AudioStreamFormat
	sourceType string // "Microphone", "File", "Stream", "PushStream"
	source     interface{}
}

// NewAudioConfigFromDefaultMicrophone creates an audio config from the default microphone
func NewAudioConfigFromDefaultMicrophone() (*AudioConfig, error) {
	// In a real implementation, we would initialize microphone here
	return &AudioConfig{
		format:     GetDefaultInputFormat(),
		sourceType: "Microphone",
		source:     nil,
	}, nil
}

// NewAudioConfigFromWavFile creates an audio config from a WAV file
func NewAudioConfigFromWavFile(filePath string) (*AudioConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	// In a real implementation, we would parse WAV header to get format
	return &AudioConfig{
		format:     GetDefaultInputFormat(),
		sourceType: "File",
		source:     file,
	}, nil
}

// NewAudioConfigFromStream creates an audio config from an audio stream
func NewAudioConfigFromStream(stream io.ReadCloser, format *AudioStreamFormat) (*AudioConfig, error) {
	if stream == nil {
		return nil, errors.New("stream cannot be nil")
	}

	if format == nil {
		format = GetDefaultInputFormat()
	}

	return &AudioConfig{
		format:     format,
		sourceType: "Stream",
		source:     stream,
	}, nil
}

// SourceType returns the type of audio source
func (c *AudioConfig) SourceType() string {
	return c.sourceType
}

// Format returns the audio format
func (c *AudioConfig) Format() *AudioStreamFormat {
	return c.format
}

// Source returns the audio source
func (c *AudioConfig) Source() interface{} {
	return c.source
}

// Close closes the audio source if applicable
func (c *AudioConfig) Close() error {
	if c.sourceType == "File" || c.sourceType == "Stream" {
		if closer, ok := c.source.(io.Closer); ok {
			return closer.Close()
		}
	}
	return nil
}

// PushAudioInputStream represents a stream that receives audio from the application
type PushAudioInputStream struct {
	format *AudioStreamFormat
	buffer chan []byte
	closed bool
}

// NewPushAudioInputStream creates a new push audio input stream
func NewPushAudioInputStream(format *AudioStreamFormat) *PushAudioInputStream {
	if format == nil {
		format = GetDefaultInputFormat()
	}

	return &PushAudioInputStream{
		format: format,
		buffer: make(chan []byte, 100), // Buffer 100 chunks
		closed: false,
	}
}

// Write writes audio data to the stream
func (s *PushAudioInputStream) Write(data []byte) (int, error) {
	if s.closed {
		return 0, errors.New("stream is closed")
	}

	// Make a copy of the data to avoid external mutations
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	s.buffer <- dataCopy
	return len(data), nil
}

// Read reads audio data from the stream
func (s *PushAudioInputStream) Read(p []byte) (int, error) {
	if s.closed {
		return 0, io.EOF
	}

	select {
	case data := <-s.buffer:
		n := copy(p, data)
		return n, nil
	default:
		// No data available
		return 0, nil
	}
}

// Close closes the stream
func (s *PushAudioInputStream) Close() error {
	s.closed = true
	close(s.buffer)
	return nil
}

// Format returns the audio format
func (s *PushAudioInputStream) Format() *AudioStreamFormat {
	return s.format
}

// NewAudioConfigFromPushStream creates an audio config from a push stream
func NewAudioConfigFromPushStream(stream *PushAudioInputStream) (*AudioConfig, error) {
	if stream == nil {
		return nil, errors.New("stream cannot be nil")
	}

	return &AudioConfig{
		format:     stream.Format(),
		sourceType: "PushStream",
		source:     stream,
	}, nil
}

// AudioOutputStream represents an audio output stream
type AudioOutputStream interface {
	io.WriteCloser
	Format() *AudioStreamFormat
}

// AudioOutputConfig represents audio output configuration
type AudioOutputConfig struct {
	format     *AudioStreamFormat
	outputType string // "DefaultSpeaker", "File", "Stream"
	output     interface{}
}

// NewAudioOutputConfigFromDefaultSpeaker creates an audio output config for the default speaker
func NewAudioOutputConfigFromDefaultSpeaker() *AudioOutputConfig {
	return &AudioOutputConfig{
		format:     GetDefaultInputFormat(),
		outputType: "DefaultSpeaker",
		output:     nil,
	}
}

// NewAudioOutputConfigFromFile creates an audio output config for a file
func NewAudioOutputConfigFromFile(filePath string) (*AudioOutputConfig, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return &AudioOutputConfig{
		format:     GetDefaultInputFormat(),
		outputType: "File",
		output:     file,
	}, nil
}

// NewAudioOutputConfigFromStream creates an audio output config from a stream
func NewAudioOutputConfigFromStream(stream io.WriteCloser) *AudioOutputConfig {
	return &AudioOutputConfig{
		format:     GetDefaultInputFormat(),
		outputType: "Stream",
		output:     stream,
	}
}

// OutputType returns the type of audio output
func (c *AudioOutputConfig) OutputType() string {
	return c.outputType
}

// Format returns the audio format
func (c *AudioOutputConfig) Format() *AudioStreamFormat {
	return c.format
}

// Output returns the audio output
func (c *AudioOutputConfig) Output() interface{} {
	return c.output
}

// Close closes the audio output if applicable
func (c *AudioOutputConfig) Close() error {
	if c.outputType == "File" || c.outputType == "Stream" {
		if closer, ok := c.output.(io.Closer); ok {
			return closer.Close()
		}
	}
	return nil
}
