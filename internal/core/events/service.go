package events

import (
	"strings"
	"time"
)

type DomainError struct {
	Code    string
	Message string
	Details map[string]string
}

type Service struct {
	Repo *EventRepository
}

const ERROR_VALIDATION = "VALIDATION"
const ERROR_NOT_FOUND = "NOT_FOUND"
const ERROR_CONFLICT = "CONFLICT"
const ERROR_INTERNAL = "INTERNAL"

type EventFilter struct {
	ProjectId string
	Type      string
	From      *time.Time
	To        *time.Time
}

func (s *Service) ListEvents(filter EventFilter) ([]EventModel, error) {
	errors := map[string]string{}
	if filter.Type != "" && filter.Type != EventPjBuildFailed && filter.Type != EventPjBuildStarted && filter.Type != EventPjBuildSuccess && filter.Type != EventRunFailed && filter.Type != EventRunStarted && filter.Type != EventRunSuccess {
		errors["type"] = "Tipo di event non valido"
	}
	if filter.From != nil && filter.To != nil && filter.From.After(*filter.To) {
		errors["dateRange"] = "La data 'from' non può essere successiva alla data 'to'"
	}
	if len(errors) > 0 {
		return []EventModel{}, NewValidationError(errors)
	}
	events, err := s.Repo.GetAll()
	if err != nil {
		return nil, NewInternalError("Errore nel recupero dell' event dal repository: " + err.Error())
	}
	filteredEvents := []EventModel{}
	for _, event := range events {
		if filter.ProjectId != "" && event.ProjectId != filter.ProjectId {
			continue
		}
		if filter.Type != "" && event.Type != filter.Type {
			continue
		}
		if filter.From != nil && event.CreatedAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && event.CreatedAt.After(*filter.To) {
			continue
		}
		filteredEvents = append(filteredEvents, event)
	}

	return filteredEvents, nil
}

func (s *Service) GetEvent(id string) (EventModel, error) {
	if id == "" {
		return EventModel{}, NewValidationError(map[string]string{"id": "id dell' event è obbligatorio"})
	}
	rm, err := s.Repo.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return EventModel{}, NewNotFoundError(id)
		}
		return EventModel{}, NewInternalError("errore interno nel recupero dell' event: " + err.Error())
	}

	return rm, nil

}

func (s *Service) ListTodayEvents(now time.Time) ([]EventModel, error) {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	fil := EventFilter{
		From:      &startOfDay,
		To:        &now,
		ProjectId: "",
		Type:      "",
	}
	return s.ListEvents(fil)

}

func NewValidationError(details map[string]string) DomainError {
	return DomainError{
		Code:    ERROR_VALIDATION,
		Message: "Errore di validazione",
		Details: details,
	}
}
func NewNotFoundError(id string) DomainError {
	return DomainError{
		Code:    ERROR_NOT_FOUND,
		Message: "event con id " + id + " non trovato",
		Details: map[string]string{"id": id},
	}
}

func NewInternalError(message string) DomainError {
	return DomainError{
		Code:    ERROR_INTERNAL,
		Message: message,
		Details: nil,
	}
}

func CreateEventsService(e *EventRepository) Service {
	return Service{
		Repo: e,
	}
}

func (d DomainError) Error() string {
	return d.Code + " Messaggio: " + d.Message
}
