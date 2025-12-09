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
	res, err := h.service.ListProjects()
	if err != nil {
		if domainErr, ok := err.(*projects.DomainError); ok {
			writeDomainError(w, *domainErr)
			return
		}

		internalErr := projects.NewInternalError(
			"errore interno durante ListProjects: " + err.Error(),
		)
		writeDomainError(w, internalErr)
		return
	}

	writeJson(w, http.StatusOK, res)
}

func (h *ProjectsHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	project, err := h.service.GetProject(id)
	if err != nil {
		if domainErr, ok := err.(*projects.DomainError); ok {
			writeDomainError(w, *domainErr)
			return
		}

		internalErr := projects.NewInternalError(
			"errore interno durante GetProject: " + err.Error(),
		)
		writeDomainError(w, internalErr)
		return

	}
	writeJson(w, 200, project)

}

func (h *ProjectsHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var input projects.CreateProjectInput

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		validationErr := projects.NewValidationError(map[string]string{
			"body": "JSON malformato o mancante",
		})
		writeDomainError(w, validationErr)
		return
	}

	createdProject, err := h.service.CreateProject(input)
	if err != nil {
		if domainErr, ok := err.(*projects.DomainError); ok {
			writeDomainError(w, *domainErr)
			return
		}

		internalErr := projects.NewInternalError(
			"errore interno durante CreateProject: " + err.Error(),
		)
		writeDomainError(w, internalErr)
		return
	}

	writeJson(w, http.StatusCreated, createdProject)

}

func (h *ProjectsHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var upd projects.UpdateProjectInput

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&upd)
	if err != nil {
		validationErr := projects.NewValidationError(map[string]string{
			"body": "JSON malformato o mancante",
		})
		writeDomainError(w, validationErr)
		return

	}

	updatePj, err1 := h.service.UpdateProject(id, upd)
	if err1 != nil {
		if domainErr, ok := err1.(*projects.DomainError); ok {
			writeDomainError(w, *domainErr)
			return
		}

		internalErr := projects.NewInternalError(
			"errore interno durante UpdateProject: " + err1.Error(),
		)

		writeDomainError(w, internalErr)
		return
	}

	writeJson(w, 200, updatePj)

}

func (h *ProjectsHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	del := h.service.DeleteProject(id)
	if del != nil {
		if domainErr, ok := del.(*projects.DomainError); ok {
			writeDomainError(w, *domainErr)
			return
		}

		internalErr := projects.NewInternalError(
			"errore interno durante DeleteProject: " + del.Error(),
		)

		writeDomainError(w, internalErr)
		return

	}

	w.WriteHeader(http.StatusNoContent)

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

//cuiaooooooo
