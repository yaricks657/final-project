package todo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yaricks657/final-project/internal/database"
	"github.com/yaricks657/final-project/internal/manager"
)

func GetTask(w http.ResponseWriter, r *http.Request) {
	manager.Mng.Log.LogInfo("Поступил запрос на получение задачи: ", r.RequestURI)

	id := r.URL.Query().Get("id")
	if id == "" {
		manager.Mng.Log.LogError("В запросе должен быть id ", fmt.Errorf("empty id"))
		sendErrorResponse(w, fmt.Sprintf("В запросе должен быть id %s", fmt.Errorf("empty id")))
		return
	}

	task, err := database.GetTask(&manager.Mng, id)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err))
		return
	}

	jsonResponse, err := json.Marshal(task)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при маршалинге данных: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при маршалинге данных: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	manager.Mng.Log.LogInfo("Отправка сообщения успешна с id: ", id)
}
