package service

import (
	"context"
	"fmt"
	"konzek-mid/models"
	"konzek-mid/prometheus"
	"konzek-mid/repository"
	"strconv"
	"time"
)

//go:generate mockgen -destination=../mocks/service/mockWorker.go -package=service konzek-mid/service Worker

type Worker interface {
	Start()
	processWithRetry(task models.Task) error
	CompleteTask(task models.Task) error
	Stop()
}
type WorkerImpl struct {
	ID        int
	TaskQueue chan models.Task
	Complete  chan bool
	Repo      repository.TaskRepository
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewWorker(id int, taskQueue chan models.Task, Complete chan bool, repo repository.TaskRepository) Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerImpl{

		ID:        id,
		TaskQueue: taskQueue,
		Complete:  Complete,
		Repo:      repo,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (w *WorkerImpl) Start() {

	go func() {

		for {
			select {
			case task := <-w.TaskQueue:
				fmt.Printf("Worker %d processing task %d: %s\n", w.ID, task.ID, task.Payload)
				if err := w.processWithRetry(task); err != nil {
					fmt.Printf("Failed to mark task %d as completed: %v\n", task.ID, err)
					println(task.Retry)
					task.Retry++

					if task.Retry > task.MaxRetries {
						err := w.Repo.UpdateErrorTask(task)
						if err == nil {
							fmt.Println("updated error task to true")
						}
					} else {
						time.Sleep(time.Second)

						w.TaskQueue <- task
					}

				} else {

					w.Complete <- true
				}

			case <-w.ctx.Done():
				fmt.Println("Worker", w.ID, "stopped")
				return
			}
		}
	}()
}
func (w *WorkerImpl) processWithRetry(task models.Task) error {
	prometheus.RequestCounter.Inc()
	const initialBackoff = 100 * time.Millisecond
	const maxTimeout = 3 * time.Second

	startTime := time.Now()

	for retries := 0; retries < task.MaxRetries; retries++ {
		err := w.CompleteTask(task)
		if err == nil {
			prometheus.TaskSuccess.WithLabelValues(strconv.Itoa(task.ID)).Set(1)
			prometheus.TaskDuration.WithLabelValues(strconv.Itoa(task.ID)).Observe(time.Since(startTime).Seconds())
			return nil
		}

		// Süre kontrolü
		if time.Since(startTime) >= maxTimeout {
			return fmt.Errorf("task %d timed out after %s", task.ID, maxTimeout)
		}

		backoff := time.Duration(retries+1) * initialBackoff
		time.Sleep(backoff)
	}
	processTıme := time.Since(startTime).Seconds()
	prometheus.TaskDuration.WithLabelValues(strconv.Itoa(task.ID)).Observe(processTıme)
	prometheus.TaskSuccess.WithLabelValues(strconv.Itoa(task.ID)).Set(0)

	return fmt.Errorf("task %d failed after maximum retries", task.ID)
}

func (w *WorkerImpl) CompleteTask(task models.Task) error {

	err := w.Repo.MarkTaskCompleted(task)

	if err != nil {
		fmt.Printf("error")
	}

	return err
}

func (w *WorkerImpl) Stop() {

	w.cancel()
}
