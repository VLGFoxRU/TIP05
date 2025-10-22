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
pz5-db/
├── .env
├── db.go
├── go.mod
├── go.sum
├── main.go
└── repository.go
```

# Скриншоты

Создание таблицы в psql:

<img width="1919" height="1030" alt="image" src="https://github.com/user-attachments/assets/52edf383-7716-4d22-8ad2-81799c95ba13" />

Успешный вывод:

<img width="507" height="200" alt="image" src="https://github.com/user-attachments/assets/ce43ad2c-a218-48ca-99ec-9b4c2e7fff81" />

SELECT * FROM tasks:

<img width="1920" height="1027" alt="image" src="https://github.com/user-attachments/assets/aab55012-be2b-40a1-b6ba-cb278d16dc70" />

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
