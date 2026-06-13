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

// GetMyProfile godoc
//
//	@Summary		Get current user profile
//	@Description	Get authenticated user profile
//	@Tags			Users
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Router			/users/me [get]
func (h *UserHandler) GetMyProfile(c *fiber.Ctx) error {
	// userId := c.Locals("userID").(string)
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

// UpdateMyProfile godoc
//
//	@Summary		Update profile
//	@Description	Update authenticated user profile
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UpdateProfileRequest	true	"Profile Data"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/users/me [patch]
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

// UpdateProfilePicture godoc
//
//	@Summary		Upload profile picture
//	@Description	Upload user profile picture
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			profile_picture	formData	file	true	"Profile Picture"
//	@Success		200				{object}	map[string]string
//	@Failure		400				{object}	map[string]interface{}
//	@Router			/users/me/profile-picture [patch]
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

// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Requires ADMIN role
//	@Tags			Admin
//	@Security		BearerAuth
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	map[string]string
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/admin/users/{id} [delete]
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

// GetAllUsers godoc
//
//	@Summary		Get all users
//	@Description	Requires ADMIN role
//	@Tags			Admin
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		403	{object}	map[string]interface{}
//	@Router			/admin/users [get]
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
