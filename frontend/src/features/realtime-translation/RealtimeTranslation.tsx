import React, { useState } from 'react';
import './RealtimeTranslation.css';
import { useTranslation } from '../../hooks/useTranslation';

export const RealtimeTranslation: React.FC = () => {
  const [sourceLanguage, setSourceLanguage] = useState<string>('ja-JP');
  const [targetLanguage, setTargetLanguage] = useState<string>('en');
  const { isRecording, translations, error, startRecording, stopRecording } = useTranslation();

  // 言語オプション
  const languageOptions = [
    { value: 'ja-JP', label: '日本語' },
    { value: 'en-US', label: '英語' },
    { value: 'es-ES', label: 'スペイン語' },
    { value: 'fr-FR', label: 'フランス語' },
    { value: 'de-DE', label: 'ドイツ語' },
    { value: 'zh-CN', label: '中国語 (簡体字)' },
  ];

  const targetLanguageOptions = [
    { value: 'en', label: '英語' },
    { value: 'ja', label: '日本語' },
    { value: 'es', label: 'スペイン語' },
    { value: 'fr', label: 'フランス語' },
    { value: 'de', label: 'ドイツ語' },
    { value: 'zh-Hans', label: '中国語 (簡体字)' },
  ];

  const handleStartTranslation = async () => {
    await startRecording(sourceLanguage, targetLanguage);
  };

  const handleStopTranslation = async () => {
    await stopRecording();
  };

  return (
    <div className="realtime-translation-container">
      <h2>リアルタイム音声翻訳</h2>
      
      <div className="controls">
        <div className="language-selectors">
          <div className="language-selector">
            <label>音声言語:</label>
            <select
              value={sourceLanguage}
              onChange={(e) => setSourceLanguage(e.target.value)}
              disabled={isRecording}
            >
              {languageOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
          
          <div className="language-selector">
            <label>翻訳言語:</label>
            <select
              value={targetLanguage}
              onChange={(e) => setTargetLanguage(e.target.value)}
              disabled={isRecording}
            >
              {targetLanguageOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
        </div>
        
        <div className="action-buttons">
          {!isRecording ? (
            <button onClick={handleStartTranslation} className="start-button">
              翻訳開始
            </button>
          ) : (
            <button onClick={handleStopTranslation} className="stop-button">
              翻訳停止
            </button>
          )}
        </div>
      </div>
      
      <div className="status-bar">
        <span className={`status-indicator ${isRecording ? 'active' : ''}`}></span>
        <span className="status-text">{isRecording ? '認識中' : '待機中'}</span>
      </div>
      
      {error && <div className="error-message">{error}</div>}
      
      <div className="translations-container">
        {translations.length > 0 ? (
          <div className="translations-list">
            {translations.map((translation, index) => (
              <div
                key={`${translation.segmentId}-${index}`}
                className={`translation-item ${translation.isFinal ? 'final' : 'interim'}`}
              >
                <div className="translation-text">{translation.translatedText}</div>
                {translation.originalText && (
                  <div className="original-text">{translation.originalText}</div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <div className="no-translations">
            {isRecording ? 'お話しください...' : '翻訳を開始するには「翻訳開始」ボタンをクリックしてください'}
          </div>
        )}
      </div>
    </div>
  );
};

export default RealtimeTranslation;