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
