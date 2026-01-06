package handler

import (
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	ticketUsecase usecase.TicketUseCase
}

func NewTicketHandler(ticketUsecase usecase.TicketUseCase) *TicketHandler {
	return &TicketHandler{ticketUsecase: ticketUsecase}
}

type CreateTicketRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateStatusRequest struct {
	Status entity.TicketStatus `json:"status" validate:"required,oneof=OPEN IN_PROGRESS DONE"`
}

func (h *TicketHandler) CreateTicket(c *fiber.Ctx) error {
	var req CreateTicketRequest
	userIdRaw := c.Locals("userID")
	if userIdRaw == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  "need useer id bro",
			"userID": userIdRaw,
		})
	}

	userId, ok := userIdRaw.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ticket, err := h.ticketUsecase.CreateTicket(userId, req.Title, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(ticket)

}

func (h *TicketHandler) GetAllTickets(c *fiber.Ctx) error {
	tickets, err := h.ticketUsecase.GetAllTicket()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Jumlah":  len(tickets),
		"tickets": tickets,
	})
}

func (h *TicketHandler) GetUserTickets(c *fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	tickets, err := h.ticketUsecase.GetUserTickets(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(tickets)
}

func (h *TicketHandler) UpdateTicketStatus(c *fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	userRole := c.Locals("userRole").(string)
	ticketID := c.Params("id")
	var req UpdateStatusRequest

	if userRole != string(entity.RoleAdmin) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Admin only",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ticket, err := h.ticketUsecase.UpdateTicketStatus(ticketID, string(req.Status), userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "update success",
		"ticket":  ticket,
	})

}
