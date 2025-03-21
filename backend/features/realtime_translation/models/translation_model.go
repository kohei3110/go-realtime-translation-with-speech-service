package models

import (
	"errors"
	"strings"
)

// 対応言語のリスト
var SupportedLanguages = map[string]bool{
	"ja": true, // 日本語
	"en": true, // 英語
	"fr": true, // フランス語
	"es": true, // スペイン語
	"de": true, // ドイツ語
	"zh": true, // 中国語
	"ko": true, // 韓国語
}

// 対応音声フォーマットのリスト
var SupportedAudioFormats = map[string]bool{
	"wav":  true,
	"mp3":  true,
	"ogg":  true,
	"flac": true,
}

// TranslationRequest はテキスト翻訳リクエストを表す構造体
type TranslationRequest struct {
	SourceLanguage string `json:"sourceLanguage"` // 翻訳元言語コード（ISO 639-1）
	TargetLanguage string `json:"targetLanguage"` // 翻訳先言語コード（ISO 639-1）
	Text           string `json:"text"`           // 翻訳するテキスト
}

// Validate はリクエストの内容を検証するメソッド
func (r *TranslationRequest) Validate() error {
	if r.SourceLanguage == "" {
		return errors.New("source language is required")
	}

	if r.TargetLanguage == "" {
		return errors.New("target language is required")
	}

	if r.Text == "" {
		return errors.New("text is required")
	}

	if !SupportedLanguages[r.SourceLanguage] {
		return errors.New("unsupported source language")
	}

	if !SupportedLanguages[r.TargetLanguage] {
		return errors.New("unsupported target language")
	}

	return nil
}

// TranslationResponse はテキスト翻訳結果を表す構造体
type TranslationResponse struct {
	SourceLanguage  string  `json:"sourceLanguage"`  // 翻訳元言語コード
	TargetLanguage  string  `json:"targetLanguage"`  // 翻訳先言語コード
	OriginalText    string  `json:"originalText"`    // 元のテキスト
	TranslatedText  string  `json:"translatedText"`  // 翻訳されたテキスト
	ConfidenceScore float64 `json:"confidenceScore"` // 翻訳の信頼度スコア（0-1）
}

// StreamingTranslationRequest は音声ストリーミング翻訳のリクエストを表す構造体
type StreamingTranslationRequest struct {
	SourceLanguage string `json:"sourceLanguage"` // 翻訳元言語コード
	TargetLanguage string `json:"targetLanguage"` // 翻訳先言語コード
	AudioFormat    string `json:"audioFormat"`    // 音声フォーマット
}

// Validate はストリーミングリクエストの内容を検証するメソッド
func (r *StreamingTranslationRequest) Validate() error {
	if r.SourceLanguage == "" {
		return errors.New("source language is required")
	}

	if r.TargetLanguage == "" {
		return errors.New("target language is required")
	}

	if r.AudioFormat == "" {
		return errors.New("audio format is required")
	}

	if !SupportedLanguages[r.SourceLanguage] {
		return errors.New("unsupported source language")
	}

	if !SupportedLanguages[r.TargetLanguage] {
		return errors.New("unsupported target language")
	}

	// 音声フォーマットの正規化（小文字に変換）
	r.AudioFormat = strings.ToLower(r.AudioFormat)

	if !SupportedAudioFormats[r.AudioFormat] {
		return errors.New("unsupported audio format")
	}

	return nil
}

// StreamingTranslationResponse はストリーミング翻訳の結果を表す構造体
type StreamingTranslationResponse struct {
	SourceLanguage string `json:"sourceLanguage"` // 翻訳元言語コード
	TargetLanguage string `json:"targetLanguage"` // 翻訳先言語コード
	TranslatedText string `json:"translatedText"` // 翻訳されたテキスト
	IsFinal        bool   `json:"isFinal"`        // 最終結果かどうか
	SegmentID      string `json:"segmentId"`      // 音声セグメントのID
}
