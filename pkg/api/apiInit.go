package api

import "net/http"

// Функция для инициализации всех API маршрутов
func InitAPI() {
	http.HandleFunc("/api/nextdate", nextDayHandler) // Обработчик для вычисления следующей даты
	http.HandleFunc("/api/task", taskHandler)        // Обработчик для работы с задачами
}
