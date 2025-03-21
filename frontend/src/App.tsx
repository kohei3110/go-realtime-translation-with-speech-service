import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import { RealtimeTranslation } from './features/realtime-translation/RealtimeTranslation'

const theme = createTheme({
  palette: {
    mode: 'light',
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <RealtimeTranslation />
    </ThemeProvider>
  )
}

export default App
