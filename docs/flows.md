# Flows

## Upload

POST /transactions/upload

handler → service
→ parser.Parse()
→ repository.SaveAll()

---

## Get Transactions

GET /transactions

handler → service
→ repository.FindAll(filters)

---

## Forecast (async)

POST /analytics/forecast

handler → service
→ repository.SaveJob()

worker:
→ repository.GetTransactions()
→ repository.GetRecurringRules()

→ today = now()

→ past   = transactions (<= today)
→ future = generate(recurring_rules, > today)

→ timeseries = past + future

→ ml.Predict()

→ repository.SaveForecastResult()

---

## Forecast status

GET /analytics/forecast/:job_id
