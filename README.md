# ExpenseMind

ExpenseMind — приложение для учёта личных финансов студента с аналитикой и прогнозированием баланса.

Проект построен на микросервисной архитектуре с разделением ответственности.

---

## Архитектура

Frontend (React)
↓
Go Backend (API + бизнес-логика) ←→ Database (PostgreSQL / SQLite)
↓
Python ML Service (аналитика и прогноз)

---

## Возможности

- Учёт транзакций
- Импорт CSV из банка
- Категоризация по MCC кодам
- Аналитика расходов
- Прогноз баланса
- Уведомления

---

## Структура проекта

expensemind/
├── backend/              # Go backend (API + логика)
│   ├── cmd/
│   ├── internal/
│   ├── api/              # oapi-codegen
│   └── go.mod
│
├── ml/                   # Python ML сервис
│   ├── app/
│   ├── models/
│   ├── requirements.txt
│   └── main.py
│
├── frontend/             # React приложение
│   ├── src/
│   └── package.json
│
├── docs/                 # документация
│   ├── architecture.md
│   └── api.yaml
│
├── docker-compose.yml
├── .env.example
├── README.md
└── CLAUDE.md

---

## Поток данных

CSV → Go → база данных → Go → ML сервис → Go → Frontend

---

## Распределение ответственности

Go backend:
- REST API (/api/v1)
- хранение данных
- парсинг CSV
- работа с MCC
- агрегация
- вызов ML
- генерация алертов

Python ML сервис:
- preprocessing (pandas)
- feature engineering
- прогноз
- возврат результата

Важно:
- не имеет доступа к БД
- не обрабатывает CSV
- не содержит бизнес-логики

Frontend:
- интерфейс
- визуализация
- работа с API

---

## Модель данных

Пример транзакции:

id: 1
amount: -500.0
description: Пятёрочка
mcc: 5411
date: 2026-04-10

Правила:
- отрицательное значение — расход
- положительное — доход
- категория определяется через MCC

---

## Запуск

1. Клонировать репозиторий

git clone https://github.com/yourname/expensemind.git
cd expensemind

2. Запустить через Docker

docker-compose up --build

3. Доступ

Backend: http://localhost:8080
ML: http://localhost:5001
Frontend: http://localhost:3000

---

## API

Публичное API: /api/v1
ML API: /internal/v1

Описание API: docs/api.yaml

---

## Важно

- Backend — единственная точка доступа к данным
- ML сервис — stateless
- все данные проходят через backend
- архитектура разделена по зонам ответственности

---

## Технологии

Go — backend
Python (FastAPI, pandas, scikit-learn) — ML
React — frontend
PostgreSQL / SQLite — база данных