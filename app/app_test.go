package app_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"konzek-mid/app"
	"konzek-mid/mocks/service"
	"konzek-mid/models"
	"konzek-mid/repository"
	x "konzek-mid/service"
	"log"
	"time"

	_ "github.com/lib/pq"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var mockService *service.MockTaskService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockTaskService(ctrl)

	return func() { ctrl.Finish() }
}

func TestTaskHandler_AddTaskHandler(t *testing.T) {
	trd := setup(t)
	defer trd()

	taskHandler := app.NewTaskHandler(mockService)

	router := fiber.New()
	router.Post("/api/tasks", taskHandler.AddTaskHandler)

	// Prepare request body
	task := models.Task{ID: 1, Payload: "Test Content"}
	taskJSON, _ := json.Marshal(task)

	// Mock TaskService to return nil error
	mockService.EXPECT().EnqueueTask(gomock.Any()).Return(models.Task{ID: 1, Payload: "Test Content"}, nil)

	// Perform request
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := router.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestTaskHandler_GetTaskStatusHandler(t *testing.T) {
	trd := setup(t)
	defer trd()

	taskHandler := app.NewTaskHandler(mockService)

	router := fiber.New()
	router.Get("/api/tasks/:id", taskHandler.GetTaskStatusHandler)

	// Mock TaskService to return a task and nil error
	mockTask := models.Task{ID: 1, Payload: "Test Content"}
	mockService.EXPECT().GetTaskStatus(1).Return(mockTask, nil)

	// Perform request
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/1", nil)
	resp, err := router.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCheckTaskCompletion(t *testing.T) {

	db, err := sql.Open("postgres", "dbname=xx user=postgres password=test host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata oluştu: %v", err)
	}
	defer db.Close()

	clearDatabase(db)
	taskRepo := repository.NewTaskRepository(db)
	taskService := x.NewTaskService(taskRepo)
	taskHandler := app.NewTaskHandler(taskService)

	router := fiber.New()
	router.Get("/api/tasks/:id", taskHandler.GetTaskStatusHandler)
	router.Post("/api/tasks", taskHandler.AddTaskHandler)

	// Task oluştur
	task := models.Task{
		Payload:  "Test Content",
		Deadline: time.Now(),
	}
	taskJSON, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	resp1, _ := router.Test(req)
	// Bekleme süresi (30 saniye)
	time.Sleep(20 * time.Second)

	req2 := httptest.NewRequest(http.MethodGet, "/api/tasks/346", nil)
	resp2, _ := router.Test(req2)

	fmt.Println(resp2)
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

}
func clearDatabase(db *sql.DB) {
	_, err := db.Exec("DELETE FROM tasks")
	if err != nil {
		log.Fatalf("Veritabanını temizlerken hata oluştu: %v", err)
	}
}
