package handler

import (
	"mime/multipart"
	"strings"

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
	Title       string `json:"title" form:"title" validate:"required"`
	Image       string `json:"image" form:"image"`
	Description string `json:"description" form:"description" validate:"required"`
}

type UpdateStatusRequest struct {
	Status entity.TicketStatus `json:"status" validate:"required,oneof=OPEN IN_PROGRESS DONE"`
}

func (h *TicketHandler) CreateTicket(c *fiber.Ctx) error {
	var req CreateTicketRequest
	var imageFile multipart.File
	var imageHeader *multipart.FileHeader

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

	if isMultipartRequest(c) {
		req.Title = c.FormValue("title")
		req.Description = c.FormValue("description")
		req.Image = c.FormValue("image")

		var err error
		imageHeader, err = getOptionalFormFile(c, "image", "file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid multipart form",
			})
		}

		if imageHeader != nil {
			imageFile, err = imageHeader.Open()
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Failed to read uploaded image",
				})
			}
			defer imageFile.Close()
		}
	} else {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
	}

	imageFilename, imageSize, imageContentType := uploadedFileMetadata(imageHeader)

	ticket, err := h.ticketUsecase.CreateTicket(userId, req.Title, req.Description, req.Image, imageFile, imageFilename, imageSize, imageContentType)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(ticket)

}

func isMultipartRequest(c *fiber.Ctx) bool {
	return strings.Contains(strings.ToLower(c.Get("Content-Type")), "multipart/form-data")
}

func getOptionalFormFile(c *fiber.Ctx, names ...string) (*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	for _, name := range names {
		files := form.File[name]
		if len(files) > 0 {
			return files[0], nil
		}
	}

	return nil, nil
}

func uploadedFileMetadata(fileHeader *multipart.FileHeader) (string, int64, string) {
	if fileHeader == nil {
		return "", 0, ""
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return fileHeader.Filename, fileHeader.Size, contentType
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

	validStatus := map[string]bool{
		string(entity.StatusDone):       true,
		string(entity.StatusInProgress): true,
		string(entity.StatusOpen):       true,
	}

	if !validStatus[string(req.Status)] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status must be  'DONE' or 'IN_PROGRESS'",
		})
	}

	ticket, err := h.ticketUsecase.UpdateTicketStatus(ticketID, userId, entity.TicketStatus(req.Status))
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

func (h *TicketHandler) GetDasboardOverview(c *fiber.Ctx) error {
	userRole := c.Locals("userRole").(string)
	if userRole != string(entity.RoleAdmin) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Admin only",
		})
	}

	stat, err := h.ticketUsecase.GetSummaryStat()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"stat": stat,
	})
}
