# hezzlizer

запуск ClickHouse и Nats серверов :
* sudo clickhouse start
* nats-server

Makefile
1. migrate - запуск миграции Postgres (один раз)
2. migrate_clickhouse - запуск миграции Clickhouse (один раз)
3. serverStart - запуск сервера
4. clientCHStart - запуск клиента ClickHouse (получает данные через NATS и делает запрос в CH)
5. clientHTTPStart - запуск демо клиента, который симулирует запросы к серверу

* В файле .env настройки для подключения к Postgres
* В папке config/config.go - можно изменить настройки по умолчанию 
* параметр newPriority был изменен на priority, для удобства чтения
* postman_collection.json - содержит образцы запросов для Postman
* Миграция и дальнейшая работа с ClickHouse - на БД "default", без пароля