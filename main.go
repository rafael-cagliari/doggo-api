package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/rafael-cagliari/doggo-api/internal/api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Dog struct {
	Id    string `json:"id"`
	Breed string `json:"breed"`
	Name  string `json:"name"`
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.LoggerMiddleware)

	router.Use(middleware.AuthMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running!"))
	})

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running test!"))
	})

	router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		response := map[string]string{
			"message": "Usuário encontrado!",
			"userID":  id,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	router.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		var dog Dog

		err := json.NewDecoder(r.Body).Decode(&dog)
		if err != nil {
			http.Error(w, "JSON parsing error", http.StatusBadRequest)
			return
		}
		if dog.Breed == "persa" || dog.Breed == "sphynx" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("%s is not a dog breed", dog.Breed),
			})
			return
		}

		response := map[string]string{
			"message": "Usuário cadastrado com sucesso!",
			"id":      dog.Id,
			"name":    dog.Name,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe("0.0.0.0:8080", router)

}
