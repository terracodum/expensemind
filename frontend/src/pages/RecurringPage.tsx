import { useState } from 'react'
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query'
import {
  Box, Button, Paper, Stack, Table, TableBody, TableCell,
  TableContainer, TableHead, TableRow, Typography,
  Alert, IconButton, Chip, Skeleton, Dialog, DialogTitle,
  DialogContent, DialogActions, TextField, Select, MenuItem,
  FormControl, InputLabel,
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline'

interface RecurringRule {
  ID: number
  SourceID: string
  Type: string
  Amount: number
  Day: number
  Label: string
}

interface RuleForm {
  Label: string
  Type: string
  Amount: string
  Day: string
}

const EMPTY_FORM: RuleForm = { Label: '', Type: 'expense', Amount: '', Day: '' }

export default function RecurringPage() {
  const qc = useQueryClient()
  const [open, setOpen] = useState(false)
  const [form, setForm] = useState<RuleForm>(EMPTY_FORM)
  const [formError, setFormError] = useState<string | null>(null)

  const { data, isLoading, error } = useQuery<RecurringRule[]>({
    queryKey: ['recurring'],
    queryFn: () => fetch('/recurring').then(r => r.json()),
  })

  const addMutation = useMutation({
    mutationFn: (rule: object) =>
      fetch('/recurring', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rule),
      }).then(async r => {
        if (!r.ok) throw new Error((await r.json().catch(() => ({}))).message ?? 'Ошибка')
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['recurring'] })
      setOpen(false)
      setForm(EMPTY_FORM)
      setFormError(null)
    },
    onError: (e: Error) => setFormError(e.message),
  })

  const deleteMutation = useMutation({
    mutationFn: (sourceID: string) =>
      fetch(`/recurring/${encodeURIComponent(sourceID)}`, { method: 'DELETE' }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['recurring'] }),
  })

  function handleSubmit() {
    const amount = parseFloat(form.Amount)
    const day = parseInt(form.Day)
    if (!form.Label.trim()) return setFormError('Введите название')
    if (isNaN(amount) || amount <= 0) return setFormError('Сумма должна быть больше 0')
    if (isNaN(day) || day < 1 || day > 28) return setFormError('День должен быть от 1 до 28')

    addMutation.mutate({
      SourceID: crypto.randomUUID(),
      Label: form.Label.trim(),
      Type: form.Type,
      Amount: amount,
      Day: day,
      StartDate: new Date().toISOString(),
    })
  }

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 3 }}>
        <Box>
          <Typography variant="h5" fontWeight={700}>Регулярные операции</Typography>
          {data && (
            <Typography variant="body2" color="text.secondary">
              {data.length} правил
            </Typography>
          )}
        </Box>
        <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpen(true)} disableElevation>
          Добавить
        </Button>
      </Stack>

      {error && <Alert severity="error" sx={{ mb: 2 }}>Не удалось загрузить правила</Alert>}

      <Paper elevation={0} sx={{ border: '1px solid #e0e0e0', borderRadius: 2 }}>
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow sx={{ '& th': { bgcolor: '#fafafa', fontWeight: 600 } }}>
                <TableCell>Название</TableCell>
                <TableCell>Тип</TableCell>
                <TableCell>День месяца</TableCell>
                <TableCell align="right">Сумма</TableCell>
                <TableCell width={48} />
              </TableRow>
            </TableHead>
            <TableBody>
              {isLoading && Array.from({ length: 3 }).map((_, i) => (
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
                    Нет регулярных операций
                  </TableCell>
                </TableRow>
              )}
              {data?.map(r => (
                <TableRow key={r.ID} hover>
                  <TableCell sx={{ fontWeight: 500 }}>{r.Label}</TableCell>
                  <TableCell>
                    <Chip
                      label={r.Type === 'income' ? 'Доход' : 'Расход'}
                      size="small"
                      color={r.Type === 'income' ? 'success' : 'error'}
                      variant="outlined"
                    />
                  </TableCell>
                  <TableCell sx={{ color: 'text.secondary' }}>{r.Day} числа</TableCell>
                  <TableCell align="right" sx={{ fontVariantNumeric: 'tabular-nums', fontWeight: 500,
                    color: r.Type === 'income' ? 'success.main' : 'text.primary' }}>
                    {r.Type === 'income' ? '+' : ''}{r.Amount.toFixed(2)}
                  </TableCell>
                  <TableCell>
                    <IconButton size="small" onClick={() => deleteMutation.mutate(r.SourceID)}
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

      <Dialog open={open} onClose={() => { setOpen(false); setFormError(null); setForm(EMPTY_FORM) }} maxWidth="xs" fullWidth>
        <DialogTitle sx={{ fontWeight: 700 }}>Новое правило</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 1 }}>
            {formError && <Alert severity="error" onClose={() => setFormError(null)}>{formError}</Alert>}
            <TextField
              label="Название"
              value={form.Label}
              onChange={e => setForm(f => ({ ...f, Label: e.target.value }))}
              fullWidth size="small"
              placeholder="Аренда, зарплата..."
            />
            <FormControl size="small" fullWidth>
              <InputLabel>Тип</InputLabel>
              <Select label="Тип" value={form.Type} onChange={e => setForm(f => ({ ...f, Type: e.target.value }))}>
                <MenuItem value="income">Доход</MenuItem>
                <MenuItem value="expense">Расход</MenuItem>
              </Select>
            </FormControl>
            <TextField
              label="Сумма"
              value={form.Amount}
              onChange={e => setForm(f => ({ ...f, Amount: e.target.value }))}
              fullWidth size="small" type="number"
              inputProps={{ min: 0, step: '0.01' }}
            />
            <TextField
              label="День месяца"
              value={form.Day}
              onChange={e => setForm(f => ({ ...f, Day: e.target.value }))}
              fullWidth size="small" type="number"
              inputProps={{ min: 1, max: 28 }}
              helperText="От 1 до 28"
            />
          </Stack>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={() => { setOpen(false); setFormError(null); setForm(EMPTY_FORM) }}>
            Отмена
          </Button>
          <Button variant="contained" disableElevation onClick={handleSubmit} disabled={addMutation.isPending}>
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
