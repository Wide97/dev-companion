package projects

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
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
