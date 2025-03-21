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
