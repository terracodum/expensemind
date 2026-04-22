# Backlog

Вещи обсуждённые но не реализованные.

---

## Backend

### service/ — forecastWorker (дописать)

Файл: `backend/internal/service/service.go`, метод `forecastWorker(id int)`

Текущее состояние: загружает транзакции и правила, дальше пусто.

Что дописать:

1. **Статус running** — сразу после старта:
   ```go
   s.forecastRepo.Update(domain.ForecastJob{ID: id, Status: "running"})
   ```

2. **Агрегация транзакций в `[]ml.TimePoint`** — группировать по дням:
   - Фильтровать `category != "transfer"` (не влияют на баланс)
   - Каждый день → одна `ml.TimePoint` с полями:
     - `T` — порядковый номер дня (1, 2, 3...)
     - `Balance` — баланс на конец дня (накопительная сумма amount)
     - `DayOfWeek` — день недели (1=пн, 7=вс)
     - `IsWeekend` — суббота или воскресенье
     - `FoodTotal` — сумма транзакций с category="food"
     - `TransportTotal` — сумма транзакций с category="transport"
     - `EntertainmentTotal` — сумма транзакций с category="entertainment"
     - `AvgTransactionSize` — средний чек за день
     - `TransactionCount` — количество транзакций за день

3. **Генерация `income_events` из recurring_rules** — каждое правило → `ml.IncomeEvent`:
   - `T` — номер дня относительно последней транзакции
   - `Amount` — rule.Amount
   - `Label` — rule.Label

4. **Сборка `ml.PredictRequest`**:
   ```go
   req := ml.PredictRequest{
       Timeseries: timepoints,
       Horizon:    30,
       Features: ml.Features{
           AvgDailyExpense: avgDailyExpense,
           IncomeEvents:    incomeEvents,
       },
   }
   ```

5. **Вызов ML и сохранение результата**:
   ```go
   forecast, err := s.ml.Predict(req)
   // fail(err); return если ошибка
   forecast.LowDataWarning = forecast.Confidence < 0.5
   s.forecastRepo.Update(domain.ForecastJob{ID: id, Status: "done", Result: &forecast})
   ```

---

### service/ — фильтрация переводов
- Категория `transfer` зарезервирована — переводы между своими счетами
- При подготовке данных для ML фильтровать `category != "transfer"`
- Пользователь сам размечает переводы в UI
- В UI показать подсказку: "переводы между своими счетами, которые не меняют общий баланс, стоит пометить как transfer"

