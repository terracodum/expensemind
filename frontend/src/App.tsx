import { useState } from 'react'
import { Box, Tab, Tabs, Typography, Container } from '@mui/material'
import TransactionsPage from './pages/TransactionsPage'
import RecurringPage from './pages/RecurringPage'
import ForecastPage from './pages/ForecastPage'

export default function App() {
  const [tab, setTab] = useState(0)

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>ExpenseMind</Typography>
      <Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ mb: 3 }}>
        <Tab label="Транзакции" />
        <Tab label="Регулярные" />
        <Tab label="Прогноз" />
      </Tabs>
      <Box>
        {tab === 0 && <TransactionsPage />}
        {tab === 1 && <RecurringPage />}
        {tab === 2 && <ForecastPage />}
      </Box>
    </Container>
  )
}
