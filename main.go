package main

import (
	"fmt"
	"konzek-mid/app"
	"konzek-mid/config"
	"konzek-mid/loggerx"
	"konzek-mid/middleware"
	"konzek-mid/prometheus"
	"konzek-mid/repository"
	"konzek-mid/service"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	prometheus.InitPrometheus()
	loggerx.Init()
	db := config.ConnectDB()

	go func() {
		if err := http.ListenAndServe(":2222", promhttp.Handler()); err != nil {
			fmt.Println("Prometheus sunucusunu başlatırken hata oluştu:", err)
		}
	}()

	defer db.Close()

	// Create task repository
	taskRepo := repository.NewTaskRepository(db)
	// Create task service
	taskService := service.NewTaskService(taskRepo)

	// Create task handler
	taskHandler := app.NewTaskHandler(taskService)

	taskService.StartWorkers()
	taskService.ScheduleTasks()
	jwtMiddleware := middleware.NewJWTMiddleware(service.NewJWTService())
	authService := service.NewAuthService(repository.NewUserRepo(db))

	jwtService := service.NewJWTService()

	userService := service.NewUserService(repository.NewUserRepo(db))

	authHandler := app.NewAuthHandler(authService, jwtService, userService)

	app := fiber.New()
	app.Use(func(ctx *fiber.Ctx) error {
		// Middleware'i atlamak istediğimiz endpointlerin adları
		skipEndpoints := []string{"/register", "/login", "/metrics", "/swagger-ui/index.html"}

		// Endpoint adını kontrol et
		for _, skipEndpoint := range skipEndpoints {
			fmt.Println(ctx.Path())
			if ctx.Path() == skipEndpoint {
				// Middleware'i atla
				fmt.Println("Atladı:", skipEndpoint)
				return ctx.Next()
			}
		}

		// Diğer durumlarda, JWT doğrulamasını yap
		return jwtMiddleware.AuthorizeJWT(ctx)
	})

	// Define routes
	app.Post("/tasks", taskHandler.AddTaskHandler)
	app.Get("/tasks/:id", taskHandler.GetTaskStatusHandler)
	app.Post("/tasks", taskHandler.AddTaskHandler)
	app.Get("/tasks/:id", taskHandler.GetTaskStatusHandler)
	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)

	// Start Fiber server
	log.Fatal(app.Listen(":8080"))
}
