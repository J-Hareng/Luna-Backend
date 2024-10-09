package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Collection struct {
	ID     *primitive.ObjectID `json:"_id" bson:"_id"`
	Name   string              `json:"name" bson:"name"`
	Files  []File              `json:"files" bson:"files"`
	GrupID string              `json:"groupid" bson:"groupid"`
	Tags   []string            `json:"tags" bson:"tags"`
}

func CreateCollection(name string, grupID string) Collection {
	id := primitive.NewObjectID()
	files := make([]File, 0)
	return Collection{
		ID:     &id,
		Name:   name,
		Files:  files,
		GrupID: grupID,
		Tags:   []string{},
	}
}
