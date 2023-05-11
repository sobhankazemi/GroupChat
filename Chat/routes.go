package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetRoutes() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/{id}" ,handler.ServeChatRooms)
	return mux
}
