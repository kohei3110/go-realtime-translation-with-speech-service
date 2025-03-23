import './App.css'
import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import { TranslationApp } from './features/realtime-translation/TranslationApp'

const theme = createTheme({
  palette: {
    mode: 'light',
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <TranslationApp />
    </ThemeProvider>
  )
}

export default App
