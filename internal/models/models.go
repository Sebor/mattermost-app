package models

import "time"

type TaskCustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TaskCustomFields struct {
	CustomField []TaskCustomField `json:"customField"`
}

type TaskExternalID struct {
	SystemID string `json:"systemId"`
}

type TaskExternalIDs struct {
	ExternalID []TaskExternalID `json:"externalId"`
}

type Task struct {
	Project      string           `json:"project"`
	IssueType    string           `json:"issueType"`
	Summary      string           `json:"summary"`
	Created      string           `json:"created"`
	Status       string           `json:"status,omitempty"`
	Priority     string           `json:"priority,omitempty"`
	Description  string           `json:"description,omitempty"`
	Reporter     string           `json:"reporter,omitempty"`
	Assignee     string           `json:"assignee,omitempty"`
	ExternalIds  TaskExternalIDs  `json:"externalIds"`
	CustomFields TaskCustomFields `json:"customFields"`
}

type CreateTask struct {
	Task         Task   `json:"task"`
	ActionAuthor string `json:"actionAuthor,omitempty"`
	SourceSystem string `json:"sourceSystem,omitempty"`
}

type NewTask struct {
	CreateTask CreateTask `json:"CreateTask"`
}

func (nt NewTask) CreateWithDefaults() NewTask {
	// Set default values
	nt.CreateTask.Task.Project = "TEST"
	nt.CreateTask.Task.IssueType = "Issue"
	nt.CreateTask.Task.Status = "new"
	nt.CreateTask.Task.Priority = "medium"
	nt.CreateTask.Task.Created = time.Now().Format("2006-01-02T15:04:05")
	return nt
}

/* Response structs */
type TaskError struct {
	Actor   string    `json:"actor"`
	Time    time.Time `json:"time"`
	Message string    `json:"message,omitempty"`
}

type TaskResponse struct {
	Status     string    `json:"status"`
	IssueKey   string    `json:"issueKey,omitempty"`
	Attachment []any     `json:"attachment,omitempty"`
	Error      TaskError `json:"error,omitempty"`
}

type NewTaskResponse struct {
	CreateTaskResponse TaskResponse `json:"CreateTaskResponse"`
}
