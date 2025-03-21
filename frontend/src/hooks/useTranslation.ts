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
  const sessionId = useRef<string | null>(null);
  const mediaStream = useRef<MediaStream | null>(null);

  const processAudioData = async (audioData: Blob) => {
    if (!sessionId.current) return;

    const reader = new FileReader();
    reader.onloadend = async () => {
      const base64Audio = (reader.result as string).split(',')[1];
      try {
        const responses = await translationService.processAudioChunk({
          sessionId: sessionId.current!,
          audioChunk: base64Audio,
        });
        if (responses.length > 0) {
          setTranslations(prev => {
            // 同じsegmentIdの場合は更新、新しいsegmentIdの場合は追加
            const newTranslations = [...prev];
            responses.forEach(response => {
              const index = newTranslations.findIndex(t => t.segmentId === response.segmentId);
              if (index !== -1) {
                newTranslations[index] = response;
              } else {
                newTranslations.push(response);
              }
            });
            return newTranslations;
          });
        }
      } catch (err) {
        console.error('音声チャンクの処理中にエラーが発生しました:', err);
      }
    };
    reader.readAsDataURL(audioData);
  };

  const startRecording = useCallback(async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ 
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          sampleRate: 16000,
        } 
      });
      mediaStream.current = stream;

      // ストリーミングセッションを開始
      sessionId.current = await translationService.startStreamingSession({
        sourceLanguage: 'ja',
        targetLanguage: 'en',
        audioFormat: 'wav',
      });

      recorder.current = new RecordRTCPromisesHandler(stream, {
        type: 'audio',
        mimeType: 'audio/wav',
        recorderType: RecordRTC.StereoAudioRecorder,
        timeSlice: 250, // 250msごとにデータを取得
        desiredSampRate: 16000,
        ondataavailable: (blob: Blob) => {
          processAudioData(blob);
        },
      });

      await recorder.current.startRecording();
      setIsRecording(true);
    } catch (err) {
      setError('マイクの使用許可が必要です。');
      console.error('録音の開始に失敗しました:', err);
    }
  }, []);

  const stopRecording = useCallback(async () => {
    if (recorder.current) {
      await recorder.current.stopRecording();
      if (sessionId.current) {
        await translationService.closeStreamingSession(sessionId.current);
        sessionId.current = null;
      }
      if (mediaStream.current) {
        mediaStream.current.getTracks().forEach(track => track.stop());
        mediaStream.current = null;
      }
      setIsRecording(false);
    }
  }, []);

  useEffect(() => {
    return () => {
      if (recorder.current) {
        recorder.current.stopRecording();
      }
      if (sessionId.current) {
        translationService.closeStreamingSession(sessionId.current);
      }
      if (mediaStream.current) {
        mediaStream.current.getTracks().forEach(track => track.stop());
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