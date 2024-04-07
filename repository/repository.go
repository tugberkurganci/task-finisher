package repository

import (
	"database/sql"
	"fmt"
	"konzek-mid/models"
	"time"
)

type TaskRepositoryImpl struct {
	DB *sql.DB
}

//go:generate mockgen -destination=../mocks/repository/mockTaskRepository.go -package=repository konzek-mid/repository TaskRepository

type TaskRepository interface {
	InsertTask(task models.Task) (models.Task, error)
	MarkTaskCompleted(task models.Task) error
	GetTaskByID(taskID int) (models.Task, error)
	GetPastScheduledTasks() ([]models.Task, error)
	UpdateErrorTask(task models.Task) error
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &TaskRepositoryImpl{DB: db}
}
func (tr *TaskRepositoryImpl) InsertTask(task models.Task) (models.Task, error) {
	// SQL sorgusu çalıştırılırken, RETURNING ifadesi kullanılarak eklenen görevin tüm alanları alınır
	var insertedTask models.Task

	// String olarak tutulan interval'i time.Duration'a dönüştürme

	err := tr.DB.QueryRow("INSERT INTO tasks (payload, deadline, retry, max_retries, priority, task_interval, completed, error) VALUES ($1, $2, $3, $4, $5, $6, $7,$8) RETURNING id, payload, deadline, retry, max_retries, priority, task_interval, completed,error",
		task.Payload, task.Deadline, 0, 5, task.Priority, task.Interval, false, false).Scan(&insertedTask.ID, &insertedTask.Payload, &insertedTask.Deadline, &insertedTask.Retry, &insertedTask.MaxRetries, &insertedTask.Priority, &insertedTask.Interval, &insertedTask.Completed, &insertedTask.Error)
	if err != nil {
		fmt.Println(err)
		return models.Task{}, err
	}

	return insertedTask, nil
}

// MarkTaskCompleted marks a task as completed in the database.
func (tr *TaskRepositoryImpl) MarkTaskCompleted(task models.Task) error {
	intervalDuration, _ := time.ParseDuration(task.Interval)
	if intervalDuration > 0 {
		_, err := tr.DB.Exec("UPDATE tasks SET deadline = $1 WHERE id = $2", task.Deadline, task.ID)

		if err != nil {
			return err
		}
	} else {
		_, err := tr.DB.Exec("UPDATE tasks SET deadline = $1, completed = $2 WHERE id = $3", task.Deadline, true, task.ID)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetTaskByID retrieves a task from the database by its ID.
func (tr *TaskRepositoryImpl) GetTaskByID(taskID int) (models.Task, error) {
	var task models.Task
	err := tr.DB.QueryRow("SELECT id, payload, deadline, retry, max_retries, priority, completed, error FROM tasks WHERE id = $1", taskID).Scan(
		&task.ID, &task.Payload, &task.Deadline, &task.Retry, &task.MaxRetries, &task.Priority, &task.Completed, &task.Error)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

// GetPastScheduledTasks method retrieves tasks with deadline within the last 1 hour,
func (tr *TaskRepositoryImpl) GetPastScheduledTasks() ([]models.Task, error) {
	// Son bir saat içinde bitmemiş görevleri al
	tasks, err := tr.DB.Query("SELECT id, payload, deadline, retry, max_retries, priority, task_interval, completed FROM tasks  WHERE deadline < NOW() AND completed = FALSE AND error = FALSE")
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer tasks.Close()
	println(tasks)
	tasksSlice := make([]models.Task, 0)
	for tasks.Next() {
		var task models.Task
		if err := tasks.Scan(&task.ID, &task.Payload, &task.Deadline, &task.Retry, &task.MaxRetries, &task.Priority, &task.Interval, &task.Completed); err != nil {
			fmt.Println("GetPastScheduledTasks - Veri çekme hatası:", err)
			return nil, err

		}
		tasksSlice = append(tasksSlice, task)
	}
	return tasksSlice, nil
}
func (tr *TaskRepositoryImpl) UpdateErrorTask(task models.Task) error {

	_, err := tr.DB.Exec("UPDATE tasks SET error = $1 WHERE id = $2", true, task.ID)

	if err != nil {
		return err
	}
	return nil
}
