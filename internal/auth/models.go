package auth

// распаковка запроса от клиента
type loginRequest struct {
	Password string `json;password`
}

// структура для ответа клиенту
type loginResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}
