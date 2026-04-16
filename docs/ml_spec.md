# ML Service — Technical Specification

Техническое задание для разработчика Python ML сервиса.

---

## Обзор

ML сервис является **вычислительным модулем** системы ExpenseMind.  
Сервис получает готовые данные от Go backend, выполняет прогнозирование баланса и возвращает результат.

Сервис **не имеет доступа к базе данных** и **не содержит бизнес-логики**.

---

## Главное правило — stateless

Каждый запрос самодостаточен и содержит все данные необходимые для вычисления.  
Сервис не хранит никакого состояния между запросами.

```
❌ нельзя:                ✅ можно:
  читать БД                принимать JSON
  писать в БД              считать
  читать файлы             возвращать JSON
  хранить сессии
  кешировать данные
  содержать бизнес-логику
```

---

## Tech Stack

```
Python       3.11+
FastAPI      — HTTP сервер
pandas       — обработка временных рядов
scikit-learn — модели прогнозирования
uvicorn      — ASGI сервер
pydantic     — валидация данных
```

---

## Структура проекта

```
ml/
├── app/
│   ├── main.py              ← точка входа FastAPI
│   ├── routes/
│   │   └── predict.py       ← роут POST /internal/v1/predict
│   ├── models/
│   │   └── forecaster.py    ← логика прогнозирования
│   └── schemas/
│       ├── request.py       ← Pydantic модель запроса
│       └── response.py      ← Pydantic модель ответа
├── requirements.txt
└── Dockerfile
```

---

## API Contract

### Endpoints

```
POST /internal/v1/predict   ← основной эндпоинт прогноза
GET  /health                ← проверка доступности сервиса
```

---

### POST /internal/v1/predict

#### Request

```json
{
  "timeseries": [
    {
      "t": 1,
      "balance": 1200.0,
      "day_of_week": 1,
      "is_weekend": false,
      "food_total": 450.0,
      "transport_total": 120.0,
      "entertainment_total": 200.0,
      "avg_transaction_size": 85.0,
      "transaction_count": 9
    },
    {
      "t": 2,
      "balance": 1000.0,
      "day_of_week": 2,
      "is_weekend": false,
      "food_total": 300.0,
      "transport_total": 80.0,
      "entertainment_total": 0.0,
      "avg_transaction_size": 60.0,
      "transaction_count": 6
    }
  ],
  "horizon": 30,
  "features": {
    "avg_daily_expense": 180.0,
    "income_events": [
      {
        "t": 15,
        "amount": 5000.0,
        "label": "стипендия"
      }
    ]
  }
}
```

#### Описание полей запроса

**timeseries[]** — временной ряд, одна точка = один день:

| Поле                  | Тип     | Описание                              |
|-----------------------|---------|---------------------------------------|
| `t`                   | int     | порядковый номер дня (1, 2, 3...)     |
| `balance`             | float   | баланс на этот день                   |
| `day_of_week`         | int     | день недели (1=пн, 7=вс)             |
| `is_weekend`          | bool    | выходной день                         |
| `food_total`          | float   | сумма трат на еду за день             |
| `transport_total`     | float   | сумма трат на транспорт за день       |
| `entertainment_total` | float   | сумма трат на развлечения за день     |
| `avg_transaction_size`| float   | средний чек за день                   |
| `transaction_count`   | int     | количество транзакций за день         |

**horizon** — на сколько дней вперёд строить прогноз (1–365)

**features{}** — агрегированные характеристики:

| Поле                        | Тип     | Описание                              |
|-----------------------------|---------|---------------------------------------|
| `avg_daily_expense`         | float   | средний дневной расход (всегда > 0)   |
| `income_events[]`           | array   | ожидаемые поступления                 |
| `income_events[].t`         | int     | на какой день ожидается               |
| `income_events[].amount`    | float   | сумма поступления                     |
| `income_events[].label`     | string  | название ("стипендия", "зарплата")    |

---

#### Response (success) — HTTP 200

```json
{
  "forecast": [
    { "t": 3, "balance": 950.0 },
    { "t": 4, "balance": 850.0 }
  ],
  "predicted_balance": 9800.0,
  "confidence": 0.82
}
```

| Поле                | Тип     | Описание                                      |
|---------------------|---------|-----------------------------------------------|
| `forecast[]`        | array   | прогноз по дням                               |
| `forecast[].t`      | int     | порядковый номер дня (продолжение timeseries) |
| `forecast[].balance`| float   | предсказанный баланс на этот день             |
| `predicted_balance` | float   | баланс на последний день горизонта            |
| `confidence`        | float   | уверенность модели от 0.0 до 1.0              |

---

#### Response (error)

