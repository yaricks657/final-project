package manager

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/yaricks657/final-project/config"
	"github.com/yaricks657/final-project/pkg/logger"
)

// методы для manager
type IManager interface {
	//Start()
	SetHandlers(r *chi.Mux)
}

// структура-контейнер для управления проектом
type Manager struct {
	//стартовый конфиг
	Cnf config.Environment

	// логирование zerolog
	Log *logger.Logger

	// бд
	Db *sql.DB
}

// структура для создания manager
type CreateConfig struct {
	Cnf config.Environment
	Log *logger.Logger
}

// Глобальная переменная с начальной конфигурацией
var Mng Manager

// старт manager
func New(config CreateConfig) (*Manager, error) {
	// добавить проверку требуемых полей методом для CreateConfig
	err := config.checkRequiredFields()
	if err != nil {
		return nil, err
	}

	Mng = Manager{
		Log: config.Log,
		Cnf: config.Cnf,
	}

	return &Mng, nil
}

// проверка наличия обязательных полей
func (cc *CreateConfig) checkRequiredFields() error {
	var missingFields []string

	if cc.Cnf.ServerPort == "" {
		missingFields = append(missingFields, "TODO_PORT")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("отсутствуют обязательные поля: %s", strings.Join(missingFields, ", "))
	}
	return nil
}
