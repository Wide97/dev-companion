package projects

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type ProjectRepository struct {
	FilePath     string
	ProjectModel []ProjectModel
	Mu           sync.Mutex
}

func NewProjectRepository(dataDir string) (*ProjectRepository, error) {

	path := dataDir + "/projects.json"

	file, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Nessun dato trovato")
		err1 := os.WriteFile(path, []byte("[]"), 0644)
		if err1 != nil {
			return nil, err1
		}
	} else if err != nil {
		return nil, err
	}

	fmt.Println(file)

	load, err2 := loadFromDisk(path)
	if err2 != nil {
		return nil, err2
	}

	fmt.Println("Dati caricati:", load)

	return &ProjectRepository{
		FilePath:     path,
		ProjectModel: load,
		Mu:           sync.Mutex{},
	}, nil
}

func loadFromDisk(path string) ([]ProjectModel, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	} else if len(file) == 0 {
		return []ProjectModel{}, nil
	}

	arr := []ProjectModel{}
	err = json.Unmarshal(file, &arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (r *ProjectRepository) saveToDisk() error {
	data, err := json.MarshalIndent(r.ProjectModel, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := r.FilePath + ".tmp"
	err = os.WriteFile(tmpPath, data, 0644)
	if err != nil {
		return err
	}
	err = os.Rename(tmpPath, r.FilePath)
	if err != nil {
		temp := os.Remove(tmpPath)
		fmt.Println("Rimosso: ", temp)
		return err
	}

	return nil
}

func (r *ProjectRepository) GetAll() ([]ProjectModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	var sl = make([]ProjectModel, len(r.ProjectModel))
	copy(sl, r.ProjectModel)

	return sl, nil
}

func (r *ProjectRepository) GetById(id string) (ProjectModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	for _, project := range r.ProjectModel {
		if project.Id == id {
			return project, nil
		}
	}

	return ProjectModel{}, fmt.Errorf("project with id %s not found", id)

}

func (r *ProjectRepository) Create(project ProjectModel) (ProjectModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	tmp := time.Now().Format("20060102150405")
	id := "pj_" + tmp + "dev_back"
	project.Id = id

	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	r.ProjectModel = append(r.ProjectModel, project)
	err := r.saveToDisk()
	if err != nil {
		return ProjectModel{}, err
	}

	return project, nil

}

func (r *ProjectRepository) Update(project ProjectModel) (ProjectModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if project.Id == "" {
		return ProjectModel{}, fmt.Errorf("id project necessario per l'aggiornamento")
	}
	foundIndex := -1
	var exisisting ProjectModel
	for i, p := range r.ProjectModel {

		if p.Id == project.Id {
			foundIndex = i
			exisisting = p
			break
		}
	}

	if foundIndex == -1 {
		return ProjectModel{}, fmt.Errorf("project con id %s non trovato", project.Id)
	}

	exisisting.Name = project.Name
	exisisting.Path = project.Path
	exisisting.Type = project.Type
	exisisting.BuildCommand = project.BuildCommand
	exisisting.TestCommand = project.TestCommand
	exisisting.Tags = project.Tags
	exisisting.UpdatedAt = time.Now()

	r.ProjectModel[foundIndex] = exisisting

	err := r.saveToDisk()
	if err != nil {
		return ProjectModel{}, err
	}

	return exisisting, nil

}

func (r *ProjectRepository) Delete(id string) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if id == "" {
		return fmt.Errorf("id project necessario per la cancellazione")
	}
	foundIndex := -1
	for i, p := range r.ProjectModel {
		if p.Id == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("project con id %s non trovato", id)
	}
	r.ProjectModel = append(r.ProjectModel[:foundIndex], r.ProjectModel[foundIndex+1:]...)

	err := r.saveToDisk()
	if err != nil {
		return err
	}

	return nil

}
