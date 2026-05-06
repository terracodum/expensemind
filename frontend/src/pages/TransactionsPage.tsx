import { useRef, useState } from 'react'
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query'
import {
  Box, Button, Paper, Stack, Table, TableBody, TableCell,
  TableContainer, TableHead, TableRow, Typography,
  CircularProgress, Alert, IconButton, Select, MenuItem,
  Chip, Skeleton,
} from '@mui/material'
import UploadFileIcon from '@mui/icons-material/UploadFile'
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline'

interface Transaction {
  ID: number
  Date: string
  Amount: number
  Description: string
  Category: string
}

const CATEGORIES = ['food', 'transport', 'entertainment', 'transfer', 'other']

const CATEGORY_COLOR: Record<string, 'success' | 'info' | 'secondary' | 'default' | 'warning'> = {
  food: 'success',
  transport: 'info',
  entertainment: 'secondary',
  transfer: 'default',
  other: 'warning',
}

function CategoryCell({ tx, onUpdate }: { tx: Transaction; onUpdate: (tx: Transaction) => void }) {
  const [editing, setEditing] = useState(false)

  if (editing) {
    return (
      <Select
        size="small"
        value={tx.Category || 'other'}
        autoFocus
        onBlur={() => setEditing(false)}
        onChange={e => {
          onUpdate({ ...tx, Category: e.target.value })
          setEditing(false)
        }}
        sx={{ minWidth: 130 }}
      >
        {CATEGORIES.map(c => <MenuItem key={c} value={c}>{c}</MenuItem>)}
      </Select>
    )
  }

  return (
    <Chip
      label={tx.Category || '—'}
      size="small"
      color={CATEGORY_COLOR[tx.Category] ?? 'default'}
      onClick={() => setEditing(true)}
      sx={{ cursor: 'pointer', fontWeight: 500 }}
    />
  )
}

export default function TransactionsPage() {
  const qc = useQueryClient()
  const inputRef = useRef<HTMLInputElement>(null)
  const [uploadError, setUploadError] = useState<string | null>(null)

  const { data, isLoading, error } = useQuery<Transaction[]>({
    queryKey: ['transactions'],
    queryFn: () => fetch('/transactions').then(r => r.json()),
  })

  const uploadMutation = useMutation({
    mutationFn: (file: File) =>
      fetch('/transactions/upload', {
        method: 'POST',
        headers: { 'Content-Type': file.type },
        body: file,
      }).then(async r => {
        if (!r.ok) throw new Error((await r.json().catch(() => ({}))).message ?? 'Ошибка загрузки')
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['transactions'] })
      setUploadError(null)
    },
    onError: (e: Error) => setUploadError(e.message),
    onSettled: () => { if (inputRef.current) inputRef.current.value = '' },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => fetch(`/transactions/${id}`, { method: 'DELETE' }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['transactions'] }),
  })

  const updateMutation = useMutation({
    mutationFn: (tx: Transaction) =>
      fetch(`/transactions/${tx.ID}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(tx),
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['transactions'] }),
  })

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 3 }}>
        <Box>
          <Typography variant="h5" fontWeight={700}>Транзакции</Typography>
          {data && (
            <Typography variant="body2" color="text.secondary">
              {data.length} записей
            </Typography>
          )}
        </Box>
        <Stack direction="row" spacing={1} alignItems="center">
          {uploadMutation.isPending && <CircularProgress size={20} />}
          <Button
            variant="contained"
            startIcon={<UploadFileIcon />}
            onClick={() => inputRef.current?.click()}
            disabled={uploadMutation.isPending}
            disableElevation
          >
            Загрузить файл
          </Button>
          <input ref={inputRef} type="file" accept=".csv,.pdf" hidden
            onChange={e => { const f = e.target.files?.[0]; if (f) uploadMutation.mutate(f) }} />
        </Stack>
      </Stack>

      {uploadError && <Alert severity="error" sx={{ mb: 2 }} onClose={() => setUploadError(null)}>{uploadError}</Alert>}
      {error && <Alert severity="error" sx={{ mb: 2 }}>Не удалось загрузить транзакции</Alert>}

      <Paper elevation={0} sx={{ border: '1px solid #e0e0e0', borderRadius: 2 }}>
        <TableContainer>
          <Table stickyHeader size="small">
            <TableHead>
              <TableRow sx={{ '& th': { bgcolor: '#fafafa', fontWeight: 600 } }}>
                <TableCell>Дата</TableCell>
                <TableCell>Описание</TableCell>
                <TableCell>Категория</TableCell>
                <TableCell align="right">Сумма</TableCell>
                <TableCell width={48} />
              </TableRow>
            </TableHead>
            <TableBody>
              {isLoading && Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={i}>
                  {Array.from({ length: 4 }).map((_, j) => (
                    <TableCell key={j}><Skeleton /></TableCell>
                  ))}
                  <TableCell />
                </TableRow>
              ))}
              {!isLoading && !data?.length && (
                <TableRow>
                  <TableCell colSpan={5} align="center" sx={{ py: 6, color: 'text.secondary' }}>
                    Загрузите PDF или CSV выписку
                  </TableCell>
                </TableRow>
              )}
              {data?.map(tx => (
                <TableRow key={tx.ID} hover>
                  <TableCell sx={{ color: 'text.secondary', whiteSpace: 'nowrap' }}>
                    {tx.Date?.slice(0, 10)}
                  </TableCell>
                  <TableCell sx={{ maxWidth: 320, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {tx.Description}
                  </TableCell>
                  <TableCell>
                    <CategoryCell tx={tx} onUpdate={t => updateMutation.mutate(t)} />
                  </TableCell>
                  <TableCell align="right" sx={{ fontVariantNumeric: 'tabular-nums', fontWeight: 500,
                    color: tx.Amount >= 0 ? 'success.main' : 'text.primary' }}>
                    {tx.Amount >= 0 ? '+' : ''}{tx.Amount.toFixed(2)}
                  </TableCell>
                  <TableCell>
                    <IconButton size="small" onClick={() => deleteMutation.mutate(tx.ID)}
                      sx={{ color: 'text.disabled', '&:hover': { color: 'error.main' } }}>
                      <DeleteOutlineIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    </Box>
  )
}
