package handler

import (
	"mime/multipart"
	"strconv"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/dto"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUsecase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUsecase: userUseCase}
}

func (h *UserHandler) GetMyProfile(c *fiber.Ctx) error {

	userIdRaw := c.Locals("userID")
	if userIdRaw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

	user, err := h.userUsecase.GetUserProfile(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) UpdateMyProfile(c *fiber.Ctx) error {
	var req dto.UpdateProfileRequest

	userIdRaw := c.Locals("userID")
	if userIdRaw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

	remote, err := strconv.ParseBool(
		c.FormValue("remote"),
	)
	if err != nil {
		return err
	}

	req.IsRemote = remote

	if err := h.userUsecase.UpdateUserProfile(userId, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Updated",
	})
}

func (h *UserHandler) UpdateMyProfilePicture(c *fiber.Ctx) error {
	var imageFile multipart.File
	var imageHeader *multipart.FileHeader

	userIDRaw := c.Locals("userID")
	if userIDRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid user id format",
		})
	}

	if !isMultipartRequest(c) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "content type must be multipart/form-data",
		})
	}

	var err error

	imageHeader, err = getOptionalFormFile(
		c,
		"profile_picture",
		"avatar",
		"image",
		"file",
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid multipart form",
		})
	}

	if imageHeader == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "profile picture is required",
		})
	}

	imageFile, err = imageHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to read uploaded image",
		})
	}
	defer imageFile.Close()

	imageFilename, imageSize, imageContentType := uploadedFileMetadata(imageHeader)

	err = h.userUsecase.UpdateProfilePicture(userID, imageFile, imageFilename, imageSize, imageContentType)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "profile picture updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	if err := h.userUsecase.DeleteUser(userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}

func (h *UserHandler) GetAllUserUser(c *fiber.Ctx) error {
	userRole := c.Locals("userRole").(string)
	if userRole != string(entity.RoleAdmin) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Admin only",
		})
	}

	users, err := h.userUsecase.GetAllUser()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Jumlah":  len(users),
		"tickets": users,
	})
}
