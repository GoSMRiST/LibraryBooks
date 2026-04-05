# 📚 LibraryBooks Service

Сервис для управления библиотекой книг, реализованный на Go.
Поддерживает **REST API** и **gRPC** для взаимодействия с другими сервисами.

Проект построен с использованием принципов **чистой архитектуры (Clean Architecture)** с разделением на слои: transport, services, repository и core.

---

## 🚀 Возможности

* 📖 Создание и получение книг
* 🔍 Работа с данными через REST и gRPC
* 🧩 Чёткое разделение слоёв (architecture-friendly)
* 🔐 Middleware для авторизации
* 🗄️ Работа с БД через repository слой
* ⚙️ Конфигурация через env-файлы
* 📦 Поддержка миграций базы данных

---

## 🛠️ Технологии

* **Go**
* **gRPC**
* **REST API**
* **PostgreSQL**
* **SQL migrations**

---

## 📂 Структура проекта

```bash
.
├── cmd/app                # Точка входа (main.go)
├── internal/
│   ├── app/               # Инициализация приложения
│   │   ├── grpc/
│   │   └── rest/
│   ├── config/            # Конфигурация
│   ├── core/              # Модели и бизнес-логика
│   ├── middleware/        # Middleware (auth)
│   ├── repository/        # Работа с БД
│   ├── services/          # Сервисный слой
│   └── transport/         # REST / gRPC обработчики
├── migrations/            # SQL миграции
├── go.mod
└── go.sum
```

---

## ⚙️ Конфигурация

Создайте `.env` файл:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=library
```

---

## 🗄️ База данных

Используется PostgreSQL.

Примените миграции из папки:

```bash
migrations/
```

(можно через любой инструмент или вручную)

---

## ▶️ Запуск

```bash
go run cmd/app/main.go
```

---

## 🌐 API

### REST

Пример эндпоинтов:

* `GET /books` — получить список книг
* `POST /books` — создать книгу

---

### gRPC

* Реализация находится в `internal/transport/grpc`
* Контракты описаны в `.proto` файлах

---

## 🔐 Аутентификация

Реализована через middleware:

```
internal/middleware/authmiddle.go
```

Используются токены для защиты эндпоинтов.

---

## 🧪 Пример запроса

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Book","author":"Author"}'
```

---

## 👤 Автор

* GitHub: https://github.com/GoSMRiST
