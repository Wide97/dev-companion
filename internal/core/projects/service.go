package projects

type DomainError struct {
	Code    string
	Message string
	Details map[string]string
}

type Service struct {
	Repo *ProjectRepository
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

func (s *Service) CreateProject(input CreateProjectInput) {}

func (s *Service) UpdateProject(id string, input UpdateProjectInput) {}

func (s *Service) DeleteProject(id string) {}

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

func (r *ProjectRepository) CreatePjService() Service {
	return Service{
		Repo: r,
	}
}

func (d DomainError) Error() string {
	return d.Code + " Messaggio: " + d.Message
}
