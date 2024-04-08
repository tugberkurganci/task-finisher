package models

import (
	"database/sql"
	"time"

	"github.com/go-redis/redis"
)

type Task struct {
	ID         int       `json:"id"`
	Payload    string    `json:"payload,omitempty" validate:"required,min=2"`
	Deadline   time.Time `json:"deadline"`
	Retry      int       `json:"retry"`
	MaxRetries int       `json:"max_retries"`
	Priority   int       `json:"priority"`
	Interval   string    `json:"interval"`
	Completed  bool      `json:"completed"`
	Error      bool      `json:"error"`
}

type Worker struct {
	ID        int
	TaskQueue chan Task
	Redis     *redis.Client
	DB        *sql.DB
}
type User struct {
	ID       int64  `json:"-"`
	Name     string `json:"name,omitempty" validate:"required,min=2"`
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=6"`
}