При любой ошибке возвращать строго в этом формате:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable description"
  }
}
```

Коды ошибок:

| Код                    | HTTP | Когда возвращать                                      |
|------------------------|------|-------------------------------------------------------|
| `VALIDATION_ERROR`     | 400  | невалидные входные данные (NaN, Inf, неверные типы)   |
| `INSUFFICIENT_DATA`    | 400  | timeseries содержит меньше 3 точек                    |
| `PREDICTION_ERROR`     | 500  | внутренняя ошибка модели (исключение scikit-learn)    |

---

### GET /health

```json
{ "status": "ok" }
```

HTTP 200 всегда если сервис запущен.  
Go backend вызывает этот эндпоинт перед отправкой запросов на прогноз.

---

## Pydantic Schemas

Go backend ожидает строгое соответствие схемам.

```python
# schemas/request.py

from pydantic import BaseModel, field_validator
from typing import List

class TimePoint(BaseModel):
    t: int
    balance: float
    day_of_week: int
    is_weekend: bool
    food_total: float
    transport_total: float
    entertainment_total: float
    avg_transaction_size: float
    transaction_count: int

class IncomeEvent(BaseModel):
    t: int
    amount: float
    label: str = ""

class Features(BaseModel):
    avg_daily_expense: float
    income_events: List[IncomeEvent] = []

class PredictRequest(BaseModel):
    timeseries: List[TimePoint]
    horizon: int
    features: Features

    @field_validator('timeseries')
    def timeseries_min_length(cls, v):
        if len(v) < 3:
            raise ValueError('timeseries must have at least 3 points')
        return v

    @field_validator('horizon')
    def horizon_valid(cls, v):
        if v < 1 or v > 365:
            raise ValueError('horizon must be between 1 and 365')
        return v
```

```python
# schemas/response.py

from pydantic import BaseModel
from typing import List

class ForecastPoint(BaseModel):
    t: int
    balance: float

class PredictResponse(BaseModel):
    forecast: List[ForecastPoint]
    predicted_balance: float
    confidence: float
```

---

## Forecaster Logic

Минимальная рабочая реализация:

```python
# models/forecaster.py

import pandas as pd
import numpy as np
from sklearn.linear_model import LinearRegression

class Forecaster:
    def predict(self, timeseries, horizon, features):
        df = pd.DataFrame([
            {
                "t":                    p.t,
                "balance":              p.balance,
                "day_of_week":          p.day_of_week,
                "is_weekend":           int(p.is_weekend),
                "food_total":           p.food_total,
                "transport_total":      p.transport_total,
                "entertainment_total":  p.entertainment_total,
                "avg_transaction_size": p.avg_transaction_size,
                "transaction_count":    p.transaction_count,
            }
            for p in timeseries
        ])

        feature_cols = [
            "t", "day_of_week", "is_weekend",
            "food_total", "transport_total", "entertainment_total",
            "avg_transaction_size", "transaction_count"
        ]

        X = df[feature_cols].values
        y = df["balance"].values

        model = LinearRegression()
        model.fit(X, y)

        last_t = int(df["t"].max())
        future_t = np.arange(last_t + 1, last_t + horizon + 1)

        # Для будущих точек используем средние значения фич
        avg_features = df[feature_cols[1:]].mean().values
        future_X = np.column_stack([
            future_t,
            np.tile(avg_features, (horizon, 1))
        ])

        predicted = model.predict(future_X)

        # Учитываем ожидаемые поступления
        income_map = {e.t: e.amount for e in features.income_events}
        for i, t in enumerate(future_t):
            if t in income_map:
                predicted[i] += income_map[t]

        r2 = model.score(X, y)
        confidence = float(max(0.0, min(1.0, r2)))

        forecast = [
            {"t": int(t), "balance": float(b)}
            for t, b in zip(future_t, predicted)
        ]

        return forecast, float(predicted[-1]), confidence
```

> Это базовая реализация. Можно улучшать модель — добавлять другие алгоритмы
> (Ridge, RandomForest), кросс-валидацию, нормализацию фич.
> Главное — не менять формат запроса и ответа без согласования с Go командой.

---

## Environment

```bash
# .env
ML_HOST=0.0.0.0
ML_PORT=8001
```

Go backend обращается к сервису по адресу:

```
http://ml-service:8001
```

---

## Что НЕ нужно реализовывать

```
❌ подключение к SQLite
❌ чтение PDF / CSV файлов
❌ хранение данных между запросами
❌ авторизация и аутентификация
❌ бизнес-логика (категоризация, агрегация транзакций)
❌ новые эндпоинты кроме /internal/v1/predict и /health
❌ изменение формата request/response без согласования с Go командой
```

---

## Checklist

```
☐ POST /internal/v1/predict работает по контракту
☐ GET /health возвращает 200
☐ Pydantic валидация на все поля запроса
☐ Все три кода ошибок реализованы
☐ HTTP статусы соответствуют спецификации
☐ Сервис не обращается к БД
☐ Dockerfile написан
☐ requirements.txt актуален
```
