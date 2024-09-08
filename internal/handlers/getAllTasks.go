package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yaricks657/final-project/internal/database"
	"github.com/yaricks657/final-project/internal/manager"
)

// получить задачи из БД
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// проверка запроса на search
	search := r.URL.Query().Get("search")
	if search != "" {
		FindTasksBySearchTerm(w, r)
		return
	}

	tasks, err := database.GetAllTasks(&manager.Mng)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err))
		return
	}

	// упаковка данных
	response := ResponseGetAllTasks{
		Tasks: tasks,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при маршалинге данных: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при маршалинге данных: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при отправке ответа: ", err)
		return
	}
}
