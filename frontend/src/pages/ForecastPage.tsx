import { useState } from 'react'
import {
  Box, Button, Typography, CircularProgress, Alert,
  Paper, Stack, LinearProgress, Chip,
} from '@mui/material'
import AutoGraphIcon from '@mui/icons-material/AutoGraph'

interface ForecastJob {
  ID: number
  Status: string
  Result?: {
    PredictedBalance: number
    Confidence: number
    Points: { T: number; Balance: number }[]
  }
}

const STATUS_LABEL: Record<string, string> = {
  pending: 'В очереди...',
  running: 'Считаем...',
  done: 'Готово',
  failed: 'Ошибка',
}

export default function ForecastPage() {
  const [job, setJob] = useState<ForecastJob | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function startForecast() {
    setLoading(true)
    setError(null)
    setJob(null)
    try {
      const res = await fetch('/analytics/forecast', { method: 'POST' })
      if (!res.ok) throw new Error('Не удалось создать задачу')
      const { job_id } = await res.json()
      pollJob(job_id)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Ошибка')
      setLoading(false)
    }
  }

  async function pollJob(id: number) {
    try {
      const res = await fetch(`/analytics/forecast/${id}`)
      const data: ForecastJob = await res.json()
      setJob(data)
      if (data.Status === 'pending' || data.Status === 'running') {
        setTimeout(() => pollJob(id), 1500)
      } else {
        setLoading(false)
      }
    } catch {
      setError('Потеряна связь с сервером')
      setLoading(false)
    }
  }

  const confidence = job?.Result?.Confidence ?? 0
  const confidencePct = Math.round(confidence * 100)

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 3 }}>
        <Box>
          <Typography variant="h5" fontWeight={700}>Прогноз</Typography>
          <Typography variant="body2" color="text.secondary">
            На 30 дней вперёд на основе истории транзакций
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={loading ? <CircularProgress size={16} color="inherit" /> : <AutoGraphIcon />}
          onClick={startForecast}
          disabled={loading}
          disableElevation
        >
          {loading ? (job ? STATUS_LABEL[job.Status] : 'Запуск...') : 'Запустить прогноз'}
        </Button>
      </Stack>

      {error && <Alert severity="error" onClose={() => setError(null)}>{error}</Alert>}

      {job && job.Status === 'failed' && (
        <Alert severity="warning">
          Прогноз не удался — ML сервис недоступен или недостаточно данных
        </Alert>
      )}

      {job?.Result && (
        <Stack spacing={2}>
          <Stack direction="row" spacing={2}>
            <Paper elevation={0} sx={{ border: '1px solid #e0e0e0', borderRadius: 2, p: 3, flex: 1 }}>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Прогноз баланса через 30 дней
              </Typography>
              <Typography variant="h4" fontWeight={700}
                color={job.Result.PredictedBalance >= 0 ? 'success.main' : 'error.main'}>
                {job.Result.PredictedBalance >= 0 ? '+' : ''}
                {job.Result.PredictedBalance.toFixed(2)}
              </Typography>
            </Paper>

            <Paper elevation={0} sx={{ border: '1px solid #e0e0e0', borderRadius: 2, p: 3, flex: 1 }}>
              <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                <Typography variant="body2" color="text.secondary">Уверенность</Typography>
                <Chip
                  label={`${confidencePct}%`}
                  size="small"
                  color={confidence >= 0.7 ? 'success' : confidence >= 0.4 ? 'warning' : 'error'}
                />
              </Stack>
              <LinearProgress
                variant="determinate"
                value={confidencePct}
                color={confidence >= 0.7 ? 'success' : confidence >= 0.4 ? 'warning' : 'error'}
                sx={{ height: 8, borderRadius: 4 }}
              />
              {confidence < 0.5 && (
                <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                  Недостаточно данных для точного прогноза
                </Typography>
              )}
            </Paper>
          </Stack>
        </Stack>
      )}
    </Box>
  )
}
