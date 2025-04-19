package api

import (
	"go_final_project/pkg/nextdate"
	"log"
	"net/http"
	"time"
)

func handleNextDate(w http.ResponseWriter, r *http.Request) {
	// Логируем запрос
	log.Println("Received request:", r.URL)

	// Получаем параметры из запроса
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	nowStr := r.URL.Query().Get("now")

	// Логируем полученные параметры
	log.Println("Parameters - date:", dateStr, "repeat:", repeat, "now:", nowStr)

	// Преобразуем строку в time.Time
	_, err := time.Parse("20060102", dateStr) // Формат: YYYYMMDD
	if err != nil {
		log.Println("Error parsing date:", err) // Логируем ошибку
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	// Аналогично можно обработать параметр "now", если он тоже передаётся как строка
	var now time.Time
	if nowStr != "" {
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			log.Println("Error parsing now:", err) // Логируем ошибку
			http.Error(w, "invalid now format", http.StatusBadRequest)
			return
		}
	} else {
		now = time.Now() // Если параметр now пустой, используем текущее время
	}

	// Преобразуем now обратно в строку
	nowFormatted := now.Format("20060102") // Преобразуем time.Time в строку в формате YYYYMMDD

	// Логируем информацию перед вызовом функции NextDate
	log.Println("Calling NextDate function with date:", dateStr, "repeat:", repeat, "now:", nowFormatted)

	// Вызываем логику из пакета nextdate
	result, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		log.Println("Error in NextDate:", err) // Логируем ошибку
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Логируем результат
	log.Println("NextDate result:", result)

	// Отдаём результат как простой текст
	w.Write([]byte(result))
}
