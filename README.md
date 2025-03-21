# 🎵 Music API 🎶

Этот проект представляет собой **RESTful API** для управления библиотекой песен.  
API позволяет **добавлять, обновлять, удалять** и **получать информацию о песнях**,  
а также **получать текст песни** с пагинацией по куплетам.

## 🚀 Функциональность

- 📜 **Получение списка песен** с фильтрацией по группе и названию песни, а также пагинацией.
- 🎤 **Получение текста песни** с пагинацией по куплетам.
- ➕ **Добавление новой песни** в формате JSON.
- ✏️ **Обновление данных песни**.
- ❌ **Удаление песни** по её ID.

## 🛠️ Технологии

- 💻 **Язык программирования**: Go
- 🌐 **Фреймворк**: Gin
- 🗄️ **База данных**: PostgreSQL
- 🔄 **Миграции**: Golang Migrate
- 📜 **Логирование**: Logrus
- 📖 **Документация API**: Swagger (Swaggo)

## 📥 Установка и настройка

### 1️⃣ Клонирование репозитория

```bash
git clone https://github.com/justcr1si/go-api.git
cd go-api
```

### 2️⃣ Установка зависимостей
Убедитесь, что у вас установлен **Go (1.16+)**, **PostgreSQL**, **pgAdmin (опционально)**.
Зависимости устанавливаются автоматически командой:
```bash
go mod tidy
```

### 3️⃣ Настройка базы данных
Заполните файл **.env** данными для подключения к базе данных **(можно воспользоваться pgAdmin)**:
```plaintext
DATABASE_URL=postgres://postgres:cdh6ed39cz@localhost:5432?sslmode=disable
API_URL=http://api.example.com/info
```
🔹 **Таблицы создавать не нужно, так как уже существует таблица songs.**

## ▶️ Запуск и использование API
Запустите сервер:
```bash
go run main.go
```

## 📖 Swagger UI
Документация API доступна по адресу:

🔗 http://localhost:8080/swagger/index.html

Либо используйте **curl-запросы / Postman** для тестирования API.