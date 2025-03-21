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
    console.log('Sending text translation request:', request);
    const response = await axios.post<TranslationResponse>(
      `${API_BASE_URL}/translate`,
      request
    );
    console.log('Received text translation response:', response.data);
    return response.data;
  }

  async startStreamingSession(
    request: StreamingTranslationRequest
  ): Promise<string> {
    console.log('Starting streaming session with request:', request);
    const response = await axios.post<{ sessionId: string }>(
      `${API_BASE_URL}/streaming/start`,
      request
    );
    console.log('Streaming session started with ID:', response.data.sessionId);
    return response.data.sessionId;
  }

  async processAudioChunk(
    request: AudioChunkRequest
  ): Promise<StreamingTranslationResponse[]> {
    console.log('Processing audio chunk for session:', request.sessionId);
    const response = await axios.post<StreamingTranslationResponse[]>(
      `${API_BASE_URL}/streaming/process`,
      request
    );
    if (response.status !== 200) {
      console.error('Failed to process audio chunk:', response);
      throw new Error('Failed to process audio chunk');
    }
    if (!Array.isArray(response.data)) {
      console.error('Invalid response format:', response.data);
      throw new Error('Invalid response format');
    }
    if (response.data.length === 0) {
      console.warn('No translation response received');
      throw new Error('No translation response received');
    }
    console.log('Received translation responses:', response.data);
    return response.data;
  }

  async closeStreamingSession(sessionId: string): Promise<void> {
    console.log('Closing streaming session:', sessionId);
    await axios.post(`${API_BASE_URL}/streaming/close`, { sessionId });
    console.log('Streaming session closed successfully:', sessionId);
  }
}