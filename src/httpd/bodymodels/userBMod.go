package bodymodels

import (
	"server/src/api/db/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// * User
type AddUserMod struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"hashPassword" binding:"required"`
	Salt     string `json:"salt" binding:"required"`
	Key      string `json:"key" binding:"required"`
}

// * Token
type RequestSessionTokenMod struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RequestEmailKeyMod struct {
	Email string `json:"email" binding:"required"`
}

type AddTeamMod struct {
	Name    string          `json:"name" binding:"required"`
	Des     string          `json:"des" binding:"required"`
	User    models.UserLink `json:"user" binding:"required"`
	Groupid string          `json:"groupid" binding:"required"`
}

type AddUserToTeamMod struct {
	Team models.TeamLink `json:"team" binding:"required"`
	User models.UserLink `json:"user" binding:"required"`
}
type GetUserTeamsMod struct {
	User models.UserLink `json:"user" binding:"required"`
}
type AddTaskMod struct {
	Team        models.TeamLink `json:"team" binding:"required"`
	Description string          `json:"des"`
	Name        string          `json:"name" binding:"required"`
	Groupid     string          `json:"groupid" binding:"required"`
	CollID      string          `json:"collID"`
	Tags        []string        `json:"tags" binding:"required"`
}
type RemoveTaskMod struct {
	Tasks  models.TaskLink    `json:"tasks"`
	TeamId primitive.ObjectID `json:"teamId"`
}
type RemoveTeamMod struct {
	Team models.TeamLink `json:"team"`
}
type GetAllTeamsMod struct {
	Groupid string `json:"groupid" binding:"required"`
}

type EditAssingmentForTaskMod struct {
	Task models.TaskLink `json:"task" binding:"required"`
	User models.UserLink `json:"user" binding:"required"`
}
type GetAllTasksFromTeamMod struct {
	Tasks []models.TaskLink `json:"tasks"`
}
type AddCollectionMod struct {
	Name    string `json:"name" binding:"required"`
	Groupid string `json:"groupid" binding:"required"`
}

type GetAllCollectionsMod struct {
	Groupid string `json:"groupid" binding:"required"`
}

type AddTagMod struct {
	Name   string             `json:"tag" binding:"required"`
	CollId primitive.ObjectID `json:"collID" binding:"required"`
}
type EditTagMod struct {
	Name    string             `json:"newtag" binding:"required"`
	OldName string             `json:"oldtag" binding:"required"`
	CollId  primitive.ObjectID `json:"collID" binding:"required"`
}

type RemoveCollectionMod struct {
	CollId primitive.ObjectID `json:"collID" binding:"required"`
}
type DonwloadFileMod struct {
	File models.File `json:"file" binding:"required"`
}
type RemoveFileFromCollectionMod struct {
	File   models.File        `json:"file" binding:"required"`
	CollId primitive.ObjectID `json:"collID" binding:"required"`
}

// * Team
type EditTeamMod struct {
	Team models.Team `json:"team" binding:"required"`
}

// * Task
type UpdateTaskRequestMod struct {
	Team models.TeamLink `json:"team" binding:"required"`
}

// * Group
type AddGrupIDTokenMod struct {
	User models.UserLink `json:"user" binding:"required"`
}
type GetGrupIDTokenMod struct {
	Key string `json:"key" binding:"required"`
}
type AddGrupIDMod struct {
	GroupidName string `json:"groupidname" binding:"required"`
}

// * File
type GetMultibleFilesMod struct {
	Files []models.File `json:"files" binding:"required"`
}
