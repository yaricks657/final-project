package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yaricks657/final-project/internal/manager"
)

type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// создать БД
func CreateDB(mng *manager.Manager) error {
	/* appPath, err := os.Executable()
	if err != nil {
		return err
	}
	mng.Log.LogInfo(appPath) */
	dbFile := filepath.Join(filepath.Dir("./"), "scheduler.db")
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	if !install {
		return nil
	}

	// создание БД
	mng.Log.LogWarn("База данных отсутствует и будет создана новая")

	// открытие БД
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		mng.Log.LogError("Ошибка при открытии БД", err)
		return err
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT CHECK(LENGTH(repeat) <= 128)
	);
	CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
	`

	// Выполнение запроса
	_, err = db.Exec(sqlStmt)
	if err != nil {
		mng.Log.LogError("Ошибка при запросе создания БД", err)
		return err
	}

	return nil
}

// Добавить задачу в БД
func AddTask(t *Task, mng *manager.Manager) (string, error) {
	// открытие БД
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		mng.Log.LogError("Ошибка при открытии БД", err)
		return "", err
	}
	defer db.Close()
	// Подготовка SQL-запроса для вставки данных
	insertTaskSQL := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	// Выполнение запроса
	statement, err := db.Prepare(insertTaskSQL)
	if err != nil {
		return "", err
	}
	defer statement.Close()

	// Вставка данных и получение ID
	result, err := statement.Exec(t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", err
	}

	// Получение ID вставленной записи
	taskID, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(taskID)), nil
}

// Получить все задачи из БД
func GetAllTasks(mng *manager.Manager) ([]Task, error) {
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Запрос для получения всех задач, отсортированных по дате
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50"

	// Выполнение SQL-запроса
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Сбор результатов в слайс задач
	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// проверка на пустой слайс
	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

// Получить задачи по поиску из БД
func GetSearchedTasks(mng *manager.Manager, search string) ([]Task, error) {
	// Подключение к базе данных
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Подготовка базового SQL-запроса
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE 1=1"
	var args []interface{}

	// Обработка параметра search
	if search != "" {
		// Проверка на соответствие формату даты 02.01.2006
		if date, err := time.Parse("02.01.2006", search); err == nil {
			// Преобразование даты в формат 20060102
			searchDate := date.Format("20060102")
			query += " AND date = ?"
			args = append(args, searchDate)
		} else {
			// Поиск подстроки в полях title и comment без изменения регистра
			searchPattern := "%" + search + "%"
			query += " AND (title LIKE ? OR comment LIKE ?)"
			args = append(args, searchPattern, searchPattern)
		}
	}

	// Завершение SQL-запроса и сортировка по дате
	query += " ORDER BY date LIMIT 50"

	// Выполнение SQL-запроса
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Сбор результатов в слайс задач
	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Если задач нет, возвращаем пустой список
	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

// Получить задачу по id
func GetTask(mng *manager.Manager, id string) (Task, error) {
	// Подключение к базе данных
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	// Подготовка SQL-запроса для поиска задачи по ID
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ? LIMIT 1"

	// Выполнение SQL-запроса
	row := db.QueryRow(query, id)

	// Сбор результата
	var task Task
	err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("Задача с id %s не найдена", id) // Если задача не найдена, возвращаем nil
		}
		return Task{}, err // В случае ошибки возвращаем ошибку
	}
	return task, nil
}

// изменение существующей задачи
func ChangeTask(t *Task, mng *manager.Manager) error {
	// Подключение к базе данных
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Проверяем существует ли задача с указанным ID
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM scheduler WHERE id = ?)", t.Id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("задача с ID %s не найдена", t.Id)
	}

	// Обновление задачи в базе данных
	_, err = db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		t.Date, t.Title, t.Comment, t.Repeat, t.Id)
	if err != nil {
		return err
	}
	return nil
}

// Удалить задачу из БД
func DeleteTask(mng *manager.Manager, id string) error {
	// Подключение к базе данных
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Проверяем существует ли задача с указанным ID
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM scheduler WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Задача с id %s не найдена", id)
	}

	// Удаление задачи из базы данных
	_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
