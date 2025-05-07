package resources

import (
	"ticatag_backend/models"
)

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"name"`
	Email     string `json:"email"`
	Role      string `bson:"role" json:"role"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
}

func NewUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func NewUserListResponse(users []models.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, NewUserResponse(user))
	}
	return responses
}
