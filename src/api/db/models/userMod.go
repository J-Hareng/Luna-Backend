package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	NAME        string             `json:"name" bson:"name"`
	EMAIL       string             `json:"email" bson:"email"`
	PASSWORD    string             `json:"password" bson:"password"`
	SALT        string             `json:"salt" bson:"salt"`
	AVALABILITY string             `json:"avalability" bson:"avalability"`
	POSTS       []PostLink         `json:"posts,omitempty" bson:"posts,omitempty"`
	TASKS       []TaskLink         `json:"tasks,omitempty" bson:"tasks,omitempty"`
	GROUPID     string             `json:"groupid,omitempty" bson:"groupid,omitempty"`
}

type PublicUser struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	NAME        string             `json:"name" bson:"name"`
	EMAIL       string             `json:"email" bson:"email"`
	AVALABILITY string             `json:"avalability" bson:"avalability"`
	POSTS       []PostLink         `json:"posts,omitempty" bson:"posts,omitempty"`
	TASKS       []TaskLink         `json:"tasks,omitempty" bson:"tasks,omitempty"`
	GROUPID     string             `json:"groupid,omitempty" bson:"groupid,omitempty"`
}
type UserLink struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	NAME string             `json:"name" bson:"name"`
}

func CreatePublicUser(name string, email string, gruID string) PublicUser {
	id := primitive.NewObjectID()

	return PublicUser{
		ID:          id,
		NAME:        name,
		EMAIL:       email,
		GROUPID:     gruID,
		AVALABILITY: "ACTIVE",
		POSTS:       []PostLink{},
		TASKS:       []TaskLink{},
	}
}
func CreateUser(name string, email string, passwd string, gruID string, salt string) User {
	id := primitive.NewObjectID()

	return User{
		ID:          id,
		NAME:        name,
		EMAIL:       email,
		PASSWORD:    passwd,
		GROUPID:     gruID,
		SALT:        salt,
		AVALABILITY: "ACTIVE",
		POSTS:       []PostLink{},
		TASKS:       []TaskLink{},
	}
}
func TransformUserToPublicUser(user User) PublicUser {
	return PublicUser{
		ID:          user.ID,
		NAME:        user.NAME,
		EMAIL:       user.EMAIL,
		AVALABILITY: user.AVALABILITY,
		POSTS:       user.POSTS,
		TASKS:       user.TASKS,
		GROUPID:     user.GROUPID,
	}
}
func TransformUserToUserLink(user User) UserLink {
	return UserLink{
		ID:   user.ID,
		NAME: user.NAME,
	}
}
