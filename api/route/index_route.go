package route

import (
	"github.com/chirag1807/task-management-system/api/controller"
	"github.com/chirag1807/task-management-system/api/middleware"
	"github.com/chirag1807/task-management-system/api/repository"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/utils/socket"
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

	router.Route("/api/auth", func(r chi.Router) {
		r.Post("/user-registration", authController.UserRegistration)
		r.Post("/user-login", authController.UserLogin)
		r.With(middleware.VerifyToken(1)).Post("/reset-token", authController.ResetToken)
	})

	router.Route("/api/task", func(r chi.Router) {
		r.Use(middleware.VerifyToken(0))
		r.Post("/create-task", taskController.CreateTask)
		r.Put("/update-task", taskController.UpdateTask)
		r.Get("/get-all-tasks/{Flag}", taskController.GetAllTasks)
		r.Get("/get-tasks-of-team/{TeamID}", taskController.GetTasksofTeam)
	})

	router.Route("/api/team", func(r chi.Router) {
		r.Use(middleware.VerifyToken(0))
		r.Post("/create-team", teamController.CreateTeam)
		r.Put("/add-members-to-team", teamController.AddMembersToTeam)
		r.Put("/remove-members-from-team", teamController.RemoveMembersFromTeam)
		r.Get("/get-all-teams/{Flag}", teamController.GetAllTeams)
		r.Get("/get-team-members/{TeamID}", teamController.GetTeamMembers)
		r.Delete("/left-team/{TeamID}", teamController.LeftTeam)
	})

	router.Route("/api/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken(0))
			r.Get("/get-public-profile-users", userController.GetAllPublicProfileUsers)
			r.Get("/get-my-details", userController.GetMyDetails)
			r.Put("/update-user-profile", userController.UpdateUserProfile)
		})
		r.Post("/send-otp-to-user", userController.SendOTPToUser)
		r.Post("/verify-otp", userController.VerifyOTP)
		r.Put("/reset-user-password", userController.ResetUserPassword)
	})

	return router
}
