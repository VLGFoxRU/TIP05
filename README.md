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

package main

import (
	"context"
	"database/sql"
	"time"
)

// Task — модель для сканирования результатов SELECT
type Task struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo { return &Repo{DB: db} }

// CreateTask — параметризованный INSERT с возвратом id
func (r *Repo) CreateTask(ctx context.Context, title string) (int, error) {
	var id int
	const q = `INSERT INTO tasks (title) VALUES ($1) RETURNING id;`
	err := r.DB.QueryRowContext(ctx, q, title).Scan(&id)
	return id, err
}

// ListTasks — базовый SELECT всех задач (демо для занятия)
func (r *Repo) ListTasks(ctx context.Context) ([]Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks ORDER BY id;`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// Фильтр задач по полю done
func (r *Repo) ListDone(ctx context.Context, done bool) ([]Task, error) {
    const q = `SELECT id, title, done, created_at FROM tasks WHERE done = $1 ORDER BY id;`
    
    rows, err := r.DB.QueryContext(ctx, q, done)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var out []Task
    for rows.Next() {
        var t Task
        if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
            return nil, err
        }
        out = append(out, t)
    }
    
    return out, rows.Err()
}

// Поиск задачи по ID
func (r *Repo) FindByID(ctx context.Context, id int) (*Task, error) {
    const q = `SELECT id, title, done, created_at FROM tasks WHERE id = $1;`
    
    var t Task
    err := r.DB.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
    
    if err == sql.ErrNoRows {
        return nil, nil // задача не найдена
    }
    if err != nil {
        return nil, err // другая ошибка
    }
    
    return &t, nil
}

// Массовая вставка через транзакцию
func (r *Repo) CreateMany(ctx context.Context, titles []string) error {
    tx, err := r.DB.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    for _, title := range titles {
        const q = `INSERT INTO tasks (title) VALUES ($1);`
        _, err := tx.ExecContext(ctx, q, title)
        if err != nil {
            return err // откат произойдёт в defer
        }
    }
    
    if err = tx.Commit(); err != nil {
        return err
    }
    
    return nil
}
```

```
main.go

// 4) Прочитаем список задач с done == false
	tasks, err = repo.ListDone(ctxList, false)
	if err != nil {
		log.Fatalf("ListDone error: %v", err)
	}

	// 5) Напечатаем
	fmt.Println("=== Tasks with done=false flag ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	// 6) Найдем задачу по ID
	ctxFind, cancelFind := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFind()

	task, err := repo.FindByID(ctxFind, 1)
	if err != nil {
		log.Fatalf("FindByID error: %v", err)
	}

	// 7) Напечатаем
	fmt.Println("=== Task with ID=1 ===")
	fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
		task.ID, task.Title, task.Done, task.CreatedAt.Format(time.RFC3339))

	// 8) Выполним массовую вставку через транзакцию
	ctxMany, cancelMany := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelMany()

	testTitles := []string{"Изучить Go", "Купить корицы", "Начать ПЗ №6"}
	if err := repo.CreateMany(ctxMany, testTitles); err != nil {
		log.Printf("CreateMany error: %v", err)
	} else {
		fmt.Println("\n Все задачи вставлены через транзакцию")
	}

	// 9) Выполним логирование статистики пула
	fmt.Println("=== Pool stats ===")
	stats := db.Stats()
	fmt.Printf("Open=%d, InUse=%d, Idle=%d, WaitCount=%d",
		stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
```

# Краткие ответы

- Что такое пул соединений *sql.DB и зачем его настраивать?

Пул соединений -- это управляемый набор соединений к базе данных. Объект *sql.DB в Go не предтавляет одно соединение, а управляет пулом нескольких соединений, которые переиспользуются между запросами.
Его настройка позволяет повысить производительность, ограничить потребление ресурсов, начать обрабатывать больше параллельных запросов.

- Почему используем плейсхолдеры $1, $2?

Для защиты от SQL-инъекций.

- Чем Query, QueryRow и Exec отличаются?

Query - возвращает несколько строк.
QueryRow - возвращает одну строку. 
Exec - возвращает количество затронутых строк (используется для команд без возвращения данных, например INSERT, DELETE)

# Обоснование транзакций и настроек пула

Транзакция -- это атомарная последовательность SQL-операций: либо все выполняются и фиксируются, либо все откатываются при ошибке.

Вместо создания нового соединения для каждого запроса (медленно), пул переиспользует существующие соединения. Это экономит ресурсы и улучшает пропускную способность.

Проблема без настройки:
* По умолчанию SetMaxOpenConns = неограничено → система может исчерпать файловые дескрипторы
* По умолчанию SetMaxIdleConns = 2 → мало соединений в резерве → медленно переиспользуются

Выбранные настройки:
* SetMaxOpenConns(10) - Максимум открытых соединений одновременно.
* SetMaxIdleConns(5) - Максимум соединений в режиме ожидания (не используются, но не закрыты).
* SetConnMaxLifetime(30 * time.Minute) - Максимальный возраст соединения перед закрытием.

Выбранные настройки пула оптимальны для локальной разработки, так как:
* Локальная машина.
* PostgreSQL на той же машине (нет сетевой задержки).
* Учебный проект: небольшое количество параллельных запросов.
