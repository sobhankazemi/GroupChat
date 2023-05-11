package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sobhankazemi/GroupChat/History/dbrepo"
)

type Repository struct {
	DB dbrepo.Repository
}

func NewHandler(db dbrepo.Repository) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) HistoryAPI(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
	room_id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	pg := r.URL.Query().Get("page")
	page := 1
	if pg != "" {
		page, _ = strconv.Atoi(pg)
	}
	response, err := repo.DB.GetHistory(room_id, page)
	if err != nil {
		log.Println(err)
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	w.Write(jsonResponse)
}
