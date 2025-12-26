package dto

type RegisterRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	FullName   string `json:"fullName" validate:"required"`
	FamilySize int    `json:"familySize" validate:"required,min=0,max=5"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         UserDTO `json:"user"`
}

type UserDTO struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	FamilySize int    `json:"familySize"`
}

type UpdateProfileRequest struct {
	FullName   string `json:"fullName" validate:"required"`
	FamilySize int    `json:"familySize" validate:"required,min=0,max=5"`
	Email      string `json:"email" validate:"required,email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
}

type ProfileResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	User    UserDTO `json:"user,omitempty"`
}
