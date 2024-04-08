package app

import (
	"fmt"
	"konzek-mid/globalerror"
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
// @Summary Add a new task
// @Description Add a new task to the task queue
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task object to be added"
// @Success 201 {string} string "Task added successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks [post]
func (th *TaskHandler) AddTaskHandler(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	if errors := globalerror.Validate(task); len(errors) > 0 && errors[0].HasError {
		return globalerror.HandleValidationErrors(c, errors)
	}
	InDB, err := th.TaskService.EnqueueTask(task)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(globalerror.ErrorResponse{
			Status: http.StatusInternalServerError,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Task",
					Description: "An error occurred while adding the task",
				},
			},
		})
	}

	return c.Status(http.StatusCreated).SendString(fmt.Sprintf("Task added successfully. Task ID: %d", InDB.ID))
}

// GetTaskStatusHandler handles getting task status via HTTP GET request.
// @Summary Get task status
// @Description Get the status of a task by its ID
// @Tags Tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task "Task status"
// @Failure 400 {string} string "Invalid task ID"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks/{id} [get]
func (th *TaskHandler) GetTaskStatusHandler(c *fiber.Ctx) error {
	taskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid task ID")
	}
	task, err := th.TaskService.GetTaskStatus(taskID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(globalerror.ErrorResponse{
			Status: http.StatusInternalServerError,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Task",
					Description: "An error occurred while fetching the task",
				},
			},
		})
	}
	return c.JSON(task)
}
