# Архитектура

## Общая схема

```
Браузер (React + MUI)
        │  HTTP/JSON
Go Backend (chi router)
        │              └── SQLite (db-data volume)
        │  HTTP/JSON
Python ML-сервис (FastAPI, stateless)
```

## Инварианты

- Go — единственная точка входа для внешних запросов
- ML-сервис stateless: не имеет доступа к БД, не парсит файлы, не содержит бизнес-логики
- Парсинг PDF — только в Go через `pdftotext` (poppler-utils)
- Взаимодействие между сервисами — HTTP + JSON

## Слои (Go backend)

| Слой | Пакет | Ответственность |
|------|-------|----------------|
| HTTP | `handler/` | Приём запросов, формирование ответов |
| Бизнес-логика | `service/` | Оркестрация, агрегация, воркер прогнозов |
| Данные | `repository/sqlite/` | Чтение и запись в SQLite |
| Парсинг | `parser/`, `parser/pdf/` | CSV и PDF → `[]Transaction` |
| ML-клиент | `ml/` | HTTP-клиент к Python-сервису |
| Модели | `domain/` | Чистые структуры без зависимостей |
| Ошибки | `errors/` | Типизированные ошибки с кодами |

## Направление зависимостей

```
handler → service → repository
                 → parser
                 → ml
domain  ← все зависят от него
errors  ← все зависят от него
```

Зависимости только вниз. `repository` не знает про `service`. `service` не знает про `handler`.

## Финансовая модель

- `transactions` — фактические операции (прошлое и настоящее), хранятся в БД
- `recurring_rules` — правила регулярных операций (зарплата и т.п.), хранятся в БД
- Будущее не хранится в БД — генерируется в service и передаётся в ML

## Воркер прогнозов

Прогноз запускается асинхронно:

```
POST /analytics/forecast
  → service создаёт ForecastJob (status=pending)
  → запускает горутину-воркер
  → возвращает job_id немедленно

Воркер:
  → читает все транзакции
  → агрегирует в timeseries по дням
  → извлекает recurring_rules как income_events
  → вызывает ml.Predict()
  → сохраняет результат в job (status=done)

GET /analytics/forecast
  → возвращает все jobs с результатами
```

## PDF-парсинг

Т-Банк использует библиотеку `ledongthuc/pdf` (координатный подход).
ВТБ, Сбербанк, Ozon — через `pdftotext -layout` (poppler-utils) + построчный парсер.

Парсеры изолированы в `parser/pdf/`, не знают о БД и сервисе.
