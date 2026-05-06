import { useState } from 'react'
import {
  AppBar, Toolbar, Typography, Tabs, Tab, Container, Box,
} from '@mui/material'
import TransactionsPage from './pages/TransactionsPage'
import RecurringPage from './pages/RecurringPage'
import ForecastPage from './pages/ForecastPage'

export default function App() {
  const [tab, setTab] = useState(0)

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: '#f5f5f7' }}>
      <AppBar position="static" elevation={0} sx={{ bgcolor: '#fff', borderBottom: '1px solid #e0e0e0' }}>
        <Toolbar>
          <Typography variant="h6" sx={{ fontWeight: 700, color: '#1a1a1a', letterSpacing: '-0.5px' }}>
            ExpenseMind
          </Typography>
          <Tabs
            value={tab}
            onChange={(_, v) => setTab(v)}
            sx={{ ml: 4 }}
            TabIndicatorProps={{ style: { height: 3 } }}
          >
            <Tab label="Транзакции" />
            <Tab label="Регулярные" />
            <Tab label="Прогноз" />
          </Tabs>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ py: 4 }}>
        {tab === 0 && <TransactionsPage />}
        {tab === 1 && <RecurringPage />}
        {tab === 2 && <ForecastPage />}
      </Container>
    </Box>
  )
}
