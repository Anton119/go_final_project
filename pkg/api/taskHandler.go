package api

import (
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Обработка добавления задачи
		addTaskHandler(w, r)
	default:
		// Если метод не поддерживается
		writeJson(w, map[string]string{
			"error": "Метод не поддерживается",
		}, http.StatusMethodNotAllowed)
	}
}
