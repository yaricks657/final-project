package todo

import "github.com/yaricks657/final-project/internal/database"

// запрос на добавление задачи addTask и запросе задачи
type Task struct {
	Id      string `json:id,omitempty`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// ответ при добавлении задачи
type ResponseAddTask struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// ответ при запросе всех задач
type ResponseGetAllTasks struct {
	Tasks []database.Task `json:"tasks"`
}
