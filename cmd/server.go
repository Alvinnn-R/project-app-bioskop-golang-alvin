package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func APiserver(route *chi.Mux) {
	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", route); err != nil {
		log.Fatal("can't run service")
	}
}
