package dto

type UpdateProfileRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Department string `json:"department"`
	IsRemote   bool   `json:"remote"`
}
