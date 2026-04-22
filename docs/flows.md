# Flows

## Upload

POST /transactions/upload

```
handler → service.UploadTransactions(contentType, file)
→ parser.Factory.Create(contentType)
→ parser.Parse(file)
→ repository.SaveAll(txs)
```

---

## Get Transactions

GET /transactions

```
handler → service.GetTransactions(filters)
→ repository.FindAll(filters)
```

---

## Forecast (async)

POST /analytics/forecast

```
handler → service.CreateForecast()
→ ForecastJobRepository.Create()        ← возвращает job_id
→ go worker()
```

worker:
```
→ TransactionRepository.FindForForecast(from, to)
→ RecurringRuleRepository.FindActive(today)

→ past    = transactions (<= today)
→ future  = generate(recurring_rules, > today)

→ timeseries = past (только прошлое, t=1..N)
→ horizon    = количество дней прогноза от последней транзакции

→ future используется как income_events в features:
   каждая future-точка → IncomeEvent{t, amount, label}
   t считается относительно последней точки past

→ ml.Predict(PredictRequest{timeseries, horizon, features{income_events: future}})

→ ForecastJobRepository.Update(job{status: done, result})
```

---

## Forecast status

GET /analytics/forecast/:job_id

```
handler → service.GetForecastJob(id)
→ ForecastJobRepository.FindByID(id)
```
