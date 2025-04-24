package api

import (
	"net/http"
)

// Обработчик для /api/task, который определяет действие в зависимости от HTTP-метода
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r) // Добавление новой задачи
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
