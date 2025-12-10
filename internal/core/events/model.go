package events

import "time"

type EventModel struct {
	Id        string            `json:"id"`
	ProjectId string            `json:"projectId"`
	Type      string            `json:"type"`
	Message   string            `json:"message"`
	CreatedAt time.Time         `json:"createdAt"`
	Metadata  map[string]string `json:"metadata"`
}

const EventPjBuildStarted = "PROJECT_BUILD_STARTED"
const EventPjBuildSuccess = "PROJECT_BUILD_SUCCESS"
const EventPjBuildFailed = "PROJECT_BUILD_FAILED"
const EventRunStarted = "RUN_STARTED"
const EventRunSuccess = "RUN_SUCCESS"
const EventRunFailed = "RUN_FAILED"
