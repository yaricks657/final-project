package todo

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
		GetSearchedTasks(w, r)
		return
	}

	manager.Mng.Log.LogInfo("Поступил запрос на получение всех задач:", r.RequestURI)

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
	w.Write(jsonResponse)
	manager.Mng.Log.LogInfo("Отправка всех сообщений завершена успешно")
}
