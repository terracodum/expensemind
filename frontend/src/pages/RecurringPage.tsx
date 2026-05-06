import { useQuery } from '@tanstack/react-query'
import {
  Table, TableHead, TableRow, TableCell, TableBody,
  Typography, CircularProgress, Alert,
} from '@mui/material'

interface RecurringRule {
  ID: number
  SourceID: string
  Type: string
  Amount: number
  Day: number
  Label: string
}

export default function RecurringPage() {
  const { data, isLoading, error } = useQuery<RecurringRule[]>({
    queryKey: ['recurring'],
    queryFn: () => fetch('/recurring').then(r => r.json()),
  })

  if (isLoading) return <CircularProgress />
  if (error) return <Alert severity="error">Ошибка загрузки</Alert>
  if (!data?.length) return <Typography>Нет регулярных правил</Typography>

  return (
    <Table size="small">
      <TableHead>
        <TableRow>
          <TableCell>Название</TableCell>
          <TableCell>Тип</TableCell>
          <TableCell>День месяца</TableCell>
          <TableCell align="right">Сумма</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {data.map(r => (
          <TableRow key={r.ID}>
            <TableCell>{r.Label}</TableCell>
            <TableCell>{r.Type}</TableCell>
            <TableCell>{r.Day}</TableCell>
            <TableCell align="right">{r.Amount.toFixed(2)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
