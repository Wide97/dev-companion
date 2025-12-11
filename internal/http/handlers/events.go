package handlers

import (
	"dev-companion/internal/core/events"
	"dev-companion/internal/utility"
	"net/http"
	"time"

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

	router.HandleFunc("/events/today", func(w http.ResponseWriter, r *http.Request) {
		h.ListTodayEvents(w, r)
	}).Methods("GET")

	router.HandleFunc("/events/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.GetEvent(w, r)
	}).Methods("GET")

}

func (h *EventsHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	projectIdStr := query.Get("projectId")
	typeStr := query.Get("type")
	fromStr := query.Get("from")
	toStr := query.Get("to")

	convFromStr, err := utility.Parser(fromStr)
	if err != nil {
		ev := events.NewValidationError(map[string]string{
			"from": "formato data non valido, usa RFC3339",
		})
		writeDomainEventError(w, ev)
		return
	}

	convToStr, err1 := utility.Parser(toStr)
	if err1 != nil {
		ev1 := events.NewValidationError(map[string]string{
			"to": "formato data non valido, usa RFC3339",
		})
		writeDomainEventError(w, ev1)
		return
	}

	filt := events.EventFilter{
		ProjectId: projectIdStr,
		Type:      typeStr,
		From:      convFromStr,
		To:        convToStr,
	}

	res, err := h.service.ListEvents(filt)
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

func (h *EventsHandler) ListTodayEvents(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	evNow, err := h.service.ListTodayEvents(now)
	if err != nil {
		if domainErr, ok := err.(events.DomainError); ok {
			writeDomainEventError(w, domainErr)
			return
		}
		internalErr := events.NewInternalError(
			"errore interno durante ListTodayEvents: " + err.Error(),
		)
		writeDomainEventError(w, internalErr)
		return
	}

	writeJson(w, 200, evNow)

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
