package main

import (
	"dev-companion/internal/config"
	"dev-companion/internal/core/projects"
	"dev-companion/internal/http/handlers"
	"dev-companion/internal/http/middleware"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	var path = "config/config.json"

	val, err := config.LoadConfig(path)
	if err != nil {
		fmt.Println("Errore durante il caricamento della configurazione: ", err)
		os.Exit(1)
	}

	dir, err1 := projects.NewProjectRepository(val.DataDir)
	if err1 != nil {
		fmt.Println("Errore durante il caricamento del repository: ", err1)
		os.Exit(1)
	}

	pjService := projects.CreatePjService(dir)

	pjHandler := handlers.NewProjectHandler(pjService)

	router := mux.NewRouter()

	authMw := middleware.NewAuthMiddleware(val.AuthToken)
	router.Use(authMw)

	pjHandler.RegisterRoutes(router)

	str := val.ListenAddress + ":" + strconv.Itoa(val.Port)
	fmt.Println("Server Dev Companion in ascolto su: " + str)

	err2 := http.ListenAndServe(str, router)
	if err2 != nil {
		fmt.Println("Errore durante l'avvio del server: ", err2)
		os.Exit(-1)
	}

}
