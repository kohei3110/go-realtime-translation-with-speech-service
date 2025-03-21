package models_test

import (
	"testing"

	"go-realtime-translation-with-speech-service/backend/features/realtime_translation/models"

	"github.com/stretchr/testify/assert"
)

func TestTranslationRequest(t *testing.T) {
	// Test case 1: Valid translation request
	t.Run("ValidTranslationRequest", func(t *testing.T) {
		req := models.TranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			Text:           "こんにちは、元気ですか？",
		}

		err := req.Validate()
		assert.NoError(t, err, "Validation should pass for valid request")
	})

	// Test case 2: Empty source language
	t.Run("EmptySourceLanguage", func(t *testing.T) {
		req := models.TranslationRequest{
			SourceLanguage: "",
			TargetLanguage: "en",
			Text:           "こんにちは、元気ですか？",
		}

		err := req.Validate()
		assert.Error(t, err, "Validation should fail with empty source language")
		assert.Contains(t, err.Error(), "source language", "Error should mention source language")
	})

	// Test case 3: Empty target language
	t.Run("EmptyTargetLanguage", func(t *testing.T) {
		req := models.TranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "",
			Text:           "こんにちは、元気ですか？",
		}

		err := req.Validate()
		assert.Error(t, err, "Validation should fail with empty target language")
		assert.Contains(t, err.Error(), "target language", "Error should mention target language")
	})

	// Test case 4: Empty text
	t.Run("EmptyText", func(t *testing.T) {
		req := models.TranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			Text:           "",
		}

		err := req.Validate()
		assert.Error(t, err, "Validation should fail with empty text")
		assert.Contains(t, err.Error(), "text", "Error should mention text")
	})

	// Test case 5: Unsupported source language
	t.Run("UnsupportedSourceLanguage", func(t *testing.T) {
		req := models.TranslationRequest{
			SourceLanguage: "xyz",
			TargetLanguage: "en",
			Text:           "こんにちは、元気ですか？",
		}

		err := req.Validate()
		assert.Error(t, err, "Validation should fail with unsupported source language")
		assert.Contains(t, err.Error(), "source language", "Error should mention source language")
	})

	// Test case 6: Unsupported target language
	t.Run("UnsupportedTargetLanguage", func(t *testing.T) {
		req := models.TranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "xyz",
			Text:           "こんにちは、元気ですか？",
		}

		err := req.Validate()
		assert.Error(t, err, "Validation should fail with unsupported target language")
		assert.Contains(t, err.Error(), "target language", "Error should mention target language")
	})
}

func TestTranslationResponse(t *testing.T) {
	// Test basic response structure
	t.Run("ResponseFields", func(t *testing.T) {
		resp := models.TranslationResponse{
			SourceLanguage:  "ja",
			TargetLanguage:  "en",
			OriginalText:    "こんにちは、元気ですか？",
			TranslatedText:  "Hello, how are you?",
			ConfidenceScore: 0.95,
		}

		// Verify fields are set correctly
		assert.Equal(t, "ja", resp.SourceLanguage)
		assert.Equal(t, "en", resp.TargetLanguage)
		assert.Equal(t, "こんにちは、元気ですか？", resp.OriginalText)
		assert.Equal(t, "Hello, how are you?", resp.TranslatedText)
		assert.Equal(t, 0.95, resp.ConfidenceScore)
	})
}

func TestStreamingTranslationRequest(t *testing.T) {
	// Test streaming request validation
	t.Run("ValidStreamingRequest", func(t *testing.T) {
		req := models.StreamingTranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			AudioFormat:    "wav",
		}

		err := req.Validate()
		assert.NoError(t, err, "Validation should pass for valid streaming request")
	})

	// Test case: Empty source language
	t.Run("EmptySourceLanguage", func(t *testing.T) {
		req := models.StreamingTranslationRequest{
			SourceLanguage: "",
			TargetLanguage: "en",
			AudioFormat:    "wav",
		}

		err := req.Validate()
		assert.Error(t, err, "Validation should fail with empty source language")
	})

	// Test case: Empty target language
	t.Run("EmptyTargetLanguage", func(t *testing.T) {
		req := models.StreamingTranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "",
			AudioFormat:    "wav",
		}

		err := req.Validate()
		assert.Error(t, err, "Validation should fail with empty target language")
	})

	// Test case: Unsupported audio format
	t.Run("UnsupportedAudioFormat", func(t *testing.T) {
		req := models.StreamingTranslationRequest{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			AudioFormat:    "unsupported",
		}

		err := req.Validate()
		assert.Error(t, err, "Validation should fail with unsupported audio format")
	})
}

func TestStreamingTranslationResponse(t *testing.T) {
	// Test streaming response structure
	t.Run("StreamingResponseFields", func(t *testing.T) {
		resp := models.StreamingTranslationResponse{
			SourceLanguage: "ja",
			TargetLanguage: "en",
			TranslatedText: "Hello, how are you?",
			IsFinal:        true,
			SegmentID:      "segment-123",
		}

		// Verify fields are set correctly
		assert.Equal(t, "ja", resp.SourceLanguage)
		assert.Equal(t, "en", resp.TargetLanguage)
		assert.Equal(t, "Hello, how are you?", resp.TranslatedText)
		assert.True(t, resp.IsFinal)
		assert.Equal(t, "segment-123", resp.SegmentID)
	})
}
