package handlers

import (
	"dev-companion/internal/core/events"
	"net/http"

	"github.com/gorilla/mux"
)

type EventsHandler struct {
	service events.Service
}

func NewEventsHandler(svc events.Service) *EventsHandler {
	newRunHan := EventsHandler{
		service: svc,
	}

	return &newRunHan
}

func (h *EventsHandler) RegisterEventsRoutes(router *mux.Router) {

	router.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		h.ListEvents(w, r)
	}).Methods("GET")

	router.HandleFunc("/events/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.GetEvent(w, r)
	}).Methods("GET")

}

func (h *EventsHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ListEvents(events.EventFilter{})
	if err != nil {
		if domainErr, ok := err.(events.DomainError); ok {
			writeDomainEventError(w, domainErr)
			return
		}

		internalErr := events.NewInternalError(
			"errore interno durante ListEvents: " + err.Error(),
		)
		writeDomainEventError(w, internalErr)
		return
	}

	writeJson(w, http.StatusOK, res)
}

func (h *EventsHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	event, err := h.service.GetEvent(id)
	if err != nil {
		if domainErr, ok := err.(events.DomainError); ok {
			writeDomainEventError(w, domainErr)
			return
		}

		internalErr := events.NewInternalError(
			"errore interno durante GetEvent: " + err.Error(),
		)
		writeDomainEventError(w, internalErr)
		return

	}
	writeJson(w, 200, event)

}

func writeDomainEventError(w http.ResponseWriter, err events.DomainError) {
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
