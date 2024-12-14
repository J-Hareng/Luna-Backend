package bodymodels

import "server/src/api/db/models"

type LunaResponse struct {
	Status string      `json:"status"`
	Body   interface{} `json:"body"`
	Err    string      `json:"error"`
}

func NewLunaResponse_OK(body interface{}, status string) LunaResponse {
	if status == "" {
		status = "OK"
	}
	return LunaResponse{
		Status: status,
		Body:   body,
		Err:    "",
	}
}

func NewLunaResponse_ERROR(err string) LunaResponse {
	return LunaResponse{
		Status: "ERROR",
		Body:   nil,
		Err:    err,
	}
}
func NewLunaResponse_ERROR_INVALID_PAYLOAD() LunaResponse {
	return LunaResponse{
		Status: "ERROR",
		Body:   nil,
		Err:    "Invalid Payload",
	}
}

// GetAllUserDataMod is a struct that represents the body of a request to get all user data
type GetAllUserDataMod struct {
	User        models.User         `json:"user" binding:"required"`
	Collections []models.Collection `json:"collections" binding:"required"`
	Teams       []models.Team       `json:"teams" binding:"required"`
	UserTasks   []models.Task       `json:"tasks" binding:"required"`
	OtherUsers  []models.OtherUsers `json:"otherusers" binding:"required"`
}
type UserTaskToTeamsMod struct {
	Team  models.TeamLink `json:"team" binding:"required"`
	Tasks models.Task     `json:"tasks" binding:"required"`
}

type UpdatedTasksMod struct {
	TaskLs []models.TaskLink `json:"tasks" binding:"required"`
}
