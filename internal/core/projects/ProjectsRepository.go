package projects

type PjRepository interface {
	GetAll() ([]ProjectModel, error)
	GetById(id string) (ProjectModel, error)
	Create(p ProjectModel) (ProjectModel, error)
	Update(p ProjectModel) (ProjectModel, error)
	Delete(id string) error
}
