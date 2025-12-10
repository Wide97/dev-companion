package events

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type EventRepository struct {
	FilePath string
	Events   []EventModel
	Mu       sync.Mutex
}

func NewEventsRepository(dataDir string) (*EventRepository, error) {

	path := dataDir + "/events.json"

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

	return &EventRepository{
		FilePath: path,
		Events:   load,
		Mu:       sync.Mutex{},
	}, nil
}

func loadFromDisk(path string) ([]EventModel, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	} else if len(file) == 0 {
		return []EventModel{}, nil
	}

	arr := []EventModel{}
	err = json.Unmarshal(file, &arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (r *EventRepository) saveToDisk() error {
	data, err := json.MarshalIndent(r.Events, "", "  ")
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

func (r *EventRepository) GetAll() ([]EventModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	var sl = make([]EventModel, len(r.Events))
	copy(sl, r.Events)

	return sl, nil
}

func (r *EventRepository) GetById(id string) (EventModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	for _, event := range r.Events {
		if event.Id == id {
			return event, nil
		}
	}

	return EventModel{}, fmt.Errorf("event with id %s not found", id)

}

func (r *EventRepository) Create(event EventModel) (EventModel, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	tmp := time.Now().Format("20060102150405")
	id := "event_" + tmp + "dev_back"
	event.Id = id
	event.CreatedAt = time.Now()
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}

	r.Events = append(r.Events, event)
	err := r.saveToDisk()
	if err != nil {
		return EventModel{}, err
	}

	return event, nil

}
