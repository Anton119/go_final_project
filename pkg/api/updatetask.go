package api

import (
	"encoding/json"
	"fmt"
	"go_final_project/pkg/db"
	"go_final_project/pkg/nextdate"
	"net/http"
	"strconv"
	"time"
)

// updateTaskHandler обрабатывает PUT-запрос на обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseInt(req.ID, 10, 64); err != nil {
		writeJson(w, map[string]string{"error": "Invalid ID"}, http.StatusBadRequest)
		return
	}

	if req.ID == "" || req.Date == "" || req.Title == "" {
		writeJson(w, map[string]string{"error": "Missing required fields"}, http.StatusBadRequest)
		return
	}

	dateParsed, err := time.Parse("20060102", req.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "Invalid date format"}, http.StatusBadRequest)
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	if dateParsed.Before(today) {
		writeJson(w, map[string]string{"error": "Date cannot be in the past"}, http.StatusBadRequest)
		return
	}

	if req.Repeat != "" {
		_, err := nextdate.NextDate(time.Now(), req.Date, req.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "Invalid repeat rule"}, http.StatusBadRequest)
			return
		}
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err = db.GetDB().Exec(query, req.Date, req.Title, req.Comment, req.Repeat, req.ID)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Failed to update task: %v", err)}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{"id": req.ID}, http.StatusOK)
}
