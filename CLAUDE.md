# ExpenseMind

PROJECT TYPE: student project
DOMAIN: personal finance + analytics + ML

---

## ARCHITECTURE

Frontend (React)
↓
Go Backend (API + business logic) ↔ Database (PostgreSQL / SQLite)
↓
Python ML Service (stateless)

---

## CORE PRINCIPLES

- Go backend is the ONLY entry point
- Python service is stateless
- Python service:
  - MUST NOT access database
  - MUST NOT parse CSV
- CSV ingestion happens ONLY in Go
- All communication is HTTP + JSON
- Business logic exists ONLY in Go

---

## DATA FLOW

CSV → Go → Database → Go → ML → Go → Frontend

---

## RESPONSIBILITIES

GO BACKEND:
- REST API (/api/v1)
- transaction storage
- CSV parsing
- MCC normalization
- data aggregation (timeseries)
- ML service calls
- alert generation

PYTHON ML SERVICE:
- preprocessing (pandas)
- feature engineering
- forecasting (regression)
- return predictions

FRONTEND:
- UI
- data visualization
- API interaction

---

## DATA MODEL

Transaction:

id: integer
amount: float
description: string
mcc: integer
date: YYYY-MM-DD

RULES:
- amount < 0 → expense
- amount > 0 → income
- category derived from MCC

---

## API (PUBLIC)

BASE: /api/v1

ENDPOINTS:
- GET /transactions
- POST /transactions/upload
- GET /analytics/forecast

---

## ML API (INTERNAL)

ENDPOINT:
POST /internal/v1/predict

REQUEST:

{
  "timeseries": [
    { "t": 1, "balance": 1200 },
    { "t": 2, "balance": 1000 },
    { "t": 3, "balance": 900 }
  ],
  "horizon": 30,
  "features": {
    "avg_daily_expense": 180.0,
    "income_events": [
      { "t": 15, "amount": 2000 }
    ]
  }
}

RESPONSE:

{
  "forecast": [
    { "t": 4, "balance": 850 }
  ],
  "predicted_balance": 9800.0,
  "confidence": 0.82
}

---

## ML SERVICE CONSTRAINTS

- stateless
- JSON only
- no DB access
- no CSV processing
- no business logic

---

## PROJECT STRUCTURE

expensemind/
├── backend/
│   ├── cmd/
│   ├── internal/
│   ├── api/
│   └── go.mod
│
├── ml/
│   ├── app/
│   ├── models/
│   ├── requirements.txt
│   └── main.py
│
├── frontend/
│   ├── src/
│   └── package.json
│
├── docs/
│   ├── architecture.md
│   └── api.yaml
│
├── docker-compose.yml
├── .env.example
├── README.md
└── CLAUDE.md

---

## TECH STACK

- Go — backend
- Python (FastAPI, pandas, scikit-learn) — ML
- React — frontend
- PostgreSQL / SQLite — database

---

## AI INSTRUCTIONS

- NEVER move business logic to Python
- NEVER allow ML service to access DB
- ALWAYS follow API contracts strictly
- KEEP solutions simple
- AVOID overengineering


## CHANGE POLICY

Architecture and rules may evolve.

If user instructions contradict this file:
- follow the user instructions
- treat this file as default, not strict

Always prioritize latest user request over this document.