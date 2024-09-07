package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yaricks657/final-project/internal/database"
	"github.com/yaricks657/final-project/internal/manager"
)

// Изменить существующую задачу
func ChangeTask(w http.ResponseWriter, r *http.Request) {
	// проверка на нужный метод
	if r.Method != http.MethodPut {
		manager.Mng.Log.LogWarn("Некорректный метод запроса")
		sendErrorResponse(w, "Некорректный метод запроса")
		return
	}

	// чтение тела запроса в слайс байт
	body, err := io.ReadAll(r.Body)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при чтении тела запроса", err)
		sendErrorResponse(w, "Ошибка при чтении тела запроса")
		return
	}
	defer r.Body.Close()

	// Распаковка ответа от клиента
	var task Task
	if err = json.Unmarshal(body, &task); err != nil {
		manager.Mng.Log.LogError("Bad request", err)
		sendErrorResponse(w, "Bad request")
		return
	}

	if task.Id == "" {
		manager.Mng.Log.LogError("Отсутствуют id", err)
		sendErrorResponse(w, fmt.Sprintf("Отсутствует id: %s", err))
		return
	}

	// проверка на наличие обязательных полей полей, если в дальнейшем прибавятся
	ok, err := task.checkRequiredFields()
	if !ok {
		manager.Mng.Log.LogError("Отсутствуют обязательные поля", err)
		sendErrorResponse(w, fmt.Sprintf("Отсутствуют обязательные поля: %s", err))
		return
	}

	// проверка на корректность даты
	err = task.isDateValid()

	if err != nil {
		manager.Mng.Log.LogError("Некорректный формат даты:", err)
		sendErrorResponse(w, fmt.Sprintf("Некорректный формат даты %s", err))
		return
	}

	// применение правила для даты, если она раньше сегодняшнего дня
	today := time.Now().Truncate(24 * time.Hour)
	parsedDate, _ := time.Parse(DateFormat, task.Date)

	if parsedDate.Before(today) {
		manager.Mng.Log.LogWarn("Дата раньше сегодняшнего числа")
		err = task.isRuleValid(today)
		if err != nil {
			manager.Mng.Log.LogError("Ошибка при применении правила повторения", err)
			sendErrorResponse(w, fmt.Sprintf("Ошибка при применении правила повторения %s", err))
			return
		}
	}

	// запись задачи в БД
	updateTask := database.Task{
		Id:      task.Id,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
	err = database.ChangeTask(&updateTask, &manager.Mng)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обновлении задачи в БД", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при обновлении задачи в БД: %s", err))
		return
	}

	// Отправляем пустой JSON в случае успеха
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при отправке ответа: ", err)
	}

}
