# ЭФМО-01-25 Буров М.А. ПР5

# Описание проекта
Подключение к PostgreSQL через database/sql. Выполнение простых запросов (INSERT, SELECT)

# Требования к проекту
* Go 1.25+
* PostgreSQL 10+
* Git

# Версия Go
<img width="340" height="55" alt="image" src="https://github.com/user-attachments/assets/e7fec853-899d-442e-9d95-94d494a5cb46" />

# Версия PostgreSQL
<img width="276" height="55" alt="image" src="https://github.com/user-attachments/assets/bca8a5ed-eb3b-4b1c-9eee-16fb8ec409e8" />

# Цели:
- Установить и настроить PostgreSQL локально.
- Подключиться к БД из Go с помощью database/sql и драйвера PostgreSQL.
- Выполнить параметризованные запросы INSERT и SELECT.
- Корректно работать с context, пулом соединений и обработкой ошибок.

# Структура проекта
Дерево структуры проекта: 
```
pz4-todo/
├── internal/
│   └── task/
│       ├── handler.go
│       ├── model.go
│       └── repo.go
├── pkg/
│   └── middleware/
│       ├── cors.go
│       └── logger.go
├── go.mod
├── go.sum
└── main.go
```

# Скриншоты

Создание таблицы в psql:



Успешный вывод:


SELECT * FROM tasks:

# Фрагменты кода
```
db.go


```

```
repository.go


```

```
main.go


```

# Краткие ответы

- Что такое пул соединений *sql.DB и зачем его настраивать?

- Почему используем плейсхолдеры $1, $2?

- Чем Query, QueryRow и Exec отличаются?

# Обоснование транзакций и настроек пула
