package route

import (
	"github.com/chirag1807/task-management-system/api/controller"
	"github.com/chirag1807/task-management-system/api/middleware"
	"github.com/chirag1807/task-management-system/api/repository"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

func InitializeRouter(dbConn *pgx.Conn, redisClient *redis.Client) *chi.Mux {
	router := chi.NewRouter()

	authRepository := repository.NewAuthRepo(dbConn, redisClient)
	authService := service.NewAuthService(authRepository)
	authController := controller.NewAuthController(authService)

	taskRepository := repository.NewTaskRepo(dbConn, redisClient)
	taskService := service.NewTaskService(taskRepository)
	taskController := controller.NewTaskController(taskService)

	teamRepository := repository.NewTeamRepo(dbConn, redisClient)
	teamService := service.NewTeamService(teamRepository)
	teamController := controller.NewTeamController(teamService)

	userRepository := repository.NewUserRepo(dbConn, redisClient)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	router.Route("/api/auth", func(r chi.Router) {
		r.Post("/user-registration", authController.UserRegistration)
		r.Post("/user-login", authController.UserLogin)
		r.With(middleware.VerifyToken(1)).Post("/reset-token", authController.ResetToken)
	})

	router.Route("/api/task", func(r chi.Router) {
		router.Post("/create-task", taskController.CreateTask)
		router.Get("/get-all-tasks/{Flag}", taskController.GetAllTasks)
		router.Get("/get-tasks-of-team/{TeamID}", taskController.GetTasksofTeam)
		router.Put("/update-task", taskController.UpdateTask)
		router.Delete("/delete-task/{TaskID}", taskController.DeleteTask)
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
