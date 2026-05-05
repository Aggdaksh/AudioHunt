import type { RecognitionResult } from '../services/api';

interface Props {
  result: RecognitionResult | null;
}

const MatchResult: React.FC<Props> = ({ result }) => {
  if (!result) return null;

  return (
    <div style={{
      marginTop: '2rem',
      padding: '2rem',
      background: 'rgba(255, 255, 255, 0.03)',
      border: '1px solid var(--glass-border)',
      borderRadius: '16px',
      width: '100%',
      animation: 'slideUp 0.5s ease-out forwards',
      opacity: 0,
      transform: 'translateY(20px)'
    }}>
      <h2 style={{ 
        color: '#a78bfa', 
        marginBottom: '1.5rem',
        fontSize: '1.5rem',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: '0.5rem'
      }}>
        🎉 Match Found!
      </h2>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.8rem' }}>
        <p style={{ margin: 0, fontSize: '1.2rem' }}>
          <span style={{ color: 'var(--text-muted)', fontSize: '0.9rem', display: 'block', marginBottom: '0.2rem' }}>TITLE</span>
          <strong style={{ color: 'var(--text-main)', letterSpacing: '0.5px' }}>{result.title}</strong>
        </p>
        <p style={{ margin: 0, fontSize: '1.1rem' }}>
          <span style={{ color: 'var(--text-muted)', fontSize: '0.9rem', display: 'block', marginBottom: '0.2rem' }}>ARTIST</span>
          <span style={{ color: '#e2e8f0' }}>{result.artist}</span>
        </p>
        <div style={{ 
          marginTop: '1rem', 
          paddingTop: '1rem', 
          borderTop: '1px solid rgba(255,255,255,0.1)',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center'
        }}>
          <span style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>Match Score</span>
          <strong style={{ 
            color: 'var(--accent-color)', 
            background: 'rgba(59, 130, 246, 0.1)',
            padding: '4px 12px',
            borderRadius: '20px',
            fontSize: '0.9rem'
          }}>
            {result.confidence} points
          </strong>
        </div>
      </div>
      <style>{`
        @keyframes slideUp {
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
      `}</style>
    </div>
  );
};

export default MatchResult;