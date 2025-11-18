# CommentTree - Древовидные комментарии

Сервис для работы с древовидными комментариями с поддержкой неограниченной вложенности, поиска, сортировки и постраничного вывода.

### HTTP API

- POST /comments - создание родительского и дочернего комментария
- GET /comments?parent={id} - получение комментария и всех вложенных
- GET /comments?search={query} - полнотекстовый поиск по комментариям
- GET /comments?page={n}&page_size={m}&order={asc|desc} - постраничная навигация и сортировка
- DELETE /comments/{id} - удаление комментария и всех вложенных под ним

## Установка и запуск проекта

### 1. Клонирование репозитория
```bash
git clone https://github.com/kstsm/wb-comment-tree
```

### 2. Настройка переменных окружения
Создайте `.env` файл, скопировав в него значения из `env.example`:
```bash
cp env.example .env
```

### 3. Запуск зависимостей (Docker)
```bash
make docker-up
```

### 4. Миграция базы данных
```bash
make goose-up
```

### 5. Запуск сервиса
```bash
make run
```

Сервис будет доступен по адресу: http://localhost:8080



# API запросы

## POST /comments - Создание комментария

### Создание корневого комментария (без родителя)

**Метод:** `POST`  
**URL:** `http://localhost:8080/comments`  
**Body (JSON):**
```json
{
  "comment": "Первый комментарий"
}
```

**Ожидаемый ответ (201 Created):**
```json
{
  "result": {
    "id": 1,
    "comment": "Первый комментарий",
    "created_at": "2025-11-18T11:49:15+06:00"
  }
}
```

### Создание ответа на комментарий (с родителем)

**Метод:** `POST`  
**URL:** `http://localhost:8080/comments`  
**Body (JSON):**
```json
{
  "parent_id": 1,
  "content": "Ответ на первый комментарий"
}
```

**Ожидаемый ответ (201 Created):**
```json
{
  "result": {
    "id": 1,
    "parent_id": 2,
    "comment": "Ответ на первый комментарий",
    "created_at": "2025-11-18T11:49:15+06:00"
  }
}
```
### Ошибки:

**Пустой контент (422 Unprocessable Entity):**
```json
{
  "comment": ""
}
```
Ответ:
```json
{
  "error": "content cannot be empty"
}
```

**Несуществующий родитель (404 Not Found):**
```json
{
  "parent_id": 999,
  "comment": "Комментарий с несуществующим родителем"
}
```
Ответ:
```json
{
  "error": "parent not found"
}
```

---

## GET /comments - Получение комментариев

### Пагинация и сортировка

**Метод:** `GET`  
**URL:** `http://localhost:8080/comments?page=1&page_size=10&order=desc`

**Query параметры:**
- `page` (опционально, по умолчанию 1) - номер страницы
- `page_size` (опционально, по умолчанию 10) - размер страницы
- `order` (опционально, по умолчанию `desc`) - порядок сортировки (`asc`, `desc`)

**Ожидаемый ответ (200 OK):**
```json
{
  "result": {
    "comments": [
      {
        "id": 1,
        "comment": "Это первый комментарий в системе",
        "created_at": "2025-11-18T11:45:45+06:00"
      }
    ],
    "page": 1,
    "page_size": 10,
    "total": 1
  }
}
```
### Полнотекстовый поиск комментариев

**Метод:** `GET`  
**URL:** `http://localhost:8080/comments?search=первый`

**Query параметры:**
- `search` - поисковый запрос
**Ожидаемый ответ (200 OK):**
```json
{
  "result": {
    "comments": [
      {
        "id": 1,
        "comment": "Это первый комментарий в системе",
        "created_at": "2025-11-18T11:45:45+06:00"
      }
    ],
    "page": 1,
    "page_size": 10,
    "total": 1
  }
}
```
---

## DELETE /comments/{id} - Удаление комментария

### Удаление комментария (с каскадным удалением дочерних)

**Метод:** `DELETE`  
**URL:** `http://localhost:8080/comments/1`

**Ожидаемый ответ (200 OK):**
```json
{
  "result": {
    "message": "comment deleted successfully"
  }
}
```

### Ошибки:

**Несуществующий комментарий (404 Not Found):**
```
DELETE http://localhost:8080/comments/999
```
Ответ:
```json
{
  "error": "comment not found"
}
```

**Некорректный ID (400 Bad Request):**
```
DELETE http://localhost:8080/comments/abc
```
Ответ:
```json
{
  "error": "invalid comment id"
}
```

---








