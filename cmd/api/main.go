package main

import (
	"log"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/config"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/delivery/http/handler"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/delivery/http/middleware"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/infrastructure/database"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/infrastructure/repository"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/infrastructure/bucket"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/usecase"
	pkg "github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.LoadConfig()

	internalMinioClient := bucket.InitMinio(cfg.Minio.Endpoint, cfg.Minio.AccessKey, cfg.Minio.SecretKey)
	publicMinioClient := bucket.InitMinio(cfg.Minio.PublicEndpoint, cfg.Minio.AccessKey, cfg.Minio.SecretKey)

	db, err := database.NewMySQLConnection(*cfg)
	if err != nil {
		log.Fatal("Failed to connect to database : ", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run database migrations: ", err)
	}

	cache, err := database.NewRedisClient(*cfg)
	if cache != nil {
		log.Println("Berhasil Terhubung ke Redis")
	} else {
		log.Panicln("Damn gak berhasil connect Redis")
	}

	jwtService := pkg.NewJWTservice(cfg.JWT.Secret, cfg.JWT.Expiry)

	userRepo := repository.NewMySQLUserRepository(db)
	ticketRepo := repository.NewMySQLTicketRepository(db)
	fileRepo := repository.NewMinioRepository(internalMinioClient, publicMinioClient, cfg.Minio.BucketName)
	cacheRepo := repository.NewRedisRepository(cache)

	authUsecase := usecase.NewAuthUseCase(userRepo, jwtService)
	ticketUsecase := usecase.NewTicketUseCase(ticketRepo, userRepo, fileRepo, cacheRepo)
	healthUsecase := usecase.NewHealthUseCase(ticketRepo, userRepo, fileRepo, cacheRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	ticketHandler := handler.NewTicketHandler(ticketUsecase)
	healthHandler := handler.NewHealthHandler(healthUsecase)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	api := app.Group("/api/v1")

	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)
	api.Get("/health/live", healthHandler.Liveness)
	api.Get("/health/ready", healthHandler.Readiness)

	ticketGroup := api.Group("/tickets")
	ticketGroup.Use(middleware.AuthMiddleware(jwtService))

	ticketGroup.Post("/", ticketHandler.CreateTicket)
	ticketGroup.Get("/my-tickets", ticketHandler.GetUserTickets)

	adminGroup := ticketGroup.Group("/admin")
	adminGroup.Use(middleware.RequireRole("ADMIN"))
	adminGroup.Get("/all", ticketHandler.GetAllTickets)
	adminGroup.Get("/dashboard", ticketHandler.GetDasboardOverview)
	adminGroup.Put("/:id/status", ticketHandler.UpdateTicketStatus)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
