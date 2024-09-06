package todo

import (
	"fmt"
	"net/http"

	"github.com/yaricks657/final-project/internal/database"
	"github.com/yaricks657/final-project/internal/manager"
)

// Удаление задачи
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	manager.Mng.Log.LogInfo("поступил запрос на удаление задачи: ", r.RequestURI)

	// проверка на наличие search в запросе
	id := r.URL.Query().Get("id")
	if id == "" {
		manager.Mng.Log.LogError("Отсутствует id в запросе ", fmt.Errorf(""))
		sendErrorResponse(w, fmt.Sprintln("Отсутствует id в запросе"))
		return
	}

	err := database.DeleteTask(&manager.Mng, id)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка удалении задачи из БД: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка удалении задачи из БД: %s", err))
		return
	}

	manager.Mng.Log.LogInfo("Задача успешно удалена из БД", id)
	// Отправляем пустой JSON в случае успеха
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
