package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "20060102"

// Функция для проверки, если дата больше now
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// Функция для вычисления следующей даты
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Преобразуем dstart в time.Time
	date, err := time.Parse(TimeFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	// Разбиваем repeat на составляющие
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("неправильный формат repeat: %v", repeat)
	}

	switch parts[0] {
	case "y": // если правило - ежегодно
		// Добавляем 1 год
		date = date.AddDate(1, 0, 0)
		// Если дата всё ещё меньше current, добавляем ещё 1 год
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(TimeFormat), nil

	case "d": // если правило - дни
		if len(parts) < 2 {
			return "", fmt.Errorf("не указан интервал для d")
		}

		// Парсим интервал в днях
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверный формат интервала дней: %v", parts[1])
		}
		// Проверка на допустимый диапазон
		if days < 1 || days > 400 {
			return "", fmt.Errorf("число дней должно быть от 1 до 400")
		}

		// Добавляем дни, пока дата не станет больше now
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(TimeFormat), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат повторения: %s", repeat)
	}
}
