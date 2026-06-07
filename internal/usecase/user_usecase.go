package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path"
	"strings"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/dto"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/repository"
)

type userUseCase struct {
	userRepo  repository.UserRepository
	cacheRepo repository.CacheRepository
	fileRepo  repository.FileRepository
}

func NewUserUseCase(userRepo repository.UserRepository, cacheRepo repository.CacheRepository, fileRepo repository.FileRepository) UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
		fileRepo:  fileRepo,
	}
}

func (uc *userUseCase) GetUserProfile(userID string) (*entity.User, error) {
	ctx := context.Background()
	user, err := uc.userRepo.FindByID(userID)
	user.ProfilePict, err = uc.fileRepo.GetPresignedURL(ctx, user.ProfilePict)
	if err != nil {
		return nil, err
	} else {
		return user, err
	}
}

func (uc *userUseCase) UpdateUserProfile(userID string, req *dto.UpdateProfileRequest) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if req.Email != "" && req.Email != user.Email {
		existing, _ := uc.userRepo.FindByEmail(req.Email)
		if existing != nil {
			return errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.Username != "" && req.Username != user.Username {
		existing, _ := uc.userRepo.FIndByUsername(req.Username)
		if existing != nil {
			return errors.New("username already exists")
		}
		user.Username = req.Username
	}

	user.Department = req.Department
	user.IsRemote = req.IsRemote

	return uc.userRepo.UpdateProfile(user)
}

func (uc *userUseCase) UpdateProfilePicture(userID string, imageFile multipart.File, imageFilename string, imageSize int64, imageContentType string) error {

	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	objectName := buildUserAvatarObjectName(
		userID,
		imageFilename,
	)

	objectKey, err := uc.fileRepo.Upload(imageFile, objectName, imageSize, imageContentType)
	if err != nil {
		return err
	}
	user.ProfilePict = objectKey

	return uc.userRepo.UpdateProfilePicture(userID, objectKey)
}

func buildUserAvatarObjectName(
	userID string,
	filename string,
) string {

	normalizedFilename := strings.ReplaceAll(
		filename,
		"\\",
		"/",
	)

	ext := strings.ToLower(
		path.Ext(normalizedFilename),
	)

	if ext == "" {
		ext = ".jpg"
	}

	return fmt.Sprintf(
		"users/%s/avatar%s",
		userID,
		ext,
	)
}

func (uc *userUseCase) GetAllUser() ([]*entity.User, error) {
	users, err := uc.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, ticket := range users {
		if ticket.ProfilePict != "" {
			imageURL, err := uc.fileRepo.GetPresignedURL(context.Background(), ticket.ProfilePict)
			if err != nil {
				return nil, err
			}

			ticket.ProfilePict = imageURL
		}
	}
	return users, nil

}

func (uc *userUseCase) DeleteUser(userID string) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user.Role == entity.RoleAdmin {
		return errors.New("admin cannot be deleted")
	}
	ctx := context.Background()
	_ = uc.cacheRepo.Delete(ctx, "dashboard:overview")
	return uc.userRepo.Delete(userID)
}
