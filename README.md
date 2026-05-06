# ExpenseMind

Локальный инструмент учёта личных финансов с ML-прогнозированием баланса.

Работает полностью на твоём компьютере — данные никуда не уходят.

---

## Быстрый старт

### Требования

- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- Git

### Установка

```bash
git clone https://github.com/terracodum/expensemind.git
cd expensemind
sudo make install
```

После этого из любого места терминала:

```bash
expensemind
```

Открой браузер: **http://localhost**

### Без установки

```bash
docker compose up --build
```

---

## Что умеет

- **Загрузка выписок** — PDF из Т-Банка, ВТБ, Сбербанка, Ozon Банка, а также CSV
- **Ручное добавление** транзакций через форму
- **Категоризация** расходов и доходов
- **Прогноз баланса** на 30 дней вперёд (Prophet, с учётом сезонности)
- **График прогноза** по дням с датами
- **История прогнозов** — все запуски сохраняются

---

## Загрузка выписок

### PDF

Поддерживаемые банки:

| Банк | Формат |
|------|--------|
| Т-Банк | PDF выписка |
| ВТБ | PDF выписка |
| Сбербанк | PDF выписка |
| Ozon Банк | PDF выписка |

При загрузке PDF система спросит, из какого банка файл.

### CSV

Формат:

```csv
date;amount;description;category
01.04.2026;-500.00;Магнит;food
15.04.2026;5000.00;Стипендия;income
```

- Разделитель — `;`
- Дата — `ДД.ММ.ГГГГ`
- Amount — отрицательное расход, положительное доход
- Category — можно оставить пустым

---

## Архитектура

```
Браузер (React + MUI)
        │
    HTTP/JSON
        │
  Go Backend  ──── SQLite
        │
    HTTP/JSON
        │
Python ML-сервис (stateless)
```

**Принципы:**
- Go — единственная точка входа, вся бизнес-логика в нём
- ML-сервис только считает прогноз, не имеет доступа к БД
- Парсинг PDF и CSV — только в Go (через `pdftotext` + собственные парсеры)

---

## Структура проекта

```
expensemind/
├── backend/
│   ├── cmd/main.go                  # точка входа
│   └── internal/
│       ├── domain/                  # Transaction, Forecast, ForecastJob
│       ├── errors/                  # типизированные ошибки
│       ├── handler/                 # HTTP-слой (chi router)
│       ├── ml/                      # HTTP-клиент к ML-сервису
│       ├── parser/
│       │   ├── csv.go
│       │   └── pdf/
│       │       ├── tbank.go         # Т-Банк
│       │       ├── vtb.go           # ВТБ (через pdftotext)
│       │       ├── sber.go          # Сбербанк (через pdftotext)
│       │       └── ozon.go          # Ozon Банк (через pdftotext)
│       ├── repository/sqlite/       # SQLite-репозиторий
│       └── service/                 # бизнес-логика, воркер прогнозов
│
├── ml/
│   └── app/
│       ├── models/forecaster.py     # Prophet-модель
│       ├── routes/predict.py        # POST /internal/v1/predict
│       └── schemas/                 # Pydantic-схемы
│
├── frontend/
│   └── src/pages/
│       ├── TransactionsPage.tsx
│       └── ForecastPage.tsx
│
├── docs/                            # архитектура, контракты, спека ML
├── docker-compose.yml
├── Makefile
└── expensemind                      # CLI-скрипт
```

---

## API

```
POST /api/v1/transactions/upload?bank=vtb   — загрузка PDF/CSV
POST /api/v1/transactions                   — добавить транзакцию вручную
GET  /api/v1/transactions                   — список транзакций
GET  /api/v1/analytics/forecast             — история прогнозов
POST /api/v1/analytics/forecast             — запустить новый прогноз
```

Подробнее: [`docs/contracts.md`](docs/contracts.md)

---

## Технологии

| Слой | Технология |
|------|-----------|
| Frontend | React, MUI, TanStack Query, Recharts |
| Backend | Go 1.26, chi, SQLite |
| ML | Python, FastAPI, Prophet |
| PDF-парсинг | pdftotext (poppler) |
| Запуск | Docker Compose |

---

## Документация

- [`docs/architecture.md`](docs/architecture.md) — слои и зависимости
- [`docs/flows.md`](docs/flows.md) — потоки данных
- [`docs/contracts.md`](docs/contracts.md) — API-контракты
- [`docs/ml_spec.md`](docs/ml_spec.md) — спецификация ML-сервиса
- [`docs/db.md`](docs/db.md) — схема БД
- [`docs/errors.md`](docs/errors.md) — коды ошибок
