package request

type RegisterUserRequest struct {
	Name     string `json:"name" form:"name" xml:"name" validate:"required"`
	Email    string `json:"email" form:"email" xml:"email" validate:"required,email"`
	Password string `json:"password" form:"password" xml:"password" validate:"required,min=8"`
	Phone    string `json:"phone" form:"phone" xml:"phone" validate:"required"`
	Gender   string `json:"gender" form:"gender" xml:"gender" validate:"required,oneof=male female"`
}

type LoginUserInput struct {
	Email    string `json:"email" form:"email" xml:"email" validate:"required,email"`
	Password string `json:"password" form:"password" xml:"password" validate:"required,min=8"`
}

type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password" xml:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" form:"new_password" xml:"new_password" validate:"required,min=8"`
}

type UserChangeProfileRequest struct {
	Name   string `json:"name" form:"name" xml:"name" validate:"required"`
	Phone  string `json:"phone" form:"phone" xml:"phone" validate:"required"`
	Gender string `json:"gender" form:"gender" xml:"gender" validate:"required,oneof=male female"`
}

type UserResetPasswordRequest struct {
	Email string `json:"email" form:"email" xml:"email" validate:"required,email"`
}

type UserResetPasswordTokenRequest struct {
	Token    string `json:"token" form:"token" xml:"token" validate:"required,min=8"`
	Email    string `json:"email" form:"email" xml:"email" validate:"required,email"`
	Password string `json:"password" form:"password" xml:"password" validate:"required,min=8"`
}
