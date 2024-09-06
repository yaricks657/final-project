package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/yaricks657/final-project/config"
	"github.com/yaricks657/final-project/internal/auth"
	"github.com/yaricks657/final-project/internal/database"
	"github.com/yaricks657/final-project/internal/manager"
	"github.com/yaricks657/final-project/internal/todo"
	"github.com/yaricks657/final-project/pkg/logger"
)

func main() {

	// создание логгера
	logFilePath := "./app.log"
	logger, err := logger.New(logFilePath)
	if err != nil {
		log.Fatal("Ошибка при создании логгера (main)", err)
		os.Exit(1)
	}

	// загрузка конфига
	config, err := config.GetEnv()
	if err != nil {
		logger.LogError("Ошибка при загрузке конфига (main)", err)
		os.Exit(1)
	}
	logger.LogInfo("Переменные окружения получены", fmt.Sprintln(config))

	// создание контейнера общего
	createConfig := manager.CreateConfig{
		Cnf: config,
		Log: logger,
	}

	// создаем менеджер для дальнейшей работы с ним
	Mng, err := manager.New(createConfig)
	if err != nil {
		logger.LogError("Ошибка при регистрации manager (main)", err)
		os.Exit(1)
	}
	logger.LogInfo("manager зарегистрирован", fmt.Sprintln(&Mng))

	// подключение БД
	err = database.CreateDB(Mng)
	if err != nil {
		os.Exit(1)
	}

	logger.LogInfo("БД подключена успешно")

	// создание и запуск сервера
	router := chi.NewRouter()

	/* Установка обработчиков */
	// web-часть
	fileServer := http.FileServer(http.Dir("./web"))
	router.Handle("/*", http.StripPrefix("/", fileServer))
	// проыерить следующую дату
	router.Get("/api/nextdate", auth.Auth(todo.HandleNextDate))
	// добавить задачу
	router.Post("/api/task", auth.Auth(todo.AddTask))
	// получить все задачи
	router.Get("/api/tasks", auth.Auth(todo.GetAllTasks))
	// получить задачи по поиску
	//	router.Get("/api/tasks/{search}", todo.GetSearchedTasks)
	// получить задачу
	router.Get("/api/task", auth.Auth(todo.GetTask))
	// изменить задачу
	router.Put("/api/task", auth.Auth(todo.ChangeTask))
	// выполнить задачу
	router.Post("/api/task/done", auth.Auth(todo.DoneTask))
	// удаление задачи
	router.Delete("/api/task", auth.Auth(todo.DeleteTask))
	// авторизация
	router.Post("/api/signin", auth.SignIn)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), router); err != nil {
		logger.LogError("Ошибка запуска сервера (main)", err)
		os.Exit(1)
	}
	logger.LogInfo("Сервер запущен на порту", config.ServerPort)

}
