package handlers

import (
	"dev-companion/internal/core/projects"
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

}

func writeDomainError(w http.ResponseWriter, err projects.DomainError) {

}
