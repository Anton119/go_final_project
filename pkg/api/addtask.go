package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/nextdate"
)

// Функция проверки и корректировки даты
func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(nextdate.TimeFormat) // если дата не указана — текущая
	}

	t, err := time.Parse(nextdate.TimeFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %v", err)
	}

	// Обрезаем время у обеих дат, сравниваем только по дню
	nowDateOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	taskDateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if !afterNow(nowDateOnly, taskDateOnly) { // если дата задачи не в будущем
		if len(task.Repeat) == 0 {
			task.Date = now.Format(nextdate.TimeFormat) // без повторений — текущая дата
		} else {
			next, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка в правиле повторения: %v", err)
			}
			task.Date = next // следующая подходящая дата
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка чтения тела запроса: %v", err)}, http.StatusBadRequest)
		return
	}
	log.Printf("Полученные данные: %s", body)

	// Декодирование JSON
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка декодирования JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Поле 'title' обязательно"}, http.StatusBadRequest)
		return
	}

	// Проверка даты
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Добавление задачи в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка добавления задачи в базу данных: %v", err)}, http.StatusInternalServerError)
		return
	}

	// Ответ с id задачи
	writeJson(w, map[string]interface{}{"id": id}, http.StatusOK)
}

// утилита для json ответов
func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("ошибка сериализации JSON: %v", err)
		http.Error(w, fmt.Sprintf("ошибка сериализации ответа"), http.StatusInternalServerError)
	}
}

// Функция для проверки, что дата позже текущего времени
func afterNow(now, date time.Time) bool {
	return date.After(now)
}
