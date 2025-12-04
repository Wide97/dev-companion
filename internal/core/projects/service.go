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

func ListProjects() {}

func GetProject(id string) {}

func CreateProject(input CreateProjectInput) {}

func UpdateProject(id string, input UpdateProjectInput) {}

func DeleteProject(id string) {}

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
