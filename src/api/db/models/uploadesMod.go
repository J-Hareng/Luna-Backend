package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Uploads struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	NAME        string             `json:"name" bson:"name"`
	DESCRIPTION string             `json:"des" bson:"des"`
	CREATOR     UserLink           `json:"creator" bosn:"creator"`
	LOCATION    string             `json:"location" bosn:"location"`

	COLLECTION primitive.ObjectID `json:"coll" bson:"coll"`
	GROUPID    string             `json:"groupid,omitempty" bson:"groupid,omitempty"`
}
