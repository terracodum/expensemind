# ExpenseMind

Локальный инструмент учёта личных финансов студента с аналитикой и ML-прогнозированием баланса.
Разворачивается на своём компьютере — данные никуда не уходят. Запускается одной командой `expensemind`.

---

## Архитектура

```
Frontend (React)
        ↓
Go Backend (API + бизнес-логика) ←→ Database (SQLite / PostgreSQL)
        ↓
Python ML Service (stateless, только вычисления)
```

Go backend — единственная точка входа. ML сервис не имеет доступа к БД, не парсит PDF, не содержит бизнес-логики.

---

## Ключевые принципы

- Go backend — единственная точка входа в систему
- Python ML сервис — stateless, только вычисления
- ML сервис НЕ обращается к БД, НЕ парсит PDF, НЕ содержит бизнес-логики
- PDF парсится только в Go
- Весь обмен — HTTP + JSON
- Бизнес-логика только в Go

---

## Структура проекта

```
expensemind/
├── backend/
│   ├── cmd/
│   │   └── main.go                        # точка входа, сборка зависимостей
│   ├── internal/
│   │   ├── errors/
│   │   │   ├── errors.go                  # интерфейс AppError
│   │   │   ├── codes.go                   # все коды ошибок проекта
│   │   │   ├── impl.go                    # реализация AppError
│   │   │   └── constructors.go            # New(), Wrap(), NotFound()...
│   │   ├── domain/
│   │   │   ├── transaction.go             # структура Transaction
│   │   │   ├── forecast.go                # структура Forecast
│   │   │   └── mcc.go                     # MCC коды → категории
│   │   ├── repository/
│   │   │   ├── interface.go               # TransactionRepository интерфейс
│   │   │   └── sqlite/
│   │   │       └── transaction.go         # реализация для SQLite
│   │   ├── pdf/
│   │   │   ├── interface.go               # Parser интерфейс
│   │   │   ├── parser.go                  # парсер выписок Т-Банка
│   │   │   └── validator.go               # валидация данных
│   │   ├── ml/
│   │   │   ├── interface.go               # MLClient интерфейс
│   │   │   ├── client.go                  # HTTP клиент к Python сервису
│   │   │   └── dto.go                     # структуры запроса и ответа
│   │   ├── service/
│   │   │   ├── interface.go               # TransactionService интерфейс
│   │   │   └── transaction_service.go     # реализация бизнес-логики
│   │   └── handler/
│   │       ├── handler.go                 # регистрация роутов
│   │       ├── transaction.go             # GET /transactions, POST /upload
│   │       ├── analytics.go               # GET /analytics/forecast
│   │       └── middleware.go              # логирование, CORS
│   └── go.mod
│
├── ml/
│   ├── app/
│   │   ├── main.py                        # точка входа FastAPI
│   │   ├── routes/
│   │   │   └── predict.py                 # POST /internal/v1/predict
│   │   ├── models/
│   │   │   └── forecaster.py              # логика прогнозирования
│   │   └── schemas/
│   │       ├── request.py                 # Pydantic модель запроса
│   │       └── response.py                # Pydantic модель ответа
│   ├── requirements.txt
│   └── Dockerfile
│
├── frontend/
│   ├── src/
│   └── package.json
│
├── docs/
│   ├── architecture.md
│   ├── ml_spec.md
│   └── api.yaml
│
├── docker-compose.yml
├── .env.example
├── README.md
└── CLAUDE.md
```

---

## Слои и ответственность

### `internal/errors/` — сквозной слой
Не зависит ни от одного другого пакета. Все остальные слои зависят от него.

Коды ошибок:
```
INTERNAL_ERROR
VALIDATION_ERROR
NOT_FOUND
PDF_PARSE_ERROR
PDF_INVALID_FORMAT
ML_SERVICE_UNAVAILABLE
ML_RESPONSE_INVALID
DB_ERROR
```

### `internal/domain/` — модели данных
Чистые структуры без логики. Не знает про HTTP, БД, ML.

```go
Transaction {
    ID          int
    Amount      float64   // < 0 расход, > 0 доход
    Description string
    MCC         int
    Category    string    // определяется из MCC
    Date        time.Time
}
```

### `internal/repository/` — работа с БД
Только сохранение и чтение. Интерфейс:
```
Save(tx Transaction) error
SaveAll(txs []Transaction) error
FindAll(filters Filters) ([]Transaction, error)
FindByDateRange(from, to time.Time) ([]Transaction, error)
```

