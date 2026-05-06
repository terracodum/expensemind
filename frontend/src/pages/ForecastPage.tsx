import { useState } from 'react'
import { Button, Typography, CircularProgress, Alert, Box } from '@mui/material'

interface ForecastJob {
  ID: number
  Status: string
  Result?: {
    PredictedBalance: number
    Confidence: number
  }
}

export default function ForecastPage() {
  const [job, setJob] = useState<ForecastJob | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function startForecast() {
    setLoading(true)
    setError(null)
    try {
      const res = await fetch('/analytics/forecast', { method: 'POST' })
      const { job_id } = await res.json()
      pollJob(job_id)
    } catch {
      setError('Не удалось запустить прогноз')
      setLoading(false)
    }
  }

  async function pollJob(id: number) {
    const res = await fetch(`/analytics/forecast/${id}`)
    const data: ForecastJob = await res.json()
    setJob(data)
    if (data.Status === 'pending' || data.Status === 'running') {
      setTimeout(() => pollJob(id), 1500)
    } else {
      setLoading(false)
    }
  }

  return (
    <Box>
      <Button variant="contained" onClick={startForecast} disabled={loading}>
        Запустить прогноз
      </Button>

      {loading && <CircularProgress sx={{ ml: 2 }} size={24} />}
      {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}

      {job && job.Status === 'done' && job.Result && (
        <Box sx={{ mt: 3 }}>
          <Typography>Прогнозируемый баланс: <b>{job.Result.PredictedBalance.toFixed(2)}</b></Typography>
          <Typography>Уверенность: <b>{(job.Result.Confidence * 100).toFixed(0)}%</b></Typography>
        </Box>
      )}

      {job && job.Status === 'failed' && (
        <Alert severity="warning" sx={{ mt: 2 }}>Прогноз не удался (ML недоступен)</Alert>
      )}
    </Box>
  )
}
