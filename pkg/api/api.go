package api

import (
	"fmt"
	"go_final_project/pkg/nextdate"
	"log"
	"net/http"
	"time"
)

// Обработчик для вычисления следующей даты
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из запроса (из URL для GET-запроса)
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	// Логируем параметры для отладки
	log.Printf("Received parameters: now=%s, date=%s, repeat=%s", nowStr, dateStr, repeatStr)

	// Проверяем, что параметры не пустые
	if nowStr == "" || dateStr == "" || repeatStr == "" {
		http.Error(w, "Ошибка: отсутствуют необходимые параметры (now, date, repeat)", http.StatusBadRequest)
		return
	}

	// Преобразуем nowStr в формат time
	now, err := time.Parse(nextdate.TimeFormat, nowStr)
	if err != nil {
		http.Error(w, "Ошибка парсинга параметра now", http.StatusBadRequest)
		return
	}

	// Передаем now как time.Time в NextDate
	nextDate, err := nextdate.NextDate(now, dateStr, repeatStr)
	if err != nil {
		// Логируем ошибку и отправляем сообщение с деталями ошибки
		log.Printf("Error calculating next date: %v", err)
		http.Error(w, fmt.Sprintf("Ошибка вычисления следующей даты: %v", err), http.StatusInternalServerError)
		return
	}

	// Логируем результат
	log.Printf("Next date: %s", nextDate)
	// Отправляем ответ
	fmt.Fprintf(w, "%s", nextDate)
}
