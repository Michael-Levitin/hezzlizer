# hezzlizer

Makefile
1. migrate - запуск миграции Postgres (один раз)
2. migrate_clickhouse - запуск миграции Clickhouse (один раз) (требуется заменить данные подключения.)
3. serverStart - запуск сервера
4. clientStart - запуск клиента


* В файле .env настройки для подключения к Postgres
* В папке config/config.go - можно изменить настройки по умолчанию 
* параметр newPriority был изменен на priority, для удобства чтения
* postman_collection.json - содержит образцы запросов для Postman