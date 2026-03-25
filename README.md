# Price Tracker Service (Go + PostgreSQL)

Простое и эффективное приложение для мониторинга цен на товары. Система автоматически отслеживает изменения цен в фоновом режиме, сохраняет историю изменений и предоставляет API для доступа к данным.

---

## 🚀 Особенности архитектуры

- **Clean Architecture**  
  Разделение на:
  - модели
  - репозитории (работа с БД)
  - сервисы (бизнес-логика)

- **Concurrency**  
  Фоновый воркер на горутинах использует `time.Ticker` для периодического обновления данных.

- **Data Integrity**  
  Использование транзакций PostgreSQL для атомарного:
  - обновления текущей цены  
  - записи в таблицу истории

- **Graceful Shutdown**  
  Безопасная остановка сервера и воркеров с завершением текущих операций при получении системных сигналов (`Ctrl+C`).

---

## 🛠 Стек технологий

- **Язык:** Go 1.21+
- **База данных:** PostgreSQL 15+
- **HTTP Framework:** Echo v4
- **DB Driver:** pgx (v5)
- **Миграции:** golang-migrate
- **Работа с деньгами:** shopspring/decimal (защита от ошибок округления `float`)

---

## ⚙️ Установка и запуск

### 1. Подготовка базы данных

Запустите PostgreSQL (например, через Docker):

```bash
docker run --name price-db \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_USER=user \
  -e POSTGRES_DB=price_tracker \
  -p 5432:5432 \
  -d postgres:15-alpine
````

---

### 2. Установка инструмента миграций

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

### 3. Применение миграций

```bash
migrate -path migrations/ \
  -database "postgres://user:password@localhost:5432/price_tracker?sslmode=disable" \
  up
```

---

### 4. Запуск приложения

```bash
go run cmd/server/main.go
```

---

## 📡 API Эндпоинты

### ➕ Добавление нового товара

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "url": "https://apple.com",
    "current_price": 799.00
  }'
```

---

### 📋 Список всех товаров

```bash
curl -X GET http://localhost:8080/products
```

---

### 📈 История изменений цены товара

```bash
curl http://localhost:8080/products/1/history
```

---

## 🏗 Структура проекта

```
/cmd/server           # Входная точка приложения + Graceful Shutdown
/internal/models      # Структуры данных (Product, PriceHistory)
/internal/repository  # Слой доступа к БД (PostgreSQL + транзакции)
/internal/service     # Бизнес-логика + фоновый воркер
/migrations           # SQL миграции схемы БД
```

---

## 📌 Примечания

* Используется безопасная работа с денежными значениями через `decimal`
* Поддерживается масштабирование за счёт разделения слоёв
* Готово к расширению (например: уведомления, парсинг сайтов, очереди)

