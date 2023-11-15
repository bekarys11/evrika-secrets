package auth

type LoginPayload struct {
	Email    string `json:"email" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"token"`
}
