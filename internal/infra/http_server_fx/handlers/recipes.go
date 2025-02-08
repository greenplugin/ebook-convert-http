package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
)

type Recipes struct {
}

func (h *Recipes) Register(mux *http.ServeMux) {
	mux.HandleFunc("/recipes", h.GetParams)
}

func NewRecipes() *Recipes {
	return &Recipes{}
}

func (h *Recipes) GetParams(w http.ResponseWriter, r *http.Request) {
	var stdout bytes.Buffer
	cmd := exec.Command("ebook-convert", "--list-recipes")
	cmd.Stdout = &stdout
	cmd.Run()

	recipes := bytes.Split(stdout.Bytes(), []byte("\n"))

	recipes = recipes[1:]

	jsonRecipes := make([]string, 0)
	for _, recipe := range recipes {
		if len(recipe) > 0 {
			jsonRecipes = append(jsonRecipes, strings.TrimSpace(string(recipe)))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jsonRecipes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
