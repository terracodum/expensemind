# ExpenseMind

Локальный инструмент для учёта личных финансов студента с аналитикой и ML-прогнозированием баланса.

Разворачивается на твоём компьютере — данные никуда не уходят.

---

## Установка

1. Установи [Docker Desktop](https://www.docker.com/products/docker-desktop)
2. Склонируй репозиторий
```bash
git clone https://github.com/terracodum/expensemind.git
cd expensemind
```
3. Запусти
```bash
docker-compose up
```
4. Открой браузер на `http://localhost:3000`

---

## Что умеет

- Загрузка выписок из Т-Банка (PDF) — можно загрузить несколько выписок из разных банков
- Ручная категоризация расходов через UI
- Аналитика по категориям и периодам
- Прогноз баланса на 30 дней вперёд (Prophet — учитывает сезонность и регулярные доходы)
- Редактирование и удаление транзакций

> Для точного прогноза рекомендуем загрузить историю за последний год при первом запуске.

---

## Архитектура

```
Браузер (React)
      ↓
Go Backend  ←→  SQLite
      ↓
Python ML сервис (stateless)
```

Go backend — единственная точка входа. ML сервис не имеет доступа к БД, не парсит PDF, не содержит бизнес-логики.

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
│   │   │   └── forecast.go                # структура Forecast
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
│   ├── ml_spec.md                         # ТЗ для ML разработчика
│   ├── architecture.md                    # детальная архитектура
│   └── api.yaml                           # OpenAPI спецификация
│
├── docker-compose.yml
├── .env.example
├── CLAUDE.md
├── TEAM.md
└── README.md
```

---

## Слои и ответственность

### `internal/errors/` — сквозной слой
Фундамент. Не зависит ни от кого, все зависят от него.

### `internal/domain/` — модели данных
Чистые структуры. Не знает про HTTP, БД, ML.

### `internal/repository/` — база данных
Только чтение и запись. Не знает про бизнес-логику.

### `internal/pdf/` — парсер PDF
Читает выписку Т-Банка, возвращает `[]Transaction`. Ничего не сохраняет.

### `internal/ml/` — клиент ML сервиса
HTTP клиент к Python. Формирует запрос, возвращает прогноз.

### `internal/service/` — бизнес-логика
Оркестрирует все слои. Зависит от интерфейсов, не от реализаций.

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
errors   ←  все зависят от него
```

Зависимости только вниз. `repository` не знает про `service`. `service` не знает про `handler`.

---

## Поток данных

**POST /transactions/upload**
```
handler → service.UploadPDF()
            → pdf.Parser.Parse()    → []Transaction
            → mcc.ToCategory()     → категория для каждой транзакции
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
            → агрегация в timeseries    ← бизнес-логика в Go
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
    category    TEXT,
    date        TEXT    NOT NULL
);

CREATE TABLE recurring_income (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL    NOT NULL,
    day    INTEGER NOT NULL,
    label  TEXT
);
```

- `amount < 0` — расход
- `amount > 0` — доход
- `category` определяется из описания операции

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

## Публичное API

```
Base: /api/v1

GET  /transactions          — список транзакций
POST /transactions/upload   — загрузка PDF выписки
GET  /analytics/forecast    — прогноз баланса
```

Полное описание: `docs/api.yaml`

---

## ML API (внутренний)

```
POST /internal/v1/predict
GET  /health
```

Полная спецификация: `docs/ml_spec.md`

---

## Как добавлять новый функционал

Новые фичи = новые файлы. Старый код не трогается.

**Добавить поддержку другого банка:**
```
1. pdf/sber_parser.go   ← новый файл
2. main.go              ← выбор парсера по типу файла
```

**Добавить новый тип аналитики:**
```
1. repository/  ← новый метод
2. service/     ← новый метод
3. handler/     ← новый роут
```

**Добавить авторизацию:**
```
1. handler/middleware.go  ← новый middleware
```

---

## Тесты

```
internal/
├── pdf/
│   ├── parser.go
│   └── parser_test.go              # белый ящик
└── service/
    ├── transaction_service.go
    └── transaction_service_test.go  # чёрный ящик, мокаем зависимости
```

---

## Технологии

| Слой     | Технология                            |
|----------|---------------------------------------|
| Backend  | Go 1.26                               |
| ML       | Python, FastAPI, pandas, scikit-learn |
| Frontend | React                                 |
| База     | SQLite                                |
| Запуск   | Docker, Docker Compose                |