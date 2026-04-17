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
prophet      — прогнозирование временных рядов (основная модель)
scikit-learn — вспомогательные метрики и утилиты
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

Используем **Prophet** (Facebook) — алгоритм для временных рядов с сезонностью.

**Почему Prophet, а не LinearRegression:**
- Временной ряд одного пользователя — мало данных, Prophet специально для этого
- Автоматически находит паттерны: стипендия каждые 15 числа, расходы по выходным
- Чем дольше пользователь использует сервис, тем точнее прогноз (больше данных в БД)
- LinearRegression построит прямую и не учтёт сезонность

```python
# models/forecaster.py

import pandas as pd
import numpy as np
from prophet import Prophet

class Forecaster:
    def predict(self, timeseries, horizon, features):
        # Prophet ожидает колонки ds (datetime) и y (значение)
        # t — порядковый номер дня, преобразуем в дату
        base_date = pd.Timestamp.today().normalize()
        df = pd.DataFrame([
            {
                "ds": base_date - pd.Timedelta(days=len(timeseries) - p.t),
                "y":  p.balance,
            }
            for p in timeseries
        ])

        model = Prophet(daily_seasonality=False, weekly_seasonality=True)

        # Добавляем ожидаемые поступления как внешние события
        for event in features.income_events:
            event_date = base_date + pd.Timedelta(days=event.t - len(timeseries))
            model.add_regressor(f"income_{event.t}")
            df[f"income_{event.t}"] = 0.0

        model.fit(df)

        future = model.make_future_dataframe(periods=horizon)
        forecast_df = model.predict(future)

        future_rows = forecast_df.tail(horizon)
        last_t = int(timeseries[-1].t)

        # Учитываем ожидаемые поступления
        income_map = {e.t: e.amount for e in features.income_events}
        predicted = future_rows["yhat"].values.copy()
        for i in range(horizon):
            t = last_t + i + 1
            if t in income_map:
                predicted[i] += income_map[t]

        forecast = [
            {"t": last_t + i + 1, "balance": float(b)}
            for i, b in enumerate(predicted)
        ]

        # Уверенность — насколько узкий доверительный интервал
        interval_width = (future_rows["yhat_upper"] - future_rows["yhat_lower"]).mean()
        avg_balance = abs(future_rows["yhat"].mean()) + 1e-9
        confidence = float(max(0.0, min(1.0, 1.0 - interval_width / avg_balance / 2)))

        return forecast, float(predicted[-1]), confidence
```

> Это базовая реализация на Prophet. Можно улучшать — добавлять yearly_seasonality
> при накоплении данных за год+, тюнить changepoint_prior_scale.
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
