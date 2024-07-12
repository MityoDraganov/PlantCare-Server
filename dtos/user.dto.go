package dtos

type CreateUserDto struct {
	Username   string `json:"username" validate:"required,min=8"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	RePassword string `json:"rePassword" validate:"required,min=8,eqfield=Password"`
}

type LoginUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	Username string
	Email string
	Token string
}
