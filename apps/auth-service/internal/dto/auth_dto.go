package dto

type RegisterRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	FullName   string `json:"fullName" validate:"required"`
	FamilySize int    `json:"familySize" validate:"gte=0,lte=5"`
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
	ID            string `json:"id"`
	Email         string `json:"email"`
	FullName      string `json:"fullName"`
	FamilySize    int    `json:"familySize"`
	EmailVerified bool   `json:"emailVerified"`
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

// Email verification DTOs
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Password reset DTOs
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
