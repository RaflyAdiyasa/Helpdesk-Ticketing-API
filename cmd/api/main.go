package main

import (
	"log"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/config"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/delivery/http/handler"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/delivery/http/middleware"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/infrastructure/database"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/infrastructure/repository"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/usecase"
	pkg "github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewMySQLConnection(*cfg)
	if err != nil {
		log.Fatal("Failed to connect to database : ", err)
	}

	database.RunMigrations(db)

	jwtService := pkg.NewJWTservice(cfg.JWT.Secret, cfg.JWT.Expiry)

	userRepo := repository.NewMySQLUserRepository(db)
	ticketRepo := repository.NewMySQLTicketRepository(db)

	authUsecase := usecase.NewAuthUseCase(userRepo, jwtService)
	ticketUsecase := usecase.NewTicketUseCase(ticketRepo, userRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	ticketHandler := handler.NewTicketHandler(ticketUsecase)

	app := fiber.New()

	api := app.Group("/api/v1")

	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	ticketGroup := api.Group("/tickets")
	ticketGroup.Use(middleware.AuthMiddleware(jwtService))

	ticketGroup.Post("/", ticketHandler.CreateTicket)
	ticketGroup.Get("/my-tickets", ticketHandler.GetUserTickets)

	adminGroup := ticketGroup.Group("/admin")
	adminGroup.Use(middleware.RequireRole("ADMIN"))
	adminGroup.Get("/all", ticketHandler.GetAllTickets)
	adminGroup.Put("/:id/status", ticketHandler.UpdateTicketStatus)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
