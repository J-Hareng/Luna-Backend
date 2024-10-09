package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	NAME        string             `json:"name" bson:"name"`
	DESCRIPTION string             `json:"des" bson:"des"`
	USERS       []UserLink         `json:"users,omitempty" bson:"users,omitempty"`
	TASKS       []TaskLink         `json:"tasks,omitempty" bson:"tasks,omitempty"`
	DONETASKS   []TaskLink         `json:"donetasks,omitempty" bson:"donetasks,omitempty"`
	POSTS       []PostLink         `json:"posts,omitempty" bson:"posts,omitempty"`
	GROUPID     string             `json:"groupid,omitempty" bson:"groupid,omitempty"`
}
type TeamLink struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
}

func CreateTeam(name string, user UserLink, Des string, gruID string) Team {
	userArr := make([]UserLink, 0)
	userArr = append(userArr, user)
	return Team{
		ID:          primitive.NewObjectID(),
		NAME:        name,
		DESCRIPTION: Des,
		USERS:       userArr,
		TASKS:       []TaskLink{},
		DONETASKS:   []TaskLink{},
		POSTS:       []PostLink{},
		GROUPID:     gruID,
	}
}
