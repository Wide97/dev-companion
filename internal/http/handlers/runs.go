package handlers

import (
	"dev-companion/internal/core/runs"
	"net/http"

	"github.com/gorilla/mux"
)

type RunsHandler struct {
	service runs.Service
}

func NewRunsHandler(svc runs.Service) *RunsHandler {
	newRunHan := RunsHandler{
		service: svc,
	}

	return &newRunHan
}

func (h *RunsHandler) RegisterRunsRoutes(router *mux.Router) {

	router.HandleFunc("/runs", func(w http.ResponseWriter, r *http.Request) {
		h.ListRuns(w, r)
	}).Methods("GET")

	router.HandleFunc("/runs/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.GetRun(w, r)
	}).Methods("GET")

}

func (h *RunsHandler) ListRuns(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListRuns(runs.RunFilter{})
	if err != nil {
		if domainErr, ok := err.(runs.DomainError); ok {
			writeDomainRunError(w, domainErr)
			return
		}

		internalErr := runs.NewInternalError(
			"errore interno durante ListRuns: " + err.Error(),
		)
		writeDomainRunError(w, internalErr)
		return
	}

	writeJson(w, http.StatusOK, res)
}

func (h *RunsHandler) GetRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	project, err := h.service.GetRun(id)
	if err != nil {
		if domainErr, ok := err.(runs.DomainError); ok {
			writeDomainRunError(w, domainErr)
			return
		}

		internalErr := runs.NewInternalError(
			"errore interno durante GetRun: " + err.Error(),
		)
		writeDomainRunError(w, internalErr)
		return

	}
	writeJson(w, 200, project)

}

func writeDomainRunError(w http.ResponseWriter, err runs.DomainError) {
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
