package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// .env не обязателен; если файла нет — ошибка игнорируется
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// fallback — прямой DSN в коде (только для учебного стенда!)
		dsn = "postgres://postgres:1234@localhost:5433/todo?sslmode=disable"
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// 1) Вставим пару задач
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	// 2) Прочитаем список задач
	ctxList, cancelList := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err := repo.ListTasks(ctxList)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

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
}
