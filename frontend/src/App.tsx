import { useState } from 'react';
import AudioRecorderComponent from './components/AudioRecorder';
import MatchResult from './components/MatchResults';
import type { RecognitionResult } from './services/api';
import './App.css';

function App() {
  const [result, setResult] = useState<RecognitionResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const testBackend = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/recognize', {
        method: 'OPTIONS'
      });
      alert(`OPTIONS status: ${res.status}`);
    } catch (err) {
      alert('Connection failed: ' + err);
    }
  };

  return (
    <div className="App">
      <div className="glass-container">
        <h1>🎵 AudioHunt</h1>
        <p className="subtitle">Tap the microphone and play a song</p>

        <button onClick={testBackend} className="test-btn">
          🧪 Test Backend Connection
        </button>

        <AudioRecorderComponent
          onResult={(res) => {
            setResult(res);
            setError(null);
          }}
          onError={(err: string) => {
            setError(err);
            setResult(null);
          }}
        />

        {error && <p style={{ color: 'var(--danger-color)', marginTop: '1rem' }}>Error: {error}</p>}
        <MatchResult result={result} />
      </div>
    </div>
  );
}

export default App;