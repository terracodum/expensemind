# Architecture

## Общая схема

Frontend → Go Backend → DB → ML

## Финансовая модель

* transactions — фактические операции
* recurring_rules — правила регулярных операций
* будущее вычисляется, не хранится

## Слои

* handler — HTTP
* service — бизнес-логика
* repository — БД
* parser — входные данные
* ml — внешний сервис

## Зависимости

handler → service
service → repository, parser, ml
repository не знает service
service не знает handler
