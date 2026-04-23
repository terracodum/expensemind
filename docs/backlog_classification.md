# [ML] Персональная классификация транзакций

## Суть
Добавить в Python ML сервис модуль автоматической категоризации транзакций на основе истории разметки пользователя.

**Стек:** TF-IDF + MultinomialNB (sklearn)

---

## Схема БД

```sql
-- В таблице transactions добавить:
category_id          INTEGER     -- null пока не размечено
is_auto_classified   BOOLEAN     DEFAULT false

-- В таблице settings:
classifier_ready         BOOLEAN   DEFAULT false
classifier_trained_at    TIMESTAMP
```

---

## Флоу

### Фаза 1 — холодный старт (модели нет)
1. Юзер загружает CSV
2. Транзакции пишутся в БД с `category_id = null`, `is_auto_classified = false`
3. Юзер вручную проставляет категории в UI
4. БД обновляется

### Фаза 2 — обучение
1. Go считает количество транзакций где `category_id IS NOT NULL AND is_auto_classified = false`
2. Как только их ≥ 200 — показать баннер "достаточно данных, можно обучить модель"
3. Юзер жмёт кнопку "Обучить"
4. Go забирает все размеченные вручную транзакции из БД
5. Отправляет `POST /classify/train` в ML сервис
6. ML сервис обучает модель, сохраняет `model.pkl`
7. Go обновляет `classifier_ready = true`, `classifier_trained_at = now()`

### Фаза 3 — автоматика
1. Юзер загружает новый CSV
2. Go видит `classifier_ready = true`
3. Каждую транзакцию отправляет в `POST /classify/predict`
4. ML сервис возвращает категорию
5. Go пишет в БД с `is_auto_classified = true`
6. В UI авто-размеченные транзакции помечены визуально (иконка / серый цвет)
7. Если юзер поправляет авто-размеченную → `is_auto_classified = false` → попадает в следующую обучающую выборку

---

## ML сервис — новые эндпоинты

```
POST /classify/train
  body: [{ description: string, category: string }]
  response: { accuracy: float, samples: int }

POST /classify/predict
  body: { description: string }
  response: { category: string, confidence: float }
```

## ML сервис — новый модуль

```
ml_service/
  forecasting.py      # уже есть, Prophet
  classification.py   # новый, sklearn
  main.py             # добавить новые роуты
```

Внутри `classification.py`:
- `train(transactions)` — обучить TF-IDF + MultinomialNB, сохранить `model.pkl`
- `predict(description)` — загрузить модель из памяти, вернуть категорию
- Модель держать в памяти после первой загрузки, не читать `pkl` на каждый запрос

---

## Почему не словарь
Модель обобщает варианты написания одного мерчанта без ручного поддержания списка. "ПЯТЁРОЧКА 1234 МСК" и "ПЯТЕРОЧКА ЭКСПРЕСС СПБ" — одна категория без явного прописывания каждого варианта.

---

## Приоритет
Добавить после завершения `parser/pdf.go` (TBankParser).
