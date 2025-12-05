package projects

import "strings"

type DomainError struct {
	Code    string
	Message string
	Details map[string]string
}

type Service struct {
	Repo PjRepository
}

const ERROR_VALIDATION = "VALIDATION"
const ERROR_NOT_FOUND = "NOT_FOUND"
const ERROR_CONFLICT = "CONFLICT"
const ERROR_INTERNAL = "INTERNAL"

type CreateProjectInput struct {
	Name         string
	Path         string
	Type         string
	BuildCommand string
	TestCommand  string
	Tags         []string
}

type UpdateProjectInput struct {
	Name         string
	Path         string
	Type         string
	BuildCommand string
	TestCommand  string
	Tags         []string
}

func (s *Service) ListProjects() ([]ProjectModel, error) {
	projects, err := s.Repo.GetAll()
	if err != nil {
		return nil, NewInternalError("Errore nel recupero dei progetti dal repository: " + err.Error())
	}
	return projects, nil
}

func (s *Service) GetProject(id string) (ProjectModel, error) {
	if id == "" {
		return ProjectModel{}, NewValidationError(map[string]string{"id": "id del progetto è obbligatorio"})
	}

	result, err := s.Repo.GetById(id)
	if err != nil {
		return ProjectModel{}, NewInternalError("Errore nella ricerca del project con id " + id + ": " + err.Error())

	}

	return result, nil

}

func (s *Service) CreateProject(input CreateProjectInput) (ProjectModel, error) {
	normalized := input
	normalized.Name = strings.TrimSpace(normalized.Name)
	normalized.Path = strings.TrimSpace(normalized.Path)
	normalized.Type = strings.TrimSpace(normalized.Type)
	normalized.Type = strings.ToLower(normalized.Type)
	normalized.BuildCommand = strings.TrimSpace(normalized.BuildCommand)
	normalized.TestCommand = strings.TrimSpace(normalized.TestCommand)
	if normalized.Tags == nil {
		normalized.Tags = []string{}
	}

	var cleanTags []string
	var seen map[string]bool = make(map[string]bool)
	for _, val := range normalized.Tags {
		trimmed := strings.TrimSpace(val)
		minus := strings.ToLower(trimmed)
		if minus == "" {
			continue
		}
		if seen[minus] {
			continue
		}

		cleanTags = append(cleanTags, minus)
		seen[minus] = true
	}

	normalized.Tags = cleanTags

	errors := make(map[string]string)
	if normalized.Name == "" {
		errors["name"] = "name è obbligatorio"
	}
	if normalized.Path == "" {
		errors["path"] = "path è obbligatorio"
	}

	if normalized.Type == "" || (normalized.Type != "maven" && normalized.Type != "gradle" && normalized.Type != "node" && normalized.Type != "docker" && normalized.Type != "other") {
		errors["type"] = "type deve essere uno tra: maven, gradle, node, docker, other"
	}
	if len(errors) > 0 {
		return ProjectModel{}, NewValidationError(errors)
	}

	all, err := s.Repo.GetAll()
	if err != nil {
		return ProjectModel{}, NewInternalError("Errore nel recupero dei progetti dal repository: " + err.Error())
	}
	for _, p := range all {
		if p.Path == normalized.Path {
			return ProjectModel{}, NewConflictError("path", normalized.Path)
		}
	}

	nor := ProjectModel{
		Name:         normalized.Name,
		Path:         normalized.Path,
		Type:         normalized.Type,
		BuildCommand: normalized.BuildCommand,
		TestCommand:  normalized.TestCommand,
		Tags:         normalized.Tags,
	}

	res, err := s.Repo.Create(nor)
	if err != nil {
		return ProjectModel{}, NewInternalError("Errore nella creazione del project: " + err.Error())
	}

	return res, nil

}

func (s *Service) UpdateProject(id string, input UpdateProjectInput) (ProjectModel, error) {
	if id == "" {
		return ProjectModel{}, NewValidationError(map[string]string{"id": "id del progetto è obbligatorio"})
	}

	normalized := input
	normalized.Name = strings.TrimSpace(normalized.Name)
	normalized.Path = strings.TrimSpace(normalized.Path)
	normalized.Type = strings.TrimSpace(normalized.Type)
	normalized.Type = strings.ToLower(normalized.Type)
	normalized.BuildCommand = strings.TrimSpace(normalized.BuildCommand)
	normalized.TestCommand = strings.TrimSpace(normalized.TestCommand)
	if normalized.Tags == nil {
		normalized.Tags = []string{}
	}
	var cleanTags []string
	var seen map[string]bool = make(map[string]bool)
	for _, val := range normalized.Tags {
		trimmed := strings.TrimSpace(val)
		minus := strings.ToLower(trimmed)
		if minus == "" {
			continue
		}
		if seen[minus] {
			continue
		}

		cleanTags = append(cleanTags, minus)
		seen[minus] = true
	}

	normalized.Tags = cleanTags

	errors := make(map[string]string)
	if normalized.Name == "" {
		errors["name"] = "name è obbligatorio"
	}
	if normalized.Path == "" {
		errors["path"] = "path è obbligatorio"
	}

	if normalized.Type == "" || (normalized.Type != "maven" && normalized.Type != "gradle" && normalized.Type != "node" && normalized.Type != "docker" && normalized.Type != "other") {
		errors["type"] = "type deve essere uno tra: maven, gradle, node, docker, other"
	}
	if len(errors) > 0 {
		return ProjectModel{}, NewValidationError(errors)
	}

	all, err := s.Repo.GetAll()
	if err != nil {
		return ProjectModel{}, NewInternalError("Errore nel recupero dei progetti dal repository: " + err.Error())
	}
	for _, p := range all {
		if p.Id == id {
			continue
		} else if p.Path == normalized.Path {
			return ProjectModel{}, NewConflictError("path", normalized.Path)
		}
	}

	exisisting, err := s.Repo.GetById(id)
	if err != nil {
		return ProjectModel{}, NewInternalError("Errore nel recupero del progetto con id: " + id)
	}

	updated := exisisting
	updated.Name = normalized.Name
	updated.Path = normalized.Path
	updated.Type = normalized.Type
	updated.BuildCommand = normalized.BuildCommand
	updated.TestCommand = normalized.TestCommand
	updated.Tags = normalized.Tags

	up, err := s.Repo.Update(updated)
	if err != nil {
		return ProjectModel{}, NewInternalError("Errore durante l' aggiornamento del project: " + err.Error())
	}

	return up, nil

}

func (s *Service) DeleteProject(id string) error {
	if id == "" {
		return NewValidationError(map[string]string{"id": "id del progetto è obbligatorio"})
	}

	err := s.Repo.Delete(id)
	if err != nil {
		return NewInternalError("Errore durante la cancellazione: " + err.Error())
	}

	return nil
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
		Message: "Project con id " + id + " non trovato",
		Details: map[string]string{"id": id},
	}
}
func NewConflictError(field string, value string) DomainError {
	return DomainError{
		Code:    ERROR_CONFLICT,
		Message: "Conflitto: " + field + " con value " + value + " esiste già",
		Details: map[string]string{field: value},
	}
}
func NewInternalError(message string) DomainError {
	return DomainError{
		Code:    ERROR_INTERNAL,
		Message: message,
		Details: nil,
	}
}

func CreatePjService(p PjRepository) Service {
	return Service{
		Repo: p,
	}
}

func (d DomainError) Error() string {
	return d.Code + " Messaggio: " + d.Message
}
