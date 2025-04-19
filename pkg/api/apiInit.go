package api

import "net/http"

// Функция для инициализации всех API маршрутов
func InitAPI() {
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler) // Добавляем обработчик для nextdate
}
