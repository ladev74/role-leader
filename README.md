# Role-Leader Service Api
Микросервис для управления встречами и обратной связью между лидерами и пользователями.
Сервис разработан для командного проектного этапа, на курсе Yandex Lyceum "Веб-разработка на Go | Специализации Яндекс Лицея | Весна 24/25"

## API Endpoints

### 1. Создание обратной связи для встречи

**Endpoint:**  
`POST: /api/create-feedback`

**Request Body (JSON):**

```json
{
  "call_id": "call_1",
  "message": "Discussed the key points"
}
```

**Response (Success):**

```json
{
  "{}"
}
```

---

### 2. Получение информации о встрече

**Endpoint:**  
`GET: "/api/get-call/{call_id}"`

**Пример запроса:**

```
GET /api/get-call/call_1
```

**Response:**

```json
{
  "call": {
    "call_id": "call_1",
    "user_id": "user_1",
    "leader_id": "leader_1",
    "title": "Planning",
    "start_time": "12:30:00",
    "status": "Completed",
    "feedback": "Discussed the key points"
  }
}
```

---

### 3. Получение всех встреч лидера

**Endpoint:**  
`GET: /api/leader-calls/{leader_id}`

**Пример запроса:**

```
GET /api/leader-calls/leader_1
```

**Response:**

```json
{
  "calls": [
    {
      "call_id": "call_1",
      "user_id": "user_1",
      "leader_id": "leader_1",
      "title": "Planning",
      "start_time": "12:30:00",
      "status": "Completed",
      "feedback": "Discussed the key points"
    },
    {
      "call_id": "call_2",
      "user_id": "user_2",
      "leader_id": "leader_1",
      "title": "title4",
      "start_time": "04:04:04",
      "status": "status4",
      "feedback": "feedback4"
    }
  ]
}
```

---

## Примеры cURL

### Создание обратной связи

```bash
curl -X POST http://localhost:8080/api/create-feedback \
-H "Content-Type: application/json" \
-d '{
  "call_id": "call_1",
  "message": "Discussed the key points"
}'
```

### Получение встречи

```bash
curl -X GET http://localhost:8080/api/get-call/call_1 
```

### Получение встреч лидера

```bash
curl -X GET http://localhost:8080/api/leader-calls/leader_1 
```

---

## Структуры данных

### Объект Meeting

| Поле       | Тип    | Описание                          |
|------------|--------|-----------------------------------|
| call_id    | string | Уникальный ID встречи             |
| user_id    | string | ID пользователя                   |
| leader_id  | string | ID лидера                         |
| title      | string | Название встречи                  |
| status     | string | Статус (planned/active/completed) |
| feedback   | string | Комментарий обратной связи        |
| start_time | string | Время начала (ISO 8601)           |

### Использованные технологии

| Название                 | Применение                          |
|--------------------------|-------------------------------------|
| PostgreSQl               | Хранение данных о звонках           |
| GRPC                     |                                     |
| REST API                 |                                     |
| Docker  + Docker Compose | Контейнеризация и изоляция сервисов |