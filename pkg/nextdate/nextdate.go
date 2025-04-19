package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "20060102"

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func NextDate(dstart time.Time, repeat string, nowStr string) (string, error) {
	// Парсим now
	_, err := time.Parse(TimeFormat, nowStr)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга now: %w", err)
	}

	// Если repeat пустой, просто возвращаем dstart
	if len(repeat) == 0 {
		return dstart.Format(TimeFormat), nil
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "y":
		// Добавляем год
		dstart = dstart.AddDate(1, 0, 0)
		// Проверяем на високосный год
		if dstart.Month() == 2 && dstart.Day() == 29 && !isLeapYear(dstart.Year()) {
			dstart = time.Date(dstart.Year(), 3, 1, 0, 0, 0, 0, dstart.Location())
		}
		return dstart.Format(TimeFormat), nil

	case "d":
		// Добавляем дни
		if len(parts) < 2 {
			return "", fmt.Errorf("нужен аргумент после d")
		}
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("ошибка преобразования в число: %w", err)
		}
		if num <= 0 || num > 400 {
			return "", fmt.Errorf("некорректное количество дней")
		}
		dstart = dstart.AddDate(0, 0, num)
		return dstart.Format(TimeFormat), nil

	case "m", "w":
		// Заглушка для неверных типов повторений
		return "", fmt.Errorf("некорректный запрос")

	default:
		return "", fmt.Errorf("неизвестное правило: %s", rule)
	}
}
