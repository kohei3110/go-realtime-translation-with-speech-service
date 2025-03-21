import axios from 'axios';
import {
  TranslationRequest,
  TranslationResponse,
  StreamingTranslationRequest,
  StreamingTranslationResponse,
  AudioChunkRequest,
} from '../types/api';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export class TranslationService {
  async translateText(request: TranslationRequest): Promise<TranslationResponse> {
    const response = await axios.post<TranslationResponse>(
      `${API_BASE_URL}/translate`,
      request
    );
    return response.data;
  }

  async startStreamingSession(
    request: StreamingTranslationRequest
  ): Promise<string> {
    const response = await axios.post<{ sessionId: string }>(
      `${API_BASE_URL}/streaming/start`,
      request
    );
    return response.data.sessionId;
  }

  async processAudioChunk(
    request: AudioChunkRequest
  ): Promise<StreamingTranslationResponse[]> {
    const response = await axios.post<StreamingTranslationResponse[]>(
      `${API_BASE_URL}/streaming/process`,
      request
    );
    return response.data;
  }

  async closeStreamingSession(sessionId: string): Promise<void> {
    await axios.post(`${API_BASE_URL}/streaming/close`, { sessionId });
  }
}