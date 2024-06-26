package route

import (
	"github.com/chirag1807/task-management-system/api/controller"
	"github.com/chirag1807/task-management-system/api/middleware"
	"github.com/chirag1807/task-management-system/api/repository"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/utils/socket"
	chi_middleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitializeRouter(dbConn *pgx.Conn, redisClient *redis.Client, rabbitmqConn *amqp.Connection, socketServer *socketio.Server) *chi.Mux {
	router := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})
	socket.SocketEvents(socketServer)
	router.Use(chi_middleware.Logger)
	router.Handle("/socket.io/", c.Handler(socketServer))

	authRepository := repository.NewAuthRepo(dbConn)
	authService := service.NewAuthService(authRepository)
	authController := controller.NewAuthController(authService)

	taskRepository := repository.NewTaskRepo(dbConn, redisClient, socketServer)
	taskService := service.NewTaskService(taskRepository)
	taskController := controller.NewTaskController(taskService)

	teamRepository := repository.NewTeamRepo(dbConn, redisClient)
	teamService := service.NewTeamService(teamRepository)
	teamController := controller.NewTeamController(teamService)

	userRepository := repository.NewUserRepo(dbConn, rabbitmqConn)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/registration", authController.UserRegistration)
			r.Post("/login", authController.UserLogin)
			r.With(middleware.VerifyToken(1)).Post("/refresh-token", authController.RefreshToken)
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Use(middleware.VerifyToken(0))
			r.Post("/", taskController.CreateTask)
			r.Put("/{TaskID}", taskController.UpdateTask)
			r.Get("/", taskController.GetAllTasks)
			r.Get("/team/{TeamID}", taskController.GetTasksofTeam)
		})

		r.Route("/teams", func(r chi.Router) {
			r.Use(middleware.VerifyToken(0))
			r.Post("/", teamController.CreateTeam)
			r.Post("/{TeamID}/members", teamController.AddMembersToTeam)
			r.Delete("/{TeamID}/members", teamController.RemoveMembersFromTeam)
			r.Get("/", teamController.GetAllTeams)
			r.Get("/{TeamID}/members", teamController.GetTeamMembers)
			r.Delete("/leave/{TeamID}", teamController.LeaveTeam)
		})

		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.VerifyToken(0))
				r.Get("/public-privacy", userController.GetAllPublicPrivacyUsers)
				r.Get("/profile", userController.GetMyDetails)
				r.Put("/profile", userController.UpdateUserProfile)
			})
			r.Post("/send-otp", userController.SendOTPToUser)
			r.Post("/verify-otp", userController.VerifyOTP)
			r.Put("/reset-password", userController.ResetUserPassword)
		})
	})

	router.Route("/socket_events", func(r chi.Router) {
		r.Get("/", socket.RenderSocketEventsDoc)
	})

	return router
}
