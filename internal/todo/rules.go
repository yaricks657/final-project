package todo

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yaricks657/final-project/internal/manager"
)

// NextDate вычисляет следующую дату на основе текущей даты, исходной даты и правила повторения.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	switch {
	case repeat == "":
		return "", fmt.Errorf("repeat rule is not specified")
	case strings.HasPrefix(repeat, "d "):
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid repeat rule")
		}
		parsedDate = parsedDate.AddDate(0, 0, days)
		for !parsedDate.After(now) {
			if parsedDate.Equal(now) {
				break
			}
			parsedDate = parsedDate.AddDate(0, 0, days)
		}
	case repeat == "y":
		parsedDate = parsedDate.AddDate(1, 0, 0)
		if parsedDate.Month() == time.February && parsedDate.Day() == 29 && !isLeapYear(parsedDate.Year()) {
			parsedDate = time.Date(parsedDate.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
		}
		for !parsedDate.After(now) {
			parsedDate = parsedDate.AddDate(1, 0, 0)
			if parsedDate.Month() == time.February && parsedDate.Day() == 29 && !isLeapYear(parsedDate.Year()) {
				parsedDate = time.Date(parsedDate.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
			}
		}
	case strings.HasPrefix(repeat, "w "):
		var daysOfWeek []int
		parts := strings.Split(strings.TrimPrefix(repeat, "w "), ",")
		for _, p := range parts {
			day, err := strconv.Atoi(p)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("invalid repeat rule")
			}
			daysOfWeek = append(daysOfWeek, day)
		}
		sort.Ints(daysOfWeek)
		for {
			weekday := int(parsedDate.Weekday())
			if weekday == 0 {
				weekday = 7 // Воскресенье равно 7
			}
			found := false
			for _, day := range daysOfWeek {
				if day >= weekday && (parsedDate.After(now) || parsedDate.Equal(now)) {
					if day == weekday && parsedDate.After(now) {
						found = true
						break
					}
					parsedDate = parsedDate.AddDate(0, 0, day-weekday)
					found = true
					break
				}
			}
			if found && parsedDate.After(now) {
				return parsedDate.Format("20060102"), nil
			}
			parsedDate = parsedDate.AddDate(0, 0, 7-int(parsedDate.Weekday())+daysOfWeek[0])
		}
	case strings.HasPrefix(repeat, "m "):
		parts := strings.Split(strings.TrimPrefix(repeat, "m "), " ")
		daysOfMonth := parseDays(parts[0])
		if daysOfMonth == nil {
			return "", fmt.Errorf("invalid repeat rule")
		}
		var months []int
		if len(parts) > 1 {
			months = parseMonths(parts[1])
			if months == nil {
				return "", fmt.Errorf("invalid repeat rule")
			}
		} else {
			for i := 1; i <= 12; i++ {
				months = append(months, i)
			}
		}
		sort.Ints(daysOfMonth)
		sort.Ints(months)
		for {
			month := int(parsedDate.Month())
			year := parsedDate.Year()
			var sortDates []time.Time
			for _, m := range months {
				lastDayOfTheMonth := getLastDayOfMonth(year, time.Month(m))
				if m >= month {
					for i, d := range daysOfMonth {
						if d <= lastDayOfTheMonth {
							newDate := calculateNewDate(year, m, d)
							if newDate.After(now) {
								sortDates = append(sortDates, newDate)
								if i == len(daysOfMonth)-1 {
									earliestDate := getEarliestDate(sortDates)
									return earliestDate, nil
								}
							}
						}
					}
				}
			}
			parsedDate = parsedDate.AddDate(0, 1, 0)
		}
	default:
		return "", fmt.Errorf("unsupported repeat rule: %s", repeat)
	}

	return parsedDate.Format("20060102"), nil
}

// получить ближайшую дату
func getEarliestDate(dates []time.Time) string {
	if len(dates) == 0 {
		return ""
	}

	earliest := dates[0] // Предполагаем, что первая дата самая ранняя

	// Проходим по остальным датам
	for _, date := range dates {
		if date.Before(earliest) { // Сравниваем с текущей минимальной
			earliest = date
		}
	}

	// Возвращаем самую раннюю дату в формате YYYYMMDD
	return earliest.Format("20060102")
}

// проверка на високосный год
func isLeapYear(year int) bool {
	if year%4 == 0 {
		if year%100 == 0 {
			return year%400 == 0
		}
		return true
	}
	return false
}

func parseDays(days string) []int {
	var result []int
	for _, day := range strings.Split(days, ",") {
		d, err := strconv.Atoi(day)
		if err != nil || d < -2 || d == 0 || d > 31 {
			return nil
		}
		result = append(result, d)
	}
	return result
}

func parseMonths(months string) []int {
	var result []int
	for _, month := range strings.Split(months, ",") {
		m, err := strconv.Atoi(month)
		if err != nil || m < 1 || m > 12 {
			return nil
		}
		result = append(result, m)
	}
	return result
}

func calculateNewDate(year, month, day int) time.Time {
	lastDay := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
	if day > 0 {
		if day > lastDay {
			day = lastDay
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	} else if day == -1 {
		return time.Date(year, time.Month(month), lastDay, 0, 0, 0, 0, time.UTC)
	} else if day == -2 {
		return time.Date(year, time.Month(month), lastDay-1, 0, 0, 0, 0, time.UTC)
	}
	return time.Time{}
}

// получить последний день месяца. Переделать под единоразовый вызов надо потом
func getLastDayOfMonth(year int, month time.Month) int {
	// Создаем временную метку на первый день следующего месяца
	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)

	// Вычисляем последний день текущего месяца как день перед первым днем следующего месяца
	lastDayOfMonth := firstDayNextMonth.Add(-24 * time.Hour)

	return lastDayOfMonth.Day()
}

// обработчик для /api/nextdate
func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	manager.Mng.Log.LogInfo("поступил запрос на получение задачи (HandleNextDate)", r.RequestURI)

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		manager.Mng.Log.LogError("Invalid now format, expected YYYYMMDD", err)
		http.Error(w, "Invalid now format, expected YYYYMMDD", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `%s`, nextDate)
	manager.Mng.Log.LogInfo("успешная обработка запроса api/nextdate (HandleNextDate)")
}
