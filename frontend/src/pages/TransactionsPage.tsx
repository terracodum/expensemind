import { useRef, useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import {
  Table, TableHead, TableRow, TableCell, TableBody,
  Typography, CircularProgress, Alert, Button, Box, Stack,
} from '@mui/material'
import UploadFileIcon from '@mui/icons-material/UploadFile'

interface Transaction {
  ID: number
  Date: string
  Amount: number
  Description: string
  Category: string
}

export default function TransactionsPage() {
  const queryClient = useQueryClient()
  const inputRef = useRef<HTMLInputElement>(null)
  const [uploading, setUploading] = useState(false)
  const [uploadError, setUploadError] = useState<string | null>(null)

  const { data, isLoading, error } = useQuery<Transaction[]>({
    queryKey: ['transactions'],
    queryFn: () => fetch('/transactions').then(r => r.json()),
  })

  async function handleUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    setUploading(true)
    setUploadError(null)
    try {
      const res = await fetch('/transactions/upload', {
        method: 'POST',
        headers: { 'Content-Type': file.type },
        body: file,
      })
      if (!res.ok) {
        const body = await res.json().catch(() => ({}))
        setUploadError(body.message ?? 'Ошибка загрузки')
      } else {
        queryClient.invalidateQueries({ queryKey: ['transactions'] })
      }
    } catch {
      setUploadError('Не удалось подключиться к серверу')
    } finally {
      setUploading(false)
      if (inputRef.current) inputRef.current.value = ''
    }
  }

  return (
    <Box>
      <Stack direction="row" alignItems="center" spacing={2} sx={{ mb: 2 }}>
        <Button
          variant="outlined"
          startIcon={<UploadFileIcon />}
          onClick={() => inputRef.current?.click()}
          disabled={uploading}
        >
          {uploading ? 'Загрузка...' : 'Загрузить файл'}
        </Button>
        <input ref={inputRef} type="file" accept=".csv,.pdf" hidden onChange={handleUpload} />
        <Typography variant="caption" color="text.secondary">PDF или CSV</Typography>
      </Stack>

      {uploadError && <Alert severity="error" sx={{ mb: 2 }}>{uploadError}</Alert>}

      {isLoading && <CircularProgress />}
      {error && <Alert severity="error">Ошибка загрузки</Alert>}
      {!isLoading && !error && !data?.length && <Typography>Нет транзакций</Typography>}

      {data && data.length > 0 && (
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Дата</TableCell>
              <TableCell>Описание</TableCell>
              <TableCell>Категория</TableCell>
              <TableCell align="right">Сумма</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map(tx => (
              <TableRow key={tx.ID}>
                <TableCell>{tx.Date?.slice(0, 10)}</TableCell>
                <TableCell>{tx.Description}</TableCell>
                <TableCell>{tx.Category}</TableCell>
                <TableCell align="right">{tx.Amount.toFixed(2)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </Box>
  )
}
