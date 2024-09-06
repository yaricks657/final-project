package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// структура для загрузки переменных окружения
type Environment struct {
	ServerPort       string `env:"TODO_PORT" env-default:"7540"`
	DatabaseFilePath string `env:"TODO_DBFILE" env-default:"./"`
	Password         string `env:"TODO_PASSWORD" env-default:"1234"`
	JWTSecret        string `env:"JWT_SECRET" env-default:"secret_key"`
}

// получить переменные окружения
func GetEnv() (Environment, error) {
	var env Environment

	if err := cleanenv.ReadEnv(&env); err != nil {
		return env, err
	}

	return env, nil
}
