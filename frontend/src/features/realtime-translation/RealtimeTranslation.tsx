import { Box, Button, Paper, Typography, CircularProgress } from '@mui/material';
import { useTranslation } from '../../hooks/useTranslation';
import { useCallback } from 'react';

export const RealtimeTranslation = () => {
  const { isRecording, translations, error, startRecording, stopRecording } = useTranslation();

  const handleToggleRecording = useCallback(() => {
    if (isRecording) {
      stopRecording();
    } else {
      startRecording();
    }
  }, [isRecording, startRecording, stopRecording]);

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto', mt: 4, p: 2 }}>
      <Typography variant="h4" component="h1" gutterBottom sx={{ mb: 3 }}>
        リアルタイム翻訳
      </Typography>

      {error && (
        <Typography color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}

      <Box sx={{ mb: 3, display: 'flex', alignItems: 'center', gap: 2 }}>
        <Button
          variant="contained"
          color={isRecording ? 'error' : 'primary'}
          onClick={handleToggleRecording}
          startIcon={isRecording ? null : <span>🎤</span>}
          disabled={!!error}
        >
          {isRecording ? '停止' : '録音開始'}
        </Button>
        {isRecording && (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <CircularProgress size={20} />
            <Typography variant="body2" color="text.secondary">
              録音中...
            </Typography>
          </Box>
        )}
      </Box>

      <Paper
        elevation={3}
        sx={{
          p: 3,
          minHeight: 200,
          maxHeight: 400,
          overflowY: 'auto',
          bgcolor: 'background.paper',
          position: 'relative',
        }}
      >
        {translations.length > 0 ? (
          translations.map((translation, index) => (
            <Box
              key={`${translation.segmentId}-${index}`}
              sx={{
                mb: 2,
                transition: 'opacity 0.3s ease-in-out',
              }}
            >
              <Typography
                variant="body1"
                sx={{
                  opacity: translation.isFinal ? 1 : 0.7,
                  fontStyle: translation.isFinal ? 'normal' : 'italic',
                  color: translation.isFinal ? 'text.primary' : 'text.secondary',
                }}
              >
                {translation.translatedText}
              </Typography>
            </Box>
          ))
        ) : (
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              height: '100%',
              minHeight: 160,
              color: 'text.secondary',
            }}
          >
            <Typography variant="body1" sx={{ mb: 1 }}>
              🎤 録音を開始すると、ここに翻訳結果が表示されます
            </Typography>
            <Typography variant="body2" color="text.secondary">
              マイクへのアクセスを許可してください
            </Typography>
          </Box>
        )}
      </Paper>
    </Box>
  );
};