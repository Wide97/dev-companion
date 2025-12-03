package projects

import "time"

type ProjectModel struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Type         string   `json:"type"`
	BuildCommand string   `json:"build_command"`
	TestCommand  string   `json:"test_command"`
	Tags         []string `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}