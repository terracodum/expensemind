# Errors

## Codes

| Code                    | Когда                                   |
|-------------------------|-----------------------------------------|
| INTERNAL_ERROR          | непредвиденная ошибка                   |
| VALIDATION_ERROR        | невалидные входные данные               |
| NOT_FOUND               | ресурс не найден                        |
| PARSE_ERROR             | ошибка парсинга CSV или PDF             |
| ML_SERVICE_UNAVAILABLE  | ML сервис недоступен                    |
| ML_RESPONSE_INVALID     | ML вернул некорректный ответ            |
| DB_ERROR                | ошибка базы данных                      |

## Format

```json
{
  "error": {
    "code": "PARSE_ERROR",
    "message": "no transactions found in pdf"
  }
}
```
