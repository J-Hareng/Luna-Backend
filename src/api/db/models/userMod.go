package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	NAME        string             `json:"name" bson:"name"`
	EMAIL       string             `json:"email" bson:"email"`
	PASSWORD    string             `json:"password" bson:"password"`
	AVALABILITY string             `json:"avalability" bson:"avalability"`
	POSTS       []PostLink         `json:"posts,omitempty" bson:"posts,omitempty"`
	TASKS       []TaskLink         `json:"tasks,omitempty" bson:"tasks,omitempty"`
	GROUPID     string             `json:"groupid,omitempty" bson:"groupid,omitempty"`
}

type UserLink struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	NAME string             `json:"name" bson:"name"`
}

func CreateUser(name string, email string, passwd string, gruID string) User {
	id := primitive.NewObjectID()

	return User{
		ID:          id,
		NAME:        name,
		EMAIL:       email,
		PASSWORD:    passwd,
		GROUPID:     gruID,
		AVALABILITY: "ACTIVE",
		POSTS:       []PostLink{},
		TASKS:       []TaskLink{},
	}
}

type OtherUsers struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	NAME        string             `json:"name" bson:"name"`
	EMAIL       string             `json:"email" bson:"email"`
	AVALABILITY string             `json:"avalability" bson:"avalability"`
	POSTS       []PostLink         `json:"posts,omitempty" bson:"posts,omitempty"`
	TASKS       []TaskLink         `json:"tasks,omitempty" bson:"tasks,omitempty"`
	GROUPID     string             `json:"groupid,omitempty" bson:"groupid,omitempty"`
}
