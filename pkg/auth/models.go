package auth

type LoginPayload struct {
	Login    string `json:"login" validate:"required" example:"kamilla.n"`
	Password string `json:"password" validate:"required" example:"password123"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDA3OTgzMzcsIm5hbWUiOiJiZWthcnlzIiwidXNlcl9pZCI6MX0.YoLwVoMnGvA7q44dFQJ4E4qBghr3zYDKbJNrhV1yrC0"`
}
