import { useState, useRef } from 'react';
import RecordRTC, { StereoAudioRecorder } from 'recordrtc';
import { recognizeSong, type RecognitionResult } from '../services/api';

interface AudioRecorderProps {
  onResult: (result: RecognitionResult) => void;
  onError: (error: string) => void;
}

const AudioRecorderComponent: React.FC<AudioRecorderProps> = ({ onResult, onError }) => {
  const [isRecording, setIsRecording] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const recordRTCRef = useRef<RecordRTC | null>(null);
  const recordingStartMsRef = useRef<number>(0);

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const recorder = new RecordRTC(stream, {
        type: 'audio',
        mimeType: 'audio/wav',
        recorderType: StereoAudioRecorder,
        numberOfAudioChannels: 1,
        sampleRate: 44100,
        desiredSampRate: 16000,
      });
      recordRTCRef.current = recorder;
      recordingStartMsRef.current = Date.now();
      // Removed agent log

      recorder.startRecording();
      setIsRecording(true);
    } catch (err) {
      onError('Microphone access denied or not available');
    }
  };

  const stopRecording = () => {
    if (recordRTCRef.current && isRecording) {
      recordRTCRef.current.stopRecording(async () => {
        const audioBlob = recordRTCRef.current!.getBlob();
        // Removed agent log
        
        setIsLoading(true);
        try {
          const result = await recognizeSong(audioBlob);
          onResult(result);
        } catch (err) {
          onError(
            err instanceof Error && err.message.includes('404')
              ? 'Song not recognized (No match found in database)'
              : err instanceof Error
              ? err.message
              : 'Recognition failed'
          );
        } finally {
          setIsLoading(false);
          // Stop all tracks
          if (recordRTCRef.current) {
            recordRTCRef.current.destroy();
            recordRTCRef.current = null;
          }
        }
      });
      setIsRecording(false);
    }
  };

  return (
    <div style={{ textAlign: 'center' }}>
      <div className="mic-wrapper">
        {isRecording && (
          <>
            <div className="ripple"></div>
            <div className="ripple"></div>
            <div className="ripple"></div>
          </>
        )}
        <button
          className={`mic-btn ${isRecording ? 'recording' : ''}`}
          onClick={isRecording ? stopRecording : startRecording}
          disabled={isLoading}
          title={isRecording ? 'Stop Recording' : 'Start Recording'}
        >
          {isRecording ? '⏹' : '🎤'}
        </button>
      </div>
      
      <div className="status-text">
        {isLoading && <p>Processing audio...</p>}
        {isRecording && !isLoading && (
          <p className="status-recording">Recording... Play a song snippet!</p>
        )}
      </div>
    </div>
  );
};

export default AudioRecorderComponent;