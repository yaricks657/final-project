package todo

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yaricks657/final-project/internal/database"
	"github.com/yaricks657/final-project/internal/manager"
)

// Процесс выполнения задачи
func DoneTask(w http.ResponseWriter, r *http.Request) {
	manager.Mng.Log.LogInfo("Поступил запрос на выполнение задачи ", r.RequestURI)

	// проверка на наличие search в запросе
	id := r.URL.Query().Get("id")
	if id == "" {
		manager.Mng.Log.LogError("Отсутствует id в запросе ", fmt.Errorf(""))
		sendErrorResponse(w, fmt.Sprintln("Отсутствует id в запросе"))
		return
	}

	// получаем задачу из БД по id
	task, err := database.GetTask(&manager.Mng, id)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err))
		return
	}

	// удаляем задачу, если отсутствует правило
	if task.Repeat == "" {
		err = database.DeleteTask(&manager.Mng, id)
		if err != nil {
			manager.Mng.Log.LogError("Ошибка удалении задачи из БД: ", err)
			sendErrorResponse(w, fmt.Sprintf("Ошибка удалении задачи из БД: %s", err))
			return
		}
		manager.Mng.Log.LogInfo("Задача успешно удалена из БД и выполнена")
		// Отправляем пустой JSON в случае успеха
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
		return
	}

	// перерасчитываем дату согласно правилу
	today := time.Now().Truncate(24 * time.Hour)
	newDate, err := NextDate(today, task.Date, task.Repeat)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при применении правила повторения", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при применении правила повторения %s", err))
		return
	}

	// запись новой даты в БД
	task.Date = newDate
	err = database.ChangeTask(&task, &manager.Mng)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обновлении задачи в БД", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при обновлении задачи в БД: %s", err))
		return
	}

	manager.Mng.Log.LogInfo("Задача успешно обновлена в БД и выполнена")
	// Отправляем пустой JSON в случае успеха
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("{}"))
}
