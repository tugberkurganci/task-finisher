package service

import (
	"fmt"
	"konzek-mid/models"
	"konzek-mid/repository"

	"sort"

	"time"
)

//go:generate mockgen -destination=../mocks/service/mockTaskService.go -package=service konzek-mid/service TaskService

type TaskService interface {
	EnqueueTask(task models.Task) (models.Task, error)
	GetTaskStatus(taskID int) (models.Task, error)
	StartWorkers()
	ScheduleTasks()
}

type TaskServiceImpl struct {
	TaskQueue chan models.Task
	Workers   []Worker
	Repo      repository.TaskRepository
	Complete  chan bool
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &TaskServiceImpl{
		TaskQueue: make(chan models.Task, 10),
		Complete:  make(chan bool, 10),
		Repo:      repo}
}

func (ts *TaskServiceImpl) EnqueueTask(task models.Task) (models.Task, error) {
	InDB, err := ts.Repo.InsertTask(task)

	return InDB, err
}

func (ts *TaskServiceImpl) ScheduleTasks() {

	go func() {
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-ticker.C:
				tasks, err := ts.Repo.GetPastScheduledTasks()

				if err != nil {
					fmt.Println("hata olustu")
				}
				sortTasksByPriority(tasks)

				for _, task := range tasks {

					duration, err := time.ParseDuration(task.Interval)
					fmt.Println(duration)

					if err != nil {

						fmt.Println("xxxxxxxxx")
						ts.TaskQueue <- task

					}
					if duration.Seconds() > 0 {

						updatedTask := task // Orijinal görevin bir kopyasını oluşturun
						updatedTask.Deadline = task.Deadline.Add(duration)

						ts.TaskQueue <- updatedTask
						continue
					}

				}

			}
		}
	}()
}
func sortTasksByPriority(tasks []models.Task) []models.Task {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority < tasks[j].Priority
	})
	return nil
}
func (ts *TaskServiceImpl) GetTaskStatus(taskID int) (models.Task, error) {
	return ts.Repo.GetTaskByID(taskID)
}

func (ts *TaskServiceImpl) StartWorkers() {

	for i := 0; i < 3; i++ {
		worker := NewWorker(i+1, ts.TaskQueue, ts.Complete, ts.Repo)
		ts.Workers = append(ts.Workers, worker)
		fmt.Println(worker)
		worker.Start()
	}

	go func() {
		for range ts.Complete {
			currentWorkerCount := len(ts.Workers)
			currentQueueLength := len(ts.TaskQueue)

			fmt.Println("zzz", currentWorkerCount)

			desiredWorkerCount := currentQueueLength / 2

			if desiredWorkerCount > currentWorkerCount {

				for i := currentWorkerCount; i < desiredWorkerCount; i++ {
					fmt.Print("arttı")
					worker := NewWorker(i+1, ts.TaskQueue, ts.Complete, ts.Repo)
					ts.Workers = append(ts.Workers, worker)
					worker.Start()
				}
			} else if desiredWorkerCount < currentWorkerCount && desiredWorkerCount > 0 {

				for i := currentWorkerCount - 1; i >= desiredWorkerCount; i-- {

					fmt.Println(i)
					ts.Workers[i].Stop()
					fmt.Println("azaldı")

				}
				ts.Workers = ts.Workers[:desiredWorkerCount]
			}
		}
	}()
}

// Monitor method monitors the task processing.
func (ts *TaskServiceImpl) Monitor() {
	for {
		select {
		case <-time.After(5 * time.Second):
			fmt.Printf("Tasks processed: ")
			time.Sleep(5 * time.Second)
		}
	}
}
