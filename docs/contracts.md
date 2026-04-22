# Contracts

## Parser

### Factory

```go
Create(contentType string) (Parser, error)
```

### Parser

```go
Parse(file io.Reader) ([]Transaction, error)
```

## Repository

### Transactions

```go
Save(tx Transaction) error
SaveAll(txs []Transaction) error
FindAll(filters Filters) ([]Transaction, error)
FindForForecast(from, to time.Time) ([]Transaction, error)
Update(tx Transaction) error
Delete(id int) error
```

### RecurringRule

```go
Save(rule RecurringRule) error
FindAll() ([]RecurringRule, error)
FindActive(today time.Time) ([]RecurringRule, error)
Delete(sourceID string) error
```

### ForecastJob

```go
Create() (int, error)
FindByID(id int) (ForecastJob, error)
FindAll() ([]ForecastJob, error)
Update(job ForecastJob) error
```

## ML Client

```go
Predict(req PredictRequest) (Forecast, error)
```

## Service

```go
UploadTransactions(contentType string, file io.Reader) error
GetTransactions(filters Filters) ([]Transaction, error)
UpdateTransaction(tx Transaction) error
DeleteTransaction(id int) error
CreateForecastJob() (int, error)
GetForecastJob(id int) (ForecastJob, error)
GetForecastJobs() ([]ForecastJob, error)
GetRecurringRules() ([]RecurringRule, error)
SaveRecurringRule(rule RecurringRule) error
DeleteRecurringRule(sourceID string) error
```

## Types

### Forecast

| Field            | Type        | Description                |
|------------------|-------------|----------------------------|
| points           | []Point     | точки прогноза             |
| predicted_balance| float       | итоговый баланс            |
| confidence       | float (0–1) | уверенность модели         |
| low_data_warning | bool        | true если confidence < 0.5 |

### ForecastJob

| Field      | Type      | Description                          |
|------------|-----------|--------------------------------------|
| id         | int       | идентификатор задачи                 |
| status     | string    | pending \| running \| done \| failed |
| result     | *Forecast | nil пока не готово                   |
| created_at | time.Time | время создания                       |
