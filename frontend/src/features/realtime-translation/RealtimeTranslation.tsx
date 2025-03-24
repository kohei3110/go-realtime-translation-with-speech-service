import React, { useState, useEffect, useRef } from 'react';
import './RealtimeTranslation.css';

// リアルタイム音声認識・翻訳コンポーネント
export const RealtimeTranslation: React.FC = () => {
  // 状態管理
  const [isActive, setIsActive] = useState<boolean>(false);
  const [sourceLanguage, setSourceLanguage] = useState<string>('en-US');
  const [targetLanguage, setTargetLanguage] = useState<string>('ja');
  const [sessionId, setSessionId] = useState<string>('');
  const [translations, setTranslations] = useState<Array<{ text: string; isFinal: boolean; original: string }>>([]);
  const [status, setStatus] = useState<string>('待機中');
  const [error, setError] = useState<string>('');

  // WebSocketの参照
  const wsRef = useRef<WebSocket | null>(null);

  // 言語オプション
  const languageOptions = [
    { value: 'en-US', label: '英語' },
    { value: 'ja-JP', label: '日本語' },
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

  // WebSocketのセットアップ
  const setupWebSocket = (url: string) => {
    const ws = new WebSocket(`ws://localhost:8080${url}`);

    ws.onopen = () => {
      console.log('WebSocket接続が開きました');
      setStatus('接続しました - 初期化中...');

      // 初期設定を送信
      ws.send(JSON.stringify({
        sourceLanguage: sourceLanguage,
        targetLanguage: targetLanguage,
        audioFormat: 'audio/wav',
      }));
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.status === 'ready') {
        setStatus('認識中 - お話しください...');
      } else if (data.translatedText) {
        setTranslations(prev => {
          // 途中結果の場合は最後のアイテムを更新、確定結果の場合は新しいアイテムを追加
          if (data.isFinal) {
            return [...prev.filter(t => t.isFinal), { 
              text: data.translatedText, 
              isFinal: true,
              original: data.originalText || ''
            }];
          } else {
            const updatedTranslations = [...prev];
            // 最後のアイテムが途中結果なら更新、そうでなければ新しい途中結果を追加
            if (updatedTranslations.length > 0 && !updatedTranslations[updatedTranslations.length - 1].isFinal) {
              updatedTranslations[updatedTranslations.length - 1] = { 
                text: data.translatedText, 
                isFinal: false,
                original: data.originalText || ''
              };
            } else {
              updatedTranslations.push({ 
                text: data.translatedText, 
                isFinal: false,
                original: data.originalText || ''
              });
            }
            return updatedTranslations;
          }
        });
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocketエラー:', error);
      setError('WebSocket接続でエラーが発生しました');
      setStatus('エラー');
    };

    ws.onclose = () => {
      console.log('WebSocket接続が閉じられました');
      setStatus('接続終了');
      wsRef.current = null;
      setIsActive(false);
    };

    wsRef.current = ws;
  };

  // 翻訳セッションの開始
  const startTranslation = async () => {
    try {
      setError('');
      setStatus('セッション開始中...');

      const response = await fetch('http://localhost:8080/api/v1/streaming/start', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          sourceLanguage,
          targetLanguage,
          audioFormat: 'audio/wav',
        }),
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      const data = await response.json();
      setSessionId(data.sessionId);
      setIsActive(true);
      setTranslations([]);

      // WebSocketに接続
      setupWebSocket(data.webSocketURL);
    } catch (err: any) {
      console.error('Error starting translation:', err);
      setError(`翻訳セッションの開始に失敗しました: ${err.message}`);
      setStatus('エラー');
    }
  };

  // 翻訳セッションの停止
  const stopTranslation = async () => {
    try {
      setStatus('セッション終了中...');

      // WebSocket接続を閉じる
      if (wsRef.current) {
        wsRef.current.close();
      }

      // APIを呼び出してセッションをクリーンアップ
      if (sessionId) {
        await fetch('http://localhost:8080/api/v1/streaming/close', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            sessionId,
          }),
        });
      }

      setIsActive(false);
      setStatus('待機中');
    } catch (err: any) {
      console.error('Error stopping translation:', err);
      setError(`翻訳セッションの停止に失敗しました: ${err.message}`);
    }
  };

  // コンポーネントのクリーンアップ
  useEffect(() => {
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

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
              disabled={isActive}
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
              disabled={isActive}
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
          {!isActive ? (
            <button onClick={startTranslation} className="start-button">
              翻訳開始
            </button>
          ) : (
            <button onClick={stopTranslation} className="stop-button">
              翻訳停止
            </button>
          )}
        </div>
      </div>
      
      <div className="status-bar">
        <span className={`status-indicator ${isActive ? 'active' : ''}`}></span>
        <span className="status-text">{status}</span>
      </div>
      
      {error && <div className="error-message">{error}</div>}
      
      <div className="translations-container">
        {translations.length > 0 ? (
          <div className="translations-list">
            {translations.map((translation, index) => (
              <div
                key={index}
                className={`translation-item ${translation.isFinal ? 'final' : 'interim'}`}
              >
                <div className="translation-text">{translation.text}</div>
                {translation.original && (
                  <div className="original-text">{translation.original}</div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <div className="no-translations">
            {isActive ? 'お話しください...' : '翻訳を開始するには「翻訳開始」ボタンをクリックしてください'}
          </div>
        )}
      </div>
    </div>
  );
};

export default RealtimeTranslation;