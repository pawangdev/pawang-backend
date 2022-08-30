package response

import (
	"pawang-backend/entity"
)

type UserProfileResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Gender string `json:"gender"`
}

func FormatUserProfileResponse(user entity.User) UserProfileResponse {
	response := UserProfileResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Phone:  user.Phone,
		Gender: user.Gender,
	}

	return response
}

type RegisterUserResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Gender string `json:"gender"`
	Token  string `json:"token"`
}

func FormatRegisterUserResponse(user entity.User, token string) RegisterUserResponse {
	response := RegisterUserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Phone:  user.Phone,
		Gender: user.Gender,
		Token:  token,
	}

	return response
}
