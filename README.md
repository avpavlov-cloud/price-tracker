Установка миграций 
``bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

``

Миграция данных 
``bash
 migrate -path migrations/ -database "postgres://user:password@localhost:5432/price_tracker?sslmode=disable" up
1/u init (19.953792ms)
``