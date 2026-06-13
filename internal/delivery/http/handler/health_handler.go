package handler

import (
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
	healthUseCase usecase.HealthUseCase
}

func NewHealthHandler(healthUseCase usecase.HealthUseCase) *HealthHandler {
	return &HealthHandler{
		healthUseCase: healthUseCase,
	}
}

// @Summary	Liveness Probe
// @Tags		Health
// @Produce	json
// @Success	200	{object}	map[string]string
// @Router		/health/live [get]
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}

// @Summary	Readiness Probe
// @Tags		Health
// @Produce	json
// @Success	200	{object}	map[string]interface{}
// @Router		/health/ready [get]
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	statusService, err := h.healthUseCase.Readiness()

	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  err,
			"service": statusService,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ready",
		"service": statusService,
	})
}

// func (h *HealthHandler) Metrics(c *fiber.Ctx) error {
// 	return adaptor.HTTPHandler(promhttp.Handler())(c)

// }
