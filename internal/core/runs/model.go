package runs

import "time"

const RunTypeBuild = "BUILD"
const RunTypeTest = "TEST"
const RunTypeCustom = "CUSTOM"

const RunStatusPending = "PENDING"
const RunStatusRunning = "RUNNING"
const RunStatusSuccess = "SUCCESS"
const RunStatusFailed = "FAILED"

type RunModel struct {
	Id        string     `json:"id"`
	ProjectId string     `json:"projectId"`
	Type      string     `json:"type"`
	Command   string     `json:"command"`
	Status    string     `json:"status"`
	StartTime time.Time  `json:"startTime"`
	EndTime   *time.Time `json: "endTime,omitempty"`
	LogPath   string     `json:"logPath"`
}
