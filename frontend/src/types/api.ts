export interface TranslationRequest {
  sourceLanguage: string;
  targetLanguage: string;
  text: string;
}

export interface TranslationResponse {
  sourceLanguage: string;
  targetLanguage: string;
  originalText: string;
  translatedText: string;
  confidenceScore: number;
}

export interface StreamingTranslationRequest {
  sourceLanguage: string;
  targetLanguage: string;
  audioFormat: string;
}

export interface StreamingTranslationResponse {
  sourceLanguage: string;
  targetLanguage: string;
  translatedText: string;
  isFinal: boolean;
  segmentId: string;
}

export interface AudioChunkRequest {
  sessionId: string;
  audioChunk: string; // Base64エンコードされた音声データ
}