### `internal/pdf/` — парсинг PDF
Читает выписку Т-Банка, возвращает `[]Transaction`. Ничего не сохраняет.
```
Parser.Parse(file io.Reader) ([]Transaction, AppError)
```

### `internal/ml/` — клиент ML сервиса
HTTP клиент к Python. Формирует запрос, возвращает прогноз.
```
MLClient.Predict(req PredictRequest) (Forecast, AppError)
```

### `internal/service/` — бизнес-логика
Оркестрирует все слои. Зависит от интерфейсов, не от реализаций.
```
TransactionService.UploadPDF(file io.Reader) (int, AppError)
TransactionService.GetTransactions(filters Filters) ([]Transaction, AppError)
TransactionService.GetForecast(horizon int) (Forecast, AppError)
```

### `internal/handler/` — HTTP слой
Принять запрос → вызвать service → вернуть ответ.

---

## Направление зависимостей

```
main.go
  │
  ├── создаёт repository
  ├── создаёт pdf.Parser
  ├── создаёт ml.Client
  ├── создаёт service (получает repository, parser, ml)
  └── создаёт handler (получает service)

handler  →  ServiceInterface
service  →  RepositoryInterface
service  →  ParserInterface
service  →  MLClientInterface
errors   ←  все слои зависят от него
```

Зависимости идут только вниз. `repository` не знает про `service`. `service` не знает про `handler`.

---

## Потоки данных

**POST /transactions/upload**
```
handler → service.UploadPDF()
            → pdf.Parser.Parse()     → []Transaction
            → mcc.ToCategory()       → категория для каждой транзакции
            → repository.SaveAll()
          → { "uploaded": 42 }
```

**GET /transactions**
```
handler → service.GetTransactions(filters)
            → repository.FindAll(filters)
          → []Transaction
```

**GET /analytics/forecast**
```
handler → service.GetForecast(horizon)
            → repository.FindAll()
            → агрегация в timeseries по дням   ← бизнес-логика в Go
            → извлечение фич
            → ml.Client.Predict()
          → Forecast
```

---

## Схема БД

```sql
CREATE TABLE transactions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    amount      REAL    NOT NULL,
    description TEXT,
    mcc         INTEGER,
    category    TEXT,
    date        TEXT    NOT NULL
);

CREATE TABLE recurring_income (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL    NOT NULL,
    day    INTEGER NOT NULL,  -- день месяца (1-31)
    label  TEXT               -- "стипендия", "зарплата"
);
```

---

## Публичное API

```
Base: /api/v1

GET  /transactions          — список транзакций
POST /transactions/upload   — загрузка PDF выписки
GET  /analytics/forecast    — прогноз баланса
```

## ML API (внутренний)

```
POST /internal/v1/predict
GET  /health
```

Запрос:
```json
{
  "timeseries": [
    { "t": 1, "balance": 1200 },
    { "t": 2, "balance": 1000 }
  ],
  "horizon": 30,
  "features": {
    "avg_daily_expense": 180.0,
    "income_events": [
      { "t": 15, "amount": 2000 }
    ]
  }
}
```

Ответ:
```json
{
  "forecast": [
    { "t": 3, "balance": 850 }
  ],
  "predicted_balance": 9800.0,
  "confidence": 0.82
}
```

---

## Формат ошибок

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "transaction not found"
  }
}
```

---

## Технологии

| Слой     | Технология                            |
|----------|---------------------------------------|
| Backend  | Go                                    |
| ML       | Python, FastAPI, pandas, scikit-learn |
| Frontend | React                                 |
| База     | SQLite                                |
| Запуск   | Docker, Docker Compose                |

---

## Тесты

Приоритет: `pdf/` и `service/` — критическая логика.
- `pdf/parser_test.go` — белый ящик, граничные случаи
- `service/transaction_service_test.go` — чёрный ящик, мокаем repository и ml

---

## Запуск

```bash
expensemind
```

Открывается браузер на `http://localhost:3000`. Требуется Docker Desktop.

```
Backend:  http://localhost:8080
ML:       http://localhost:8001
Frontend: http://localhost:3000
```

---

## Инструкции для ИИ

- НЕ переносить бизнес-логику в Python
- НЕ давать ML сервису доступ к БД
- Строго соблюдать API контракты
- Держать решения простыми, не усложнять

## Политика изменений

Если инструкции пользователя противоречат этому файлу — следовать инструкциям пользователя. Этот файл — дефолт, не жёсткое ограничение.
