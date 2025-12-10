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

func NewRunRepository(dataDir string) (*RunRepository, error) {

	path := dataDir + "/runs.json"

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

	return &RunRepository{
		FilePath: path,
		Runs:     load,
		Mu:       sync.Mutex{},
	}, nil
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
	for _, run := range r.Runs {
		if run.Id == id {
			return run, nil
		}
	}

	return RunModel{}, fmt.Errorf("run with id %s not found", id)

}

func (r *RunRepository) Create(run RunModel) (RunModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	tmp := time.Now().Format("20060102150405")
	id := "run_" + tmp + "dev_back"
	run.Id = id
	run.Status = RunStatusPending
	run.StartTime = time.Now()
	run.EndTime = nil

	r.Runs = append(r.Runs, run)
	err := r.saveToDisk()
	if err != nil {
		return RunModel{}, err
	}

	return run, nil

}

func (r *RunRepository) Update(run RunModel) (RunModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if run.Id == "" {
		return RunModel{}, fmt.Errorf("id run necessario per l'aggiornamento")
	}
	foundIndex := -1
	var exisisting RunModel
	for i, p := range r.Runs {

		if p.Id == run.Id {
			foundIndex = i
			exisisting = p
			break
		}
	}

	if foundIndex == -1 {
		return RunModel{}, fmt.Errorf("run con id %s non trovato", run.Id)
	}

	exisisting.Status = run.Status
	exisisting.EndTime = run.EndTime
	exisisting.LogPath = run.LogPath
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
		return fmt.Errorf("id run necessario per la cancellazione")
	}
	foundIndex := -1
	for i, p := range r.Runs {
		if p.Id == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("run con id %s non trovato", id)
	}
	r.Runs = append(r.Runs[:foundIndex], r.Runs[foundIndex+1:]...)

	err := r.saveToDisk()
	if err != nil {
		return err
	}

	return nil

}
