package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Task represents a task in the database
type Task struct {
	ID        int
	Name      string
	Completed bool
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "TaskManager"
)

func main() {
	// Connect to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	fmt.Println("Connected to the database")

	// Create a table for tasks if it does not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        completed BOOLEAN NOT NULL DEFAULT FALSE
    )`)
	if err != nil {
		log.Fatal("Failed to create tasks table:", err)
	}

	err = CreateTask(db, "Task 1")
	if err != nil {
		log.Fatal("Failed to create task:", err)
	}
	fmt.Println("Task created successfully")

	// Read all tasks
	tasks, err := ReadTasks(db)
	if err != nil {
		log.Fatal("Failed to read tasks:", err)
	}
	fmt.Println("Tasks:")
	for _, task := range tasks {
		fmt.Printf("ID: %d, Name: %s, Completed: %t\n", task.ID, task.Name, task.Completed)
	}

	// Update a task
	err = UpdateTask(db, 1, true)
	if err != nil {
		log.Fatal("Failed to update task:", err)
	}
	fmt.Println("Task updated successfully")

	// Delete a task
	err = DeleteTask(db, 1)
	if err != nil {
		log.Fatal("Failed to delete task:", err)
	}
	fmt.Println("Task deleted successfully")
}
func CreateTask(db *sql.DB, name string) error {
	_, err := db.Exec("INSERT INTO tasks (name) VALUES ($1)", name)
	return err
}

func ReadTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, name, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func UpdateTask(db *sql.DB, id int, completed bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE tasks SET completed = $1 WHERE id = $2", completed, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteTask(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
