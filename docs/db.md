# Database

## transactions

* id
* amount (float, + доход / - расход)
* description
* category
* date

## recurring_rules

* id
* source_id
* type (income | expense)
* amount (положительный)
* day
* start_date
* label

Уникальность:
(source_id, start_date)

## forecast_jobs

* id
* status
* result
* created_at
