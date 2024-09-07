package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yaricks657/final-project/internal/manager"
	"github.com/yaricks657/final-project/internal/todo"
)

// обработчик для /api/nextdate
func HandleNextDate(w http.ResponseWriter, r *http.Request) {

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		manager.Mng.Log.LogError("Invalid now format, expected YYYYMMDD", err)
		http.Error(w, "Invalid now format, expected YYYYMMDD", http.StatusBadRequest)
		return
	}

	nextDate, err := todo.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `%s`, nextDate)
}
