# Backlog

Вещи обсуждённые но не реализованные.

---

## Backend

### parser/ — валидатор
- Добавить общий интерфейс `Validator` в `parser/interface.go`
- `CSVValidator` — проверяет наличие колонок `date`, `amount`, `description`, `category`, возвращает `INVALID_CSV_FORMAT`
- `PDFValidator` — проверяет что PDF от нужного банка (Т-Банк), возвращает `INVALID_PDF_FORMAT`
- Каждый парсер получает свой валидатор и вызывает его перед парсингом

### parser/csv.go
- Колонка `category` опциональна — если отсутствует или пустая, подставлять `"unknown"`. Пользователь проставит категории в UI.

### repository/sqlite/transaction.go
- `FindByDateRange` — переименовать в `FindForForecast`, чтобы была понятна семантика: подбор данных для ML без пользовательских фильтров.
