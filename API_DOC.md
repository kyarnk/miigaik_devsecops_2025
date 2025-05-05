# Task Manager API - Полная документация

## Оглавление
1. [Общая информация](#общая-информация)
2. [Аутентификация](#аутентификация)
3. [Пользователи](#пользователи)
   - [Регистрация](#регистрация)
   - [Авторизация](#авторизация)
   - [Выход](#выход)
   - [Удаление](#удаление-пользователя)
4. [Задачи](#задачи)
   - [Получить все](#получить-все-задачи)
   - [Получить по ID](#получить-задачу-по-id)
   - [Создать](#создать-задачу)
   - [Массовое создание](#массовое-создание-задач)
   - [Обновить](#обновить-задачу)
   - [Удалить](#удалить-задачу)
5. [Ошибки](#ошибки)
6. [Примеры](#примеры-использования)

## Общая информация

Базовый URL: `http://your-server-address:8080`

Формат данных: `application/json`

## Аутентификация

Большинство endpoints требуют JWT-токен в заголовке:
```http
Authorization: Bearer ваш_токен
```

Токен получается при авторизации и действует 24 часа.

## Пользователи

### Регистрация

**Endpoint:** `POST /user/register`

**Тело запроса:**
```json
{
    "username": "имя_пользователя",
    "password": "пароль"
}
```

**Успешный ответ (201):**
```json
{
    "id": 1,
    "username": "имя_пользователя"
}
```

### Авторизация

**Endpoint:** `POST /user/login`

**Тело запроса:**
```json
{
    "username": "имя_пользователя",
    "password": "пароль"
}
```

**Успешный ответ (200):**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Выход

**Endpoint:** `PUT /user/logout`

**Заголовки:**
```
Authorization: Bearer ваш_токен
```

**Успешный ответ (200):**
```
User logged out
```

### Удаление пользователя

**Endpoint:** `DELETE /user/delete`

**Заголовки:**
```
Authorization: Bearer ваш_токен
```

**Успешный ответ (200):**
```
User deleted
```

## Задачи

### Получить все задачи

**Endpoint:** `GET /tasks`

**Успешный ответ (200):**
```json
[
    {
        "id": 1,
        "title": "Купить продукты",
        "done": false,
        "user_id": 1
    }
]
```

### Получить задачу по ID

**Endpoint:** `GET /tasks/{id}`

**Пример:** `GET /tasks/1`

**Успешный ответ (200):**
```json
{
    "id": 1,
    "title": "Купить продукты",
    "done": false,
    "user_id": 1
}
```

### Создать задачу

**Endpoint:** `POST /tasks/create`

**Заголовки:**
```
Authorization: Bearer ваш_токен
Content-Type: application/json
```

**Тело запроса:**
```json
{
    "title": "Новая задача",
    "done": false,
    "user_id": 1
}
```

**Успешный ответ (201):**
```json
{
    "id": 2,
    "title": "Новая задача",
    "done": false,
    "user_id": 1
}
```

### Массовое создание задач

**Endpoint:** `POST /tasks/bulkupload`

**Заголовки:**
```
Authorization: Bearer ваш_токен
Content-Type: application/json
```

**Тело запроса:**
```json
[
    {
        "title": "Задача 1",
        "done": false,
        "user_id": 1
    },
    {
        "title": "Задача 2",
        "done": true,
        "user_id": 1
    }
]
```

**Успешный ответ (201):**
```json
[
    {
        "id": 3,
        "title": "Задача 1",
        "done": false,
        "user_id": 1
    },
    {
        "id": 4,
        "title": "Задача 2",
        "done": true,
        "user_id": 1
    }
]
```

### Обновить задачу

**Endpoint:** `PUT /tasks/update/{id}`

**Пример:** `PUT /tasks/update/1`

**Заголовки:**
```
Authorization: Bearer ваш_токен
Content-Type: application/json
```

**Тело запроса:**
```json
{
    "title": "Обновленная задача",
    "done": true
}
```

**Успешный ответ (200):**
```json
{
    "id": 1,
    "title": "Обновленная задача",
    "done": true,
    "user_id": 1
}
```

### Удалить задачу

**Endpoint:** `DELETE /tasks/delete/{id}`

**Пример:** `DELETE /tasks/delete/1`

**Заголовки:**
```
Authorization: Bearer ваш_токен
```

**Успешный ответ (200):**
```
Task with ID 1 has been deleted
```

## Ошибки

| Код | Статус         | Описание                     |
|-----|----------------|-----------------------------|
| 400 | Bad Request    | Неверные данные             |
| 401 | Unauthorized   | Неверный/отсутствующий токен|
| 404 | Not Found      | Задача/пользователь не найден|
| 500 | Server Error   | Ошибка сервера              |

## Примеры использования

### Пример с cURL

Авторизация:
```bash
curl -X POST http://localhost:8080/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test"}'
```

Создание задачи:
```bash
curl -X POST http://localhost:8080/tasks/create \
  -H "Authorization: Bearer ваш_токен" \
  -H "Content-Type: application/json" \
  -d '{"title":"Новая задача","done":false,"user_id":1}'
```

### Пример с JavaScript (Fetch API)

```javascript
async function createTask(taskData) {
  const response = await fetch('http://localhost:8080/tasks/create', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ваш_токен',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(taskData)
  });
  return await response.json();
}
```
