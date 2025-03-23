import { useState } from 'react';
import {
  Box,
  TextField,
  Button,
  Paper,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
} from '@mui/material';
import { TranslationService } from '../../services/translationService';

const translationService = new TranslationService();

const LANGUAGES = [
  { code: 'ja', name: '日本語' },
  { code: 'en', name: '英語' },
  { code: 'zh', name: '中国語' },
  { code: 'ko', name: '韓国語' },
  { code: 'es', name: 'スペイン語' },
  { code: 'fr', name: 'フランス語' },
];

export const TextTranslation = () => {
  const [text, setText] = useState('');
  const [sourceLanguage, setSourceLanguage] = useState('ja');
  const [targetLanguage, setTargetLanguage] = useState('en');
  const [translatedText, setTranslatedText] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleTranslate = async () => {
    if (!text) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await translationService.translateText({
        text,
        sourceLanguage,
        targetLanguage,
      });
      setTranslatedText(response.translatedText);
    } catch (err) {
      setError('翻訳中にエラーが発生しました。');
      console.error('Translation error:', err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto', mt: 4, p: 2 }}>
      <Typography variant="h4" component="h1" gutterBottom sx={{ mb: 3 }}>
        テキスト翻訳
      </Typography>

      <Box sx={{ display: 'flex', gap: 2, mb: 3 }}>
        <FormControl sx={{ minWidth: 120 }}>
          <InputLabel>元の言語</InputLabel>
          <Select
            value={sourceLanguage}
            label="元の言語"
            onChange={(e) => setSourceLanguage(e.target.value)}
          >
            {LANGUAGES.map((lang) => (
              <MenuItem key={lang.code} value={lang.code}>
                {lang.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <FormControl sx={{ minWidth: 120 }}>
          <InputLabel>翻訳先</InputLabel>
          <Select
            value={targetLanguage}
            label="翻訳先"
            onChange={(e) => setTargetLanguage(e.target.value)}
          >
            {LANGUAGES.map((lang) => (
              <MenuItem key={lang.code} value={lang.code}>
                {lang.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>

      {error && (
        <Typography color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}

      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <TextField
          fullWidth
          multiline
          rows={4}
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="翻訳したいテキストを入力してください"
          variant="outlined"
        />

        <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
          <Button
            variant="contained"
            onClick={handleTranslate}
            disabled={!text || isLoading}
          >
            {isLoading ? (
              <>
                <CircularProgress size={20} sx={{ mr: 1 }} />
                翻訳中...
              </>
            ) : (
              '翻訳'
            )}
          </Button>
        </Box>
      </Paper>

      {translatedText && (
        <Paper elevation={3} sx={{ p: 3 }}>
          <Typography variant="subtitle2" color="text.secondary" gutterBottom>
            翻訳結果:
          </Typography>
          <Typography variant="body1">{translatedText}</Typography>
        </Paper>
      )}
    </Box>
  );
};