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

## ML Client

Predict(req PredictRequest) (Forecast, error)

## Service

UploadPDF(file io.Reader) (int, error)
GetTransactions(filters Filters) ([]Transaction, error)
GetForecast(horizon int) (Forecast, error)
