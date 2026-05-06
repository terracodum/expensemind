import { useRef, useState } from 'react'
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query'
import {
  Box, Button, Paper, Stack, Table, TableBody, TableCell,
  TableContainer, TableHead, TableRow, Typography,
  CircularProgress, Alert, IconButton, Select, MenuItem,
  Chip, Skeleton, Dialog, DialogTitle, DialogContent,
  DialogActions, TextField, Accordion, AccordionSummary, AccordionDetails,
} from '@mui/material'
import UploadFileIcon from '@mui/icons-material/UploadFile'
import AddIcon from '@mui/icons-material/Add'
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline'
import ExpandMoreIcon from '@mui/icons-material/ExpandMore'

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

const MONTH_NAMES = ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь',
  'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь']

function monthKey(date: string) {
  return date.slice(0, 7)
}

function monthLabel(key: string) {
  const [year, month] = key.split('-')
  return `${MONTH_NAMES[parseInt(month) - 1]} ${year}`
}

function groupByMonth(txs: Transaction[]) {
  const map = new Map<string, Transaction[]>()
  for (const tx of txs) {
    const key = monthKey(tx.Date)
    if (!map.has(key)) map.set(key, [])
    map.get(key)!.push(tx)
  }
  return Array.from(map.entries()).sort((a, b) => b[0].localeCompare(a[0]))
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

const EMPTY_FORM = { Date: '', Amount: '', Description: '', Category: 'other' }

function AddTransactionDialog({ open, onClose, onSave }: {
  open: boolean
  onClose: () => void
  onSave: (tx: Omit<Transaction, 'ID'>) => void
}) {
  const [form, setForm] = useState(EMPTY_FORM)

  function handleSave() {
    onSave({ ...form, Amount: parseFloat(form.Amount) })
    setForm(EMPTY_FORM)
  }

  const valid = form.Date && form.Amount && !isNaN(parseFloat(form.Amount)) && form.Description

  return (
    <Dialog open={open} onClose={onClose} maxWidth="xs" fullWidth>
      <DialogTitle>Добавить транзакцию</DialogTitle>
      <DialogContent>
        <Stack spacing={2} sx={{ mt: 1 }}>
          <TextField label="Дата" type="date" size="small" InputLabelProps={{ shrink: true }}
            value={form.Date} onChange={e => setForm(f => ({ ...f, Date: e.target.value }))} />
          <TextField label="Описание" size="small"
            value={form.Description} onChange={e => setForm(f => ({ ...f, Description: e.target.value }))} />
          <TextField label="Сумма" type="number" size="small" helperText="Положительная — доход, отрицательная — расход"
            value={form.Amount} onChange={e => setForm(f => ({ ...f, Amount: e.target.value }))} />
          <Select size="small" value={form.Category} onChange={e => setForm(f => ({ ...f, Category: e.target.value }))}>
            {CATEGORIES.map(c => <MenuItem key={c} value={c}>{c}</MenuItem>)}
          </Select>
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Отмена</Button>
        <Button variant="contained" disableElevation disabled={!valid} onClick={handleSave}>Добавить</Button>
      </DialogActions>
    </Dialog>
  )
}

function MonthGroup({ monthKey, txs, defaultExpanded, onDelete, onUpdate }: {
  monthKey: string
  txs: Transaction[]
  defaultExpanded: boolean
  onDelete: (id: number) => void
  onUpdate: (tx: Transaction) => void
}) {
  return (
    <Accordion defaultExpanded={defaultExpanded} elevation={0}
      sx={{ border: '1px solid #e0e0e0', borderRadius: '8px !important', mb: 1, '&:before': { display: 'none' } }}>
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        <Stack direction="row" spacing={2} alignItems="center">
          <Typography fontWeight={600}>{monthLabel(monthKey)}</Typography>
          <Typography variant="body2" color="text.secondary">{txs.length} транзакций</Typography>
        </Stack>
      </AccordionSummary>
      <AccordionDetails sx={{ p: 0 }}>
        <TableContainer>
          <Table size="small">
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
              {txs.map(tx => (
                <TableRow key={tx.ID} hover>
                  <TableCell sx={{ color: 'text.secondary', whiteSpace: 'nowrap' }}>
                    {tx.Date?.slice(0, 10)}
                  </TableCell>
                  <TableCell sx={{ maxWidth: 320, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {tx.Description}
                  </TableCell>
                  <TableCell>
                    <CategoryCell tx={tx} onUpdate={onUpdate} />
                  </TableCell>
                  <TableCell align="right" sx={{ fontVariantNumeric: 'tabular-nums', fontWeight: 500,
                    color: tx.Amount >= 0 ? 'success.main' : 'text.primary' }}>
                    {tx.Amount >= 0 ? '+' : ''}{tx.Amount.toFixed(2)}
                  </TableCell>
                  <TableCell>
                    <IconButton size="small" onClick={() => onDelete(tx.ID)}
                      sx={{ color: 'text.disabled', '&:hover': { color: 'error.main' } }}>
                      <DeleteOutlineIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </AccordionDetails>
    </Accordion>
  )
}

const BANKS = [
  { value: 'tbank', label: 'Т-Банк' },
  { value: 'vtb', label: 'ВТБ' },
  { value: 'sber', label: 'Сбер' },
  { value: 'ozon', label: 'Озон' },
]

function BankSelectDialog({ file, onClose, onConfirm }: {
  file: File
  onClose: () => void
  onConfirm: (file: File, bank: string) => void
}) {
  const [bank, setBank] = useState('tbank')
  return (
    <Dialog open onClose={onClose} maxWidth="xs" fullWidth>
      <DialogTitle>Выберите банк</DialogTitle>
      <DialogContent>
        <Stack spacing={2} sx={{ mt: 1 }}>
          <Typography variant="body2" color="text.secondary">{file.name}</Typography>
          <Select size="small" value={bank} onChange={e => setBank(e.target.value)}>
            {BANKS.map(b => <MenuItem key={b.value} value={b.value}>{b.label}</MenuItem>)}
          </Select>
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Отмена</Button>
        <Button variant="contained" disableElevation onClick={() => onConfirm(file, bank)}>Загрузить</Button>
      </DialogActions>
    </Dialog>
  )
}

export default function TransactionsPage() {
  const qc = useQueryClient()
  const inputRef = useRef<HTMLInputElement>(null)
  const [uploadError, setUploadError] = useState<string | null>(null)
  const [addOpen, setAddOpen] = useState(false)
  const [pendingFile, setPendingFile] = useState<File | null>(null)

  const { data, isLoading, error } = useQuery<Transaction[]>({
    queryKey: ['transactions'],
    queryFn: () => fetch('/transactions').then(r => r.json()),
  })

  const uploadMutation = useMutation({
    mutationFn: ({ file, bank }: { file: File; bank: string }) =>
      fetch(`/transactions/upload?bank=${bank}`, {
        method: 'POST',
        headers: { 'Content-Type': file.type },
        body: file,
      }).then(async r => {
        if (!r.ok) throw new Error((await r.json().catch(() => ({}))).message ?? 'Ошибка загрузки')
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['transactions'] })
      setUploadError(null)
      setPendingFile(null)
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

  const createMutation = useMutation({
    mutationFn: (tx: Omit<Transaction, 'ID'>) =>
      fetch('/transactions', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(tx),
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['transactions'] })
      setAddOpen(false)
    },
  })

  const groups = data ? groupByMonth(data) : []
  const latestKey = groups[0]?.[0]

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
          <Button variant="outlined" startIcon={<AddIcon />} onClick={() => setAddOpen(true)} disableElevation>
            Добавить
          </Button>
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
            onChange={e => {
              const f = e.target.files?.[0]
              if (!f) return
              if (f.type === 'application/pdf') {
                setPendingFile(f)
              } else {
                uploadMutation.mutate({ file: f, bank: '' })
              }
            }} />
        </Stack>
      </Stack>

      <AddTransactionDialog
        open={addOpen}
        onClose={() => setAddOpen(false)}
        onSave={tx => createMutation.mutate(tx)}
      />
      {pendingFile && (
        <BankSelectDialog
          file={pendingFile}
          onClose={() => { setPendingFile(null); if (inputRef.current) inputRef.current.value = '' }}
          onConfirm={(file, bank) => uploadMutation.mutate({ file, bank })}
        />
      )}

      {uploadError && <Alert severity="error" sx={{ mb: 2 }} onClose={() => setUploadError(null)}>{uploadError}</Alert>}
      {error && <Alert severity="error" sx={{ mb: 2 }}>Не удалось загрузить транзакции</Alert>}

      {isLoading && Array.from({ length: 3 }).map((_, i) => (
        <Skeleton key={i} variant="rectangular" height={56} sx={{ mb: 1, borderRadius: 2 }} />
      ))}

      {!isLoading && !data?.length && (
        <Paper elevation={0} sx={{ border: '1px solid #e0e0e0', borderRadius: 2, py: 6, textAlign: 'center' }}>
          <Typography color="text.secondary">Загрузите PDF или CSV выписку</Typography>
        </Paper>
      )}

      {groups.map(([key, txs]) => (
        <MonthGroup
          key={key}
          monthKey={key}
          txs={txs}
          defaultExpanded={key === latestKey}
          onDelete={id => deleteMutation.mutate(id)}
          onUpdate={tx => updateMutation.mutate(tx)}
        />
      ))}
    </Box>
  )
}
