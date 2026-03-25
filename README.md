Установка миграций 
```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

```

Миграция данных 
```bash
 migrate -path migrations/ -database "postgres://user:password@localhost:5432/price_tracker?sslmode=disable" up
1/u init (19.953792ms)
```

POST-запрос для добавления товаров в БД
```bash
curl -X POST http://localhost:8080/products \
     -H "Content-Type: application/json" \
     -d '{"name": "iPhone 15", "url": "://apple.com", "current_price": 799.00}'
```
GET-запрос для просмотра всех товаров
```bash
curl -X GET http://localhost:8080/products
```