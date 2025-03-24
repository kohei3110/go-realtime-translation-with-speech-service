import { useState, useCallback, useRef, useEffect } from 'react';
import RecordRTC, { RecordRTCPromisesHandler } from 'recordrtc';
import { TranslationService } from '../services/translationService';
import { StreamingTranslationResponse } from '../types/api';

const translationService = new TranslationService();

export const useTranslation = () => {
  const [isRecording, setIsRecording] = useState(false);
  const [translations, setTranslations] = useState<StreamingTranslationResponse[]>([]);
  const [error, setError] = useState<string | null>(null);

  const recorder = useRef<RecordRTCPromisesHandler | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const mediaStream = useRef<MediaStream | null>(null);

  // WebSocketの設定と管理
  const setupWebSocket = (sourceLanguage: string, targetLanguage: string) => {
    const ws = new WebSocket(`ws://localhost:8080/api/v1/streaming/ws/${Date.now()}`);

    ws.onopen = () => {
      console.log('WebSocket connected, sending initial setup message');
      // 初期設定を送信
      const setupMsg = {
        sourceLanguage,
        targetLanguage,
        audioFormat: 'audio/wav',
      };
      console.log('Sending setup message:', setupMsg);
      ws.send(JSON.stringify(setupMsg));
    };

    ws.onmessage = (event) => {
      console.log('Received WebSocket message:', event.data);
      const data = JSON.parse(event.data);
      if (data.status === 'ready') {
        console.log('WebSocket ready for streaming');
      } else if (data.translatedText) {
        console.log('Received translation result:', data);
        setTranslations(prev => {
          const newTranslation: StreamingTranslationResponse = {
            sourceLanguage: data.sourceLanguage,
            targetLanguage: data.targetLanguage,
            translatedText: data.translatedText,
            originalText: data.originalText,
            isFinal: data.isFinal,
            segmentId: data.segmentId
          };

          if (data.isFinal) {
            console.log('Adding final translation result');
            return [...prev, newTranslation];
          } else {
            console.log('Updating interim translation result');
            const updatedTranslations = [...prev];
            if (updatedTranslations.length > 0 && !updatedTranslations[updatedTranslations.length - 1].isFinal) {
              updatedTranslations[updatedTranslations.length - 1] = newTranslation;
            } else {
              updatedTranslations.push(newTranslation);
            }
            return updatedTranslations;
          }
        });
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      setError('WebSocketの接続中にエラーが発生しました');
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
      wsRef.current = null;
    };

    return ws;
  };

  const startRecording = useCallback(async (sourceLanguage: string, targetLanguage: string) => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ 
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          sampleRate: 16000,
        } 
      });
      
      mediaStream.current = stream;
      
      // WebSocket接続の確立
      wsRef.current = setupWebSocket(sourceLanguage, targetLanguage);

      // RecordRTCの設定
      recorder.current = new RecordRTCPromisesHandler(stream, {
        type: 'audio',
        mimeType: 'audio/wav',
        recorderType: RecordRTC.StereoAudioRecorder,
        timeSlice: 1000, // 1秒ごとにデータを送信
        desiredSampRate: 16000,
        numberOfAudioChannels: 1,
        ondataavailable: async (blob) => {
          if (wsRef.current?.readyState === WebSocket.OPEN) {
            const reader = new FileReader();
            reader.onloadend = () => {
              const base64Audio = (reader.result as string).split(',')[1];
              wsRef.current?.send(JSON.stringify({
                audio: {
                  data: base64Audio
                }
              }));
            };
            reader.readAsDataURL(blob);
          }
        }
      });

      await recorder.current.startRecording();
      setIsRecording(true);
      setError(null);
    } catch (err) {
      console.error('Error starting recording:', err);
      setError('マイクへのアクセスに失敗しました。');
    }
  }, []);

  const stopRecording = useCallback(async () => {
    try {
      if (!recorder.current) {
        return;
      }
      
      await recorder.current.stopRecording();
      recorder.current = null;

      // WebSocket接続を閉じる
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }

      // メディアストリームを停止
      if (mediaStream.current) {
        mediaStream.current.getTracks().forEach(track => track.stop());
        mediaStream.current = null;
      }

      setError(null);
    } catch (err) {
      console.error('Error stopping recording:', err);
      setError('録音の停止中にエラーが発生しました。');
    } finally {
      setIsRecording(false);
    }
  }, []);

  useEffect(() => {
    return () => {
      if (recorder.current) {
        recorder.current.stopRecording().catch(console.error);
        recorder.current = null;
      }
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
      if (mediaStream.current) {
        mediaStream.current.getTracks().forEach(track => track.stop());
        mediaStream.current = null;
      }
    };
  }, []);

  return {
    isRecording,
    translations,
    error,
    startRecording,
    stopRecording,
  };
};