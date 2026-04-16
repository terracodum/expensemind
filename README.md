# ExpenseMind

ExpenseMind — приложение для учёта личных финансов студента с аналитикой и прогнозированием баланса на основе ML.

---

## Архитектура

```
Frontend (React)
        ↓
Go Backend (API + бизнес-логика) ←→ Database (SQLite / PostgreSQL)
        ↓
Python ML Service (stateless, только вычисления)
```

Go backend — единственная точка входа в систему. ML сервис не имеет доступа к БД и не содержит бизнес-логики.

---

## Возможности

- Учёт транзакций (доходы и расходы)
- Импорт CSV из банковских выписок
- Категоризация по MCC кодам
- Аналитика расходов по категориям
- Прогноз баланса на основе ML
- Уведомления при превышении расходов

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
│   │   ├── csv/
│   │   │   ├── interface.go               # Parser интерфейс
│   │   │   ├── parser.go                  # реализация
│   │   │   └── validator.go               # валидация строк CSV
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
│   │   │   └── forecaster.go              # логика прогнозирования
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
│   ├── ml_spec.md                         # ТЗ для ML разработчика
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
Фундамент системы. Не зависит ни от одного другого пакета.
Все остальные слои зависят от него.

```
AppError интерфейс — что умеет любая ошибка в системе
codes.go           — единый список кодов ошибок
constructors.go    — удобные функции создания ошибок
```

Коды ошибок проекта:
```
INTERNAL_ERROR
VALIDATION_ERROR
NOT_FOUND
CSV_PARSE_ERROR
CSV_INVALID_FORMAT
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

### `internal/repository/` — работа с базой данных
Только сохранение и чтение данных. Не знает про HTTP и бизнес-логику.

Интерфейс:
```
Save(tx Transaction) error
SaveAll(txs []Transaction) error
FindAll(filters Filters) ([]Transaction, error)
FindByDateRange(from, to time.Time) ([]Transaction, error)
```

Хочешь добавить PostgreSQL — создаёшь `repository/postgres/transaction.go`,
меняешь одну строчку в `main.go`. Всё остальное не трогаешь.

### `internal/csv/` — парсинг CSV
Читает файл, возвращает `[]Transaction`. Ничего не сохраняет.

```
Parser интерфейс:
  Parse(file io.Reader) ([]Transaction, AppError)
```

### `internal/ml/` — клиент ML сервиса
HTTP клиент к Python сервису. Формирует запрос, возвращает прогноз.

```
MLClient интерфейс:
  Predict(req PredictRequest) (Forecast, AppError)
```

### `internal/service/` — бизнес-логика
Оркестрирует все остальные слои. Здесь живут правила приложения.

```
TransactionService интерфейс:
  UploadCSV(file io.Reader) (int, AppError)
  GetTransactions(filters Filters) ([]Transaction, AppError)
  GetForecast(horizon int) (Forecast, AppError)
```

Зависит от интерфейсов `RepositoryI`, `ParserI`, `MLClientI` —
не от конкретных реализаций.

### `internal/handler/` — HTTP слой
Принять запрос → вызвать service → вернуть ответ. Не знает про БД и ML.

---

## Направление зависимостей

```
main.go
  │
  ├── создаёт repository
  ├── создаёт csv.Parser
  ├── создаёт ml.Client
  ├── создаёт service (получает repository, parser, ml)
  └── создаёт handler (получает service)

handler  →  ServiceInterface
service  →  RepositoryInterface
service  →  ParserInterface
service  →  MLClientInterface
errors   ←  все слои зависят от него
```

Зависимости идут только вниз. `repository` никогда не знает про `service`.
`service` никогда не знает про `handler`.

---

## Поток данных

**POST /transactions/upload**
```
handler → service.UploadCSV()
            → csv.Parser.Parse()     → []Transaction
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

## Модель данных

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

Правила:
- `amount < 0` — расход
- `amount > 0` — доход
- `category` определяется через MCC код

---

## Формат ошибок

Все слои возвращают ошибки в едином формате:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "transaction not found"
  }
}
```

---

## Публичное API

```
Base path: /api/v1

GET  /transactions              — список транзакций
POST /transactions/upload       — загрузка CSV
GET  /analytics/forecast        — прогноз баланса
```

Полное описание: `docs/api.yaml`

---

## ML API (внутренний)

```
POST /internal/v1/predict
GET  /health
```

Полная спецификация для ML разработчика: `docs/ml_spec.md`

---

## Как добавлять новый функционал

Архитектура построена так, чтобы новые фичи требовали только новых файлов.
Старый код не трогается.

**Пример — добавить PostgreSQL:**
```
1. создать repository/postgres/transaction.go
2. в main.go изменить одну строчку
```

**Пример — добавить новый тип аналитики:**
```
1. repository/ — новый метод в интерфейсе + реализация
2. service/    — новый метод в интерфейсе + реализация
3. handler/    — новый роут
```

**Пример — добавить авторизацию:**
```
1. handler/middleware.go — новый middleware
```

---

## Тесты

```
internal/
├── csv/
│   ├── parser.go
│   └── parser_test.go          # белый ящик — граничные случаи парсинга
└── service/
    ├── transaction_service.go
    └── transaction_service_test.go  # чёрный ящик — мокаем repository и ml
```

Приоритет покрытия: `csv/` и `service/` — там живёт критическая логика.

---

## Запуск

1. Клонировать репозиторий

```bash
git clone https://github.com/yourname/expensemind.git
cd expensemind
```

2. Запустить через Docker

```bash
docker-compose up --build
```

3. Доступ

```
Backend:  http://localhost:8080
ML:       http://localhost:8001
Frontend: http://localhost:3000
```

---

## Технологии

| Слой     | Технология                              |
|----------|-----------------------------------------|
| Backend  | Go                                      |
| ML       | Python, FastAPI, pandas, scikit-learn   |
| Frontend | React                                   |
| База     | SQLite (dev) / PostgreSQL (prod)        |