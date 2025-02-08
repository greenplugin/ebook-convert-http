package handlers

import "net/http"

type Health struct {
}

func (h *Health) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func NewHealth() *Health {
	return &Health{}
}
