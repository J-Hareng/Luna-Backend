package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostLink struct {
	NAME string `name:"title,omitempty" bson:"title,omitempty"`
	URL  string `json:"url,omitempty" bson:"url,omitempty"`
}

type Post struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	NAME        string             `json:"name" bson:"name"`
	DESCRIPTION string             `json:"des" bson:"des"`
	OUTHOR      string             `json:"outhor" bosn:"outhor"`
	GROUPID     string             `json:"groupid,omitempty" bson:"groupid,omitempty"`
}

func CreatePost(name string, des string, outhor string, gruID string) Post {
	id := primitive.NewObjectID()

	return Post{
		ID:          id,
		NAME:        name,
		DESCRIPTION: des,
		OUTHOR:      outhor,
		GROUPID:     gruID,
	}
}
