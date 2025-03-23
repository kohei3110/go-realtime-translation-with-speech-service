import { useState } from 'react';
import { Box, Tabs, Tab } from '@mui/material';
import { TextTranslation } from './TextTranslation';
import { RealtimeTranslation } from './RealtimeTranslation';

export const TranslationApp = () => {
  const [currentTab, setCurrentTab] = useState(0);

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };

  return (
    <Box>
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Tabs value={currentTab} onChange={handleTabChange} centered>
          <Tab label="テキスト翻訳" />
          <Tab label="音声翻訳" />
        </Tabs>
      </Box>

      {currentTab === 0 ? <TextTranslation /> : <RealtimeTranslation />}
    </Box>
  );
};