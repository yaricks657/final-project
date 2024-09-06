package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yaricks657/final-project/internal/manager"
)

// обработчик для api/sign
func SignIn(w http.ResponseWriter, r *http.Request) {
	manager.Mng.Log.LogInfo("поступил запрос на авторизацию ", r.RemoteAddr)

	// проверка метода
	if r.Method != http.MethodPost {
		manager.Mng.Log.LogWarn("Некорректный метод запроса")
		sendErrorResponse(w, "Некорректный метод запроса", http.StatusBadRequest)
		return
	}

	// декодирование тела запроса с паролем
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		manager.Mng.Log.LogError("Неверный формат запроса", err)
		sendErrorResponse(w, fmt.Sprintf("Некорректный метод запроса %s", err), http.StatusBadRequest)
		return
	}

	// проверка наличия пароля на сервере
	storedPassword := manager.Mng.Cnf.Password
	fmt.Println(storedPassword, " ", req.Password)
	if storedPassword == "" {
		manager.Mng.Log.LogError("Пароль не настроен на сервере", fmt.Errorf("No password"))
		sendErrorResponse(w, fmt.Sprintln("Пароль не настроен на сервере"), http.StatusInternalServerError)
		return
	}

	// проверка хешей паролей на совпадение
	if hashPassword(req.Password) != hashPassword(storedPassword) {
		manager.Mng.Log.LogWarn("Пользователь ввел неверный пароль")
		sendErrorResponse(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	// выдача подписанного JWT-токена
	JWTSecret := manager.Mng.Cnf.JWTSecret
	if JWTSecret == "" {
		manager.Mng.Log.LogError("Не настроен JWT на сервере", fmt.Errorf("No JWT-token"))
		sendErrorResponse(w, fmt.Sprintln("Не настроен JWT на сервере"), http.StatusInternalServerError)
		return
	}
	token, err := generateJWT(hashPassword(req.Password), []byte(JWTSecret))
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при создании токена", err)
		sendErrorResponse(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, token)

}

// Хэширование пароля
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// Генерация JWT токена
func generateJWT(passwordHash string, secret []byte) (string, error) {
	claims := &jwt.MapClaims{
		"passwordHash": passwordHash,
		"exp":          time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// отправка на клиент ответа об ошибке
func sendErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := loginResponse{
		Error: errorMsg,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при декодировании (sendErrorResponse)", err)
	}

}

// отправка на клиент ответа об успехе
func sendSuccessResponse(w http.ResponseWriter, token string) {
	response := loginResponse{
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при декодировании (sendSuccessResponse)", err)
	}

	// Устанавливаем куку с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(8 * time.Hour),
	})
}
