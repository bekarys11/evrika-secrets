package auth

type LoginPayload struct {
	Login    string `json:"login" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
