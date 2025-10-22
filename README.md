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

<img width="1918" height="1026" alt="image" src="https://github.com/user-attachments/assets/a78dfd1d-5750-4b86-a889-a8d2309df69c" />

# Фрагменты кода
```
db.go

package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// настройки пула — достаточно для локалки
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// проверка соединения с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	log.Println("Connected to PostgreSQL")
	return db, nil
}
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
