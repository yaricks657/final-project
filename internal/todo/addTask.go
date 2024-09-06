package todo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/yaricks657/final-project/internal/database"
	"github.com/yaricks657/final-project/internal/manager"
)

// добавить задачу в БД
func AddTask(w http.ResponseWriter, r *http.Request) {
	manager.Mng.Log.LogInfo("поступил запрос на добавление задачи (AddTask)")
	// проверка на нужный метод
	if r.Method != http.MethodPost {
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

	// проверка на наличие обязательных полей полей, если в дальнейшем прибавятся
	ok, err := task.checkRequiredFields()
	if !ok {
		manager.Mng.Log.LogError("Отсутствуют обязательные поля", err)
		sendErrorResponse(w, fmt.Sprintf("Отсутствуют обязательные поля: %s", err))
		return
	}

	// проверка на корректность даты
	ok, err = task.isDateValid()

	if !ok {
		manager.Mng.Log.LogError("Некорректный формат даты:", err)
		sendErrorResponse(w, fmt.Sprintf("Некорректный формат даты %s", err))
		return
	}

	// применение правила для даты, если она раньше сегодняшнего дня
	today := time.Now().Truncate(24 * time.Hour)
	Ttime, _ := time.Parse("20060102", task.Date)

	if Ttime.Before(today) {
		manager.Mng.Log.LogWarn("Дата раньше сегодняшнего числа")
		ok, err = task.isRuleValid(today)
		if !ok {
			manager.Mng.Log.LogError("Ошибка при применении правила повторения", err)
			sendErrorResponse(w, fmt.Sprintf("Ошибка при применении правила повторения %s", err))
			return
		}
	}

	// запись задачи в БД
	addTask := database.Task{
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
	id, err := database.AddTask(&addTask, &manager.Mng)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при записи задачи в БД", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при записи задачи в БД: %s", err))
		return
	}

	// отправка успешного ответа с id
	manager.Mng.Log.LogInfo("Задача успешно записана в БД")
	sendSuccessResponse(w, id)

}

// проверка правила
func (t *Task) isRuleValid(now time.Time) (bool, error) {
	if t.Repeat == "" {
		manager.Mng.Log.LogWarn("Правило отсутствует. Будет проставлено сегодняшнее число")
		t.Date = time.Now().Format("20060102")
		return true, nil
	}

	newDate, err := NextDate(now, t.Date, t.Repeat)
	if err != nil {
		return false, err
	}
	t.Date = newDate

	return true, nil
}

// проверка даты на валидность
func (t *Task) isDateValid() (bool, error) {
	if t.Date == "" {
		t.Date = time.Now().Format("20060102")
		manager.Mng.Log.LogWarn("Дата отсутствует. Будет проставлено сегодняшнее число")
		return true, nil
	}

	_, err := time.Parse("20060102", t.Date)
	if err != nil {
		return false, fmt.Errorf("Некорректный формат даты. Ожидался YYYYMMDD, получили %s", t.Date)
	}

	return true, nil
}

// проверить наличие обязательных полей
func (t *Task) checkRequiredFields() (bool, error) {
	var missingFields []string

	if t.Title == "" {
		missingFields = append(missingFields, "Title")
	}

	if len(missingFields) > 0 {
		return false, fmt.Errorf("отсутствуют обязательные поля: %s", strings.Join(missingFields, ", "))
	}
	return true, nil
}

// отправка на клиент ответа об ошибке
func sendErrorResponse(w http.ResponseWriter, errorMsg string) {
	response := ResponseAddTask{
		Error: errorMsg,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при декодировании (sendErrorResponse)", err)
	}
}

// отправка на клиент ответа об успехе
func sendSuccessResponse(w http.ResponseWriter, id string) {
	response := ResponseAddTask{
		ID: id,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при декодировании (sendSuccessResponse)", err)
	}
}
