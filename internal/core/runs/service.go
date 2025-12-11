package runs

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
	Repo *RunRepository
}

const ERROR_VALIDATION = "VALIDATION"
const ERROR_NOT_FOUND = "NOT_FOUND"
const ERROR_CONFLICT = "CONFLICT"
const ERROR_INTERNAL = "INTERNAL"

type RunFilter struct {
	ProjectId string
	Type      string
	Status    string
	From      *time.Time
	To        *time.Time
}

func (s *Service) ListRuns(filter RunFilter) ([]RunModel, error) {
	errors := map[string]string{}
	if filter.Type != "" && filter.Type != RunTypeBuild && filter.Type != RunTypeCustom && filter.Type != RunTypeTest {
		errors["type"] = "Tipo di run non valido"
	}
	if filter.Status != "" && filter.Status != RunStatusFailed && filter.Status != RunStatusPending && filter.Status != RunStatusRunning && filter.Status != RunStatusSuccess {
		errors["status"] = "Stato di run non valido"
	}
	if filter.From != nil && filter.To != nil && filter.From.After(*filter.To) {
		errors["dateRange"] = "La data 'from' non può essere successiva alla data 'to'"
	}
	if len(errors) > 0 {
		return []RunModel{}, NewValidationError(errors)
	}
	runs, err := s.Repo.GetAll()
	if err != nil {
		return nil, NewInternalError("Errore nel recupero delle run dal repository: " + err.Error())
	}
	filteredRuns := []RunModel{}
	for _, run := range runs {
		if filter.ProjectId != "" && run.ProjectId != filter.ProjectId {
			continue
		}
		if filter.Type != "" && run.Type != filter.Type {
			continue
		}
		if filter.Status != "" && run.Status != filter.Status {
			continue
		}
		if filter.From != nil && run.StartTime.Before(*filter.From) {
			continue
		}
		if filter.To != nil && run.StartTime.After(*filter.To) {
			continue
		}
		filteredRuns = append(filteredRuns, run)
	}

	return filteredRuns, nil
}

func (s *Service) GetRun(id string) (RunModel, error) {
	if id == "" {
		return RunModel{}, NewValidationError(map[string]string{"id": "id della run è obbligatorio"})
	}
	rm, err := s.Repo.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return RunModel{}, NewNotFoundError(id)
		}
		return RunModel{}, NewInternalError("errore interno nel recupero della run: " + err.Error())
	}

	return rm, nil

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
		Message: "run con id " + id + " non trovato",
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

func CreateRunSerivce(r *RunRepository) Service {
	return Service{
		Repo: r,
	}
}

func (d DomainError) Error() string {
	return d.Code + " Messaggio: " + d.Message
}
