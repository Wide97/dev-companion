package main

import (
	"dev-companion/internal/config"
	"dev-companion/internal/core/events"
	"dev-companion/internal/core/projects"
	"dev-companion/internal/core/runs"
	"dev-companion/internal/http/handlers"
	"dev-companion/internal/http/middleware"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	var path = "internal/config/config.json"

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

	dir1, err2 := runs.NewRunRepository(val.DataDir)
	if err2 != nil {
		fmt.Println("Errore durante il caricamento della configurazione: ", err2)
		os.Exit(1)
	}

	dir2, err3 := events.NewEventsRepository(val.DataDir)
	if err3 != nil {
		fmt.Println("Errore durante il caricamento della configurazione: ", err3)
		os.Exit(1)
	}

	pjService := projects.CreatePjService(dir)

	pjHandler := handlers.NewProjectHandler(pjService)

	runsService := runs.CreateRunSerivce(dir1)

	runsHandler := handlers.NewRunsHandler(runsService)

	eventsService := events.CreateEventsService(dir2)

	eventsHandler := handlers.NewEventsHandler(eventsService)

	router := mux.NewRouter()

	recMw := middleware.NewRecoveryMiddleware()
	router.Use(recMw)

	loggingMw := middleware.NewLoggingMiddleware()
	router.Use(loggingMw)

	authMw := middleware.NewAuthMiddleware(val.AuthToken)
	router.Use(authMw)

	pjHandler.RegisterRoutes(router)

	runsHandler.RegisterRunsRoutes(router)

	eventsHandler.RegisterEventsRoutes(router)

	str := val.ListenAddress + ":" + strconv.Itoa(val.Port)
	fmt.Println("Server Dev Companion in ascolto su: " + str)

	err4 := http.ListenAndServe(str, router)
	if err4 != nil {
		fmt.Println("Errore durante l'avvio del server: ", err2)
		os.Exit(-1)
	}

}

//Comando per lanciare: go run ./cmd/dev-companion
