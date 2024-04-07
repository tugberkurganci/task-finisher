package app

import (
	"konzek-mid/models"
	"konzek-mid/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	TaskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{TaskService: taskService}
}

// AddTaskHandler handles adding tasks via HTTP POST request.
func (th *TaskHandler) AddTaskHandler(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	if _, err := th.TaskService.EnqueueTask(task); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(http.StatusCreated).SendString("Task added successfully")
}

// GetTaskStatusHandler handles getting task status via HTTP GET request.
func (th *TaskHandler) GetTaskStatusHandler(c *fiber.Ctx) error {
	taskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid task ID")
	}
	task, err := th.TaskService.GetTaskStatus(taskID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(task)
}
