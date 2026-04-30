# AceClub
## О проекте:
Ace Club - cервис для управления теннисным клубом. 


## Технологический стек

| Компонент | Технология |
|----------|-----------|
| Язык | Go 1.23+ |
| HTTP роутер | Стандартный net/http |
| База данных | PostgreSQL 16 |
| Драйвер БД | database/sql |
| Миграции | golang-migrate/migrate CLI |
| Хэширование паролей | golang.org/x/crypto/bcrypt |
| Логирование | log/slog (структурированное) |
| Тестирование | testing + testcontainers |
| Конфигурация | .env через godotenv |

## Слои приложения

| Слой | Назначение | Пример |
|-----|-----------|--------|
| Transport | Обработка HTTP запросов, маршрутизация, middleware | UserHandler, ErrorAdapter |
| Service | Бизнес-логика, валидация, хэширование | UserService.RegisterUser |
| Repository | Доступ к БД, SQL запросы | UserRepository.Create |
| Domain | Модели данных и бизнес-ошибки | User, ValidationError |

**Поток запроса**
HTTP Request → LoggingMiddleware → ErrorAdapter → Handler → Service → Repository → DB
↓
ошибка → ErrorAdapter → HTTP Response


## API Endpoints

### GET /health

Проверка работоспособности сервиса.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2026-04-30T10:00:00Z"
}
```

## POST /api/users
Регистрация нового пользователя

### Request:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Validation rules:
- Email: не пустой, корректный формат (RFC 5322)
- Password: длина от 8 до 16 символов

### Success response (201 Created):
```json
{
  "id": 1,
  "email": "user@example.com"
}
```

### Error responses:

| Сценарий | Статус | Пример ответа |
|-----|-----------|--------|
| Невалидный email | 422 | {"error":"invalid email format: mail"} |
| Короткий пароль | 422 | {"error":"password length must be >= 8"} |
| Длинный пароль | 422 | {"error":"password length must be <= 16"} |
| Email уже существует | 409 | {"error":"user with id user@example.com already exists"} |
| Некорректный json | 422 | {"error":"validation failed for field: body"} |
| Внетренняя ошибка | 500 | {"error":"internal server error"} |

## Обработка ошибок
| Тип | Использование | HTTP статус |
|-----|-----------|--------|
| ValidationError | Невалидные входные данные | 422 | 
| ConflictError | Ресурс уже существует | 409 |
| NotFoundError | Ресурс не найден | 404 |

## ErrorAdapter

Централизованная обработки всех ошибок, возвращаемых из handlers:
- Определяет HTTP статус через statusCodeFromError()

- Логирует ошибку с уровнем ERROR

- Возвращает JSON с полем error