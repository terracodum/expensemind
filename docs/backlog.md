# Backlog

Вещи обсуждённые но не реализованные.

---

## Backend

### parser/ — валидатор
- Добавить общий интерфейс `Validator` в `parser/interface.go`
- `CSVValidator` — проверяет наличие колонок `date`, `amount`, `description`, `category`, возвращает `INVALID_CSV_FORMAT`
- `PDFValidator` — проверяет что PDF от нужного банка (Т-Банк), возвращает `INVALID_PDF_FORMAT`
- Каждый парсер получает свой валидатор и вызывает его перед парсингом

### repository/sqlite/transaction.go
- `FindByDateRange` — дублирует `FindAll(filters)`, рассмотреть удаление
