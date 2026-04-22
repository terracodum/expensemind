# Contracts

## Parser

Parse(file io.Reader) ([]Transaction, error)

## Repository

### Transactions

Save(tx Transaction) error
SaveAll(txs []Transaction) error
FindAll(filters Filters) ([]Transaction, error)

### Recurring

SaveRule(rule RecurringRule) error
FindAllRules() ([]RecurringRule, error)
DeleteRule(id string) error

## Types

### Forecast

| Field            | Type        | Description               |
|------------------|-------------|---------------------------|
| confidence       | float (0–1) | уверенность модели        |
| low_data_warning | bool        | true если confidence < 0.5 |

## ML Client

Predict(req PredictRequest) (Forecast, error)

## Service

UploadPDF(file io.Reader) (int, error)
GetTransactions(filters Filters) ([]Transaction, error)
GetForecast(horizon int) (Forecast, error)
