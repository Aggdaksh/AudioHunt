import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface RecognitionResult {
  song_id: number;
  title: string;
  artist: string;
  confidence: number;
}

export const recognizeSong = async (audioBlob: Blob): Promise<RecognitionResult> => {
  const formData = new FormData();
  formData.append('audio_snippet', audioBlob, 'recording.wav');

  const response = await axios.post(
    `${API_URL}/api/recognize`,
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }
  );

  return response.data;
};