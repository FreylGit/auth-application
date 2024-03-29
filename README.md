# Сервис Авторизации

Этот проект создан для облегчения процесса авторизации в других проектах, предоставляя стандартный сервис авторизации с методами Registration, Login и Refresh.

Сервис работает по протоколу gRPC.

Проект с моделями proto: [https://github.com/FreylGit/protoModel](https://github.com/FreylGit/protoModel)

## Запуск в Docker

```bash
docker-compose up --build
```

### Описание схемы базы данных

Для данного проекта используется база данных PostgreSQL с следующей структурой:

---

#### **Таблица `users`**

- **Поля:**
  - `id` (Serial, Первичный ключ): Уникальный идентификатор пользователя.
  - `email` (Varchar(256), Уникальный): Email адрес пользователя (уникальное ограничение).
  - `password_hash` (Bytea, Не Null): Хэш пароля пользователя.
  - `name` (Varying(256)): Имя пользователя.
  - `create_date` (Date, По умолчанию: Текущая дата): Дата создания записи пользователя.

#### **Таблица `refresh_tokens`**

- **Поля:**
  - `id` (Serial, Первичный ключ): Уникальный идентификатор обновляющего токена.
  - `user_id` (Integer, Не Null): Внешний ключ, ссылается на столбец `id` в таблице `users`.
  - `token` (Varchar(256), Не Null): Обновляющий токен, связанный с пользователем.
  - `create_date` (Date, По умолчанию: Текущая дата): Дата создания обновляющего токена.
  - `expired_date` (Date): Дата истечения срока действия обновляющего токена.

---

**Описание:**
- Схема базы данных состоит из двух таблиц: `users` и `refresh_tokens`.
- В таблице `users` хранятся данные о пользователях, такие как их идентификатор, электронная почта, хэш пароля, имя и дата создания записи.
- Таблица `refresh_tokens` предназначена для хранения обновляющих токенов, связанных с пользователями. Каждый токен имеет уникальный идентификатор, ссылку на соответствующего пользователя, сам токен, дату его создания и дату истечения срока действия.

Эта структура базы данных обеспечивает безопасную аутентификацию пользователей и управление их токенами в рамках проекта.

