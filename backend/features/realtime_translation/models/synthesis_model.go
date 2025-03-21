package models

import (
	"errors"
)

// SynthesisRequest はテキストを音声に合成するためのリクエスト
type SynthesisRequest struct {
	Language string `json:"language"` // 言語コード（例: "ja-JP", "en-US"）
	Text     string `json:"text"`     // 合成するテキスト
}

// Validate はリクエストの内容を検証するメソッド
func (r *SynthesisRequest) Validate() error {
	if r.Language == "" {
		return errors.New("language is required")
	}
	if r.Text == "" {
		return errors.New("text is required")
	}
	return nil
}

// SynthesisResponse はテキストから音声への合成結果のレスポンス
type SynthesisResponse struct {
	Language  string `json:"language"`            // 言語コード
	Text      string `json:"text"`                // 合成されたテキスト
	AudioData []byte `json:"audioData,omitempty"` // 音声データ（バイナリ）
}
