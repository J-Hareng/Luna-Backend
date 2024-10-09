package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskLink struct {
	NAME string             `json:"name,omitempty" bson:"name,omitempty"`
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}
type Task struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	NAME        string             `json:"name" bson:"name"`
	DESCRIPTION string             `json:"des" bson:"des"`
	INPROGRESS  []UserLink         `json:"inprogress,omitempty" bosn:"inprogress,omitempty"`
	TEAM        TeamLink           `json:"team,omitempty" bosn:"team,omitempty"`

	COLLECTION string   `json:"coll" bson:"coll"`
	TAGS       []string `json:"tags" bson:"tags"`
	GROUPID    string   `json:"groupid,omitempty" bson:"groupid,omitempty"`
}
type TaskWithStringId struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	NAME        string             `json:"name" bson:"name"`
	DESCRIPTION string             `json:"des" bson:"des"`
	INPROGRESS  []UserLink         `json:"inprogress,omitempty" bosn:"inprogress,omitempty"`
	TEAM        TeamLink           `json:"team,omitempty" bosn:"team,omitempty"`

	COLLECTION string   `json:"coll" bson:"coll"`
	TAGS       []string `json:"tags" bson:"tags"`
	GROUPID    string   `json:"groupid,omitempty" bson:"groupid,omitempty"`
}

func CreateTask(name string, des string, team TeamLink, gruID string, collID string, tags []string) Task {
	id := primitive.NewObjectID()
	inprog := []UserLink{}
	return Task{
		ID:          id,
		NAME:        name,
		DESCRIPTION: des,
		TEAM:        team,
		INPROGRESS:  inprog,
		GROUPID:     gruID,

		COLLECTION: collID,
		TAGS:       tags,
	}
}
