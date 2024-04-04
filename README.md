# APIGateway (AI Marketplace)

***тестовое задание***

Микросервис, представляющий REST API ручки для связи с [Statistic Service](https://github.com/shamank/ai-marketplace-stats-service)

Полный проект: [AI Marketplace](https://github.com/shamank/ai-marketplace)

### Использованные библиотеки:
- **gin** - для работы с HTTP
- **cleanenv** - для работы с конфигом

### Запуск микросервиса

**Для старта микросервиса**:

```sh
make run
```
или
```sh
go run ./cmd/app/main.go --cfg=./configs/prod.yaml
```

### REST API


**Запрос на получение статистики:**
```http request
GET /api/calls
```
Доступные query-параметры:
- `user_uid` - фильтрация по конкретному пользователю
- `service_uid` - фильтрация по конкретному сервису
- `order` - (asc или desc) - сортировка по дате времени запроса
- `page_number` - номер страницы (для пагинации)
- `page_size` - размер данных (для пагинации)

Пример выходных данных:
```json
{
  "calls": [
    {
      "user_uid": "123e4567-e89b-12d3-a456-426655440000",
      "service_uid": "223e4567-e89b-12d3-a456-426655440001",
      "count": 2,
      "full_amount": 400
    },
    {
      "user_uid": "123e4567-e89b-12d3-a456-426655440002",
      "service_uid": "223e4567-e89b-12d3-a456-426655440001",
      "count": 1,
      "full_amount": 200
    },
    {
      "user_uid": "123e4567-e89b-12d3-a456-426655440000",
      "service_uid": "223e4567-e89b-12d3-a456-426655440000",
      "count": 1,
      "full_amount": 100
    }
  ]
}
```
--- 
**Запрос на создание сервиса**
```http request
POST /api/service
```
Пример тела запроса:
```json
{
    "title": "my_service",
    "description": "my description", // не обязательное поле
    "price": 0.003
}
```
Пример ответа:
```json
{
    "uid": "caaca9c0-a27d-40b5-b661-32dbe24b2d7d" // UID нового сервиса
}
```

**Запрос на использование сервиса**
```http request
POST /api/call
```
Пример тела запроса:
```json
{
    "user_uid": "123e4567-e89b-12d3-a456-426655440002",
    "service_uid": "e22a027d-44b8-4848-b972-b96c9e7c4690"
}
```
Пример тела ответа:
```json
{
    "message": "OK"
}
```
--- 


### Структура проекта

```
├───cmd
│   └───app
├───configs
└───internal
    ├───app
    ├───clients // клиенты, к которым кидаются запросы
    │   └───stats-service
    │       └───grpc
    ├───config
    ├───delivery
    │   └───http // описанные REST API хэндлеры
    └───domain
        └───models
```

*Примечание*: решил не делать слой сервиса, так как логика тут нет смысла опускаться ниже на уровень (если не выходить за рамки задания)
