package dto

type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	ExpiresIn int    `json:"expiresIn"`
}
