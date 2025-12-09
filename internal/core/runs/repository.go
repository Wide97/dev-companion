package runs

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type RunRepository struct {
	FilePath string
	Runs     []RunModel
	Mu       sync.Mutex
}

func loadFromDisk(path string) ([]RunModel, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	} else if len(file) == 0 {
		return []RunModel{}, nil
	}

	arr := []RunModel{}
	err = json.Unmarshal(file, &arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (r *RunRepository) saveToDisk() error {
	data, err := json.MarshalIndent(r.Runs, "", "  ")
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

func (r *RunRepository) GetAll() ([]RunModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	var sl = make([]RunModel, len(r.Runs))
	copy(sl, r.Runs)

	return sl, nil
}

func (r *RunRepository) GetById(id string) (RunModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	for _, project := range r.Runs {
		if project.Id == id {
			return project, nil
		}
	}

	return RunModel{}, fmt.Errorf("project with id %s not found", id)

}

func (r *RunRepository) Create(project RunModel) (RunModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	tmp := time.Now().Format("20060102150405")
	id := "pj_" + tmp + "dev_back"
	project.Id = id

	project.StartTime = time.Now()
	project.EndTime = time.Now()

	r.Runs = append(r.Runs, project)
	err := r.saveToDisk()
	if err != nil {
		return RunModel{}, err
	}

	return project, nil

}

func (r *RunRepository) Update(project RunModel) (RunModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if project.Id == "" {
		return RunModel{}, fmt.Errorf("id project necessario per l'aggiornamento")
	}
	foundIndex := -1
	var exisisting RunModel
	for i, p := range r.Runs {

		if p.Id == project.Id {
			foundIndex = i
			exisisting = p
			break
		}
	}

	if foundIndex == -1 {
		return RunModel{}, fmt.Errorf("project con id %s non trovato", project.Id)
	}

	exisisting.Command = project.Command
	exisisting.Status = project.Status
	exisisting.EndTime = project.EndTime
	exisisting.StartTime = project.StartTime
	exisisting.LogPath = project.LogPath
	exisisting.Type = project.Type
	exisisting.ProjectId = project.ProjectId
	r.Runs[foundIndex] = exisisting
	err := r.saveToDisk()
	if err != nil {
		return RunModel{}, err
	}

	return exisisting, nil

}

func (r *RunRepository) Delete(id string) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if id == "" {
		return fmt.Errorf("id project necessario per la cancellazione")
	}
	foundIndex := -1
	for i, p := range r.Runs {
		if p.Id == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("project con id %s non trovato", id)
	}
	r.Runs = append(r.Runs[:foundIndex], r.Runs[foundIndex+1:]...)

	err := r.saveToDisk()
	if err != nil {
		return err
	}

	return nil

}
