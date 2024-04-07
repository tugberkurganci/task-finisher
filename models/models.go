package models

import (
	"database/sql"
	"time"

	"github.com/go-redis/redis"
)

type Task struct {
	ID         int
	Payload    string
	Deadline   time.Time
	Retry      int
	MaxRetries int
	Priority   int
	Interval   string
	Completed  bool
	Error      bool
}

type Worker struct {
	ID        int
	TaskQueue chan Task
	Redis     *redis.Client
	DB        *sql.DB
}
type User struct {
	ID       int64  `gorm:"primary_key:auto_increment" json:"-"`
	Name     string `gorm:"type:varchar(100)" json:"name,omitempty" validate:"required,min=2"`
	Email    string `gorm:"type:varchar(100);unique;" json:"email,omitempty" validate:"required,email"`
	Password string `gorm:"type:varchar(100)" json:"password,omitempty" validate:"required,min=6"`
}
