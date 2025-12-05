package handlers

import (
	"dev-companion/internal/core/projects"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ProjectsHandler struct {
	service projects.Service
}

func NewProjectHandler(svc projects.Service) *ProjectsHandler {
	newPjHan := ProjectsHandler{
		service: svc,
	}

	return &newPjHan
}

func (h *ProjectsHandler) RegisterRoutes(router *mux.Router) {

	router.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		h.ListProjects(w, r)
	}).Methods("GET")

	router.HandleFunc("/projects/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.GetProject(w, r)
	}).Methods("GET")

	router.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		h.CreateProject(w, r)
	}).Methods("POST")

	router.HandleFunc("/projects/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.UpdateProject(w, r)
	}).Methods("PUT")

	router.HandleFunc("/projects/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeleteProject(w, r)
	}).Methods("DELETE")

}

func (h *ProjectsHandler) ListProjects(w http.ResponseWriter, r *http.Request) {

}

func (h *ProjectsHandler) GetProject(w http.ResponseWriter, r *http.Request) {

}

func (h *ProjectsHandler) CreateProject(w http.ResponseWriter, r *http.Request) {

}

func (h *ProjectsHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {

}

func (h *ProjectsHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {

}

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Errore nella serializzazione JSON:", err)

		errorBody := map[string]string{
			"code":    "INTERNAL",
			"message": "Errore interno durante la serializzazione",
		}

		w.WriteHeader(http.StatusInternalServerError)

		fallback, err2 := json.MarshalIndent(errorBody, "", "  ")
		if err2 != nil {
			w.Write([]byte(`{"code":"INTERNAL","message":"Errore interno"}`))
			return
		}

		w.Write(fallback)
		return
	}

	w.WriteHeader(status)
	w.Write(j)
}

func writeDomainError(w http.ResponseWriter, err projects.DomainError) {
	code := err.Code

	statusHttp := 500

	switch {
	case code == "VALIDATION":
		statusHttp = 400
	case code == "NOT_FOUND":
		statusHttp = 404
	case code == "CONFLICT":
		statusHttp = 409
	case code == "INTERNAL":
		statusHttp = 500
	default:
		statusHttp = 500
	}

	type response struct {
		Code    string
		Message string
		Details map[string]string
	}

	var body response
	body.Code = err.Code
	body.Message = err.Message
	body.Details = err.Details

	writeJson(w, statusHttp, body)

}
