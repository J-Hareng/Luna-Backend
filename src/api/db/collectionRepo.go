package db

import (
	"context"
	"server/src/api/db/models"
	"server/src/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) AddCollection(name string, GruID string) (*mongo.InsertOneResult, error) {
	coll := models.CreateCollection(name, GruID)
	result, err := db.Coll.InsertOne(context.TODO(), coll)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *DB) AddFileToCollection(file models.File, collectionId primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: collectionId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "files", Value: file}}}}
	result, err := db.Coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (db *DB) RemoveFileFromCollection(file models.File, collectionId primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: collectionId}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "files", Value: file}}}}
	result, err := db.Coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *DB) RemoveCollection(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	GetFullVal, err := db.GetCollection(id)
	if err != nil {
		return nil, err
	}
	for _, file := range GetFullVal.Files {
		if err := helper.RemoveFileRaw(file.Path); err != nil {
			return nil, err
		}
	}
	result, err := db.Coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (db *DB) GetCollection(id primitive.ObjectID) (models.Collection, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	cursor := db.Coll.FindOne(context.Background(), filter)
	var collection models.Collection
	if err := cursor.Decode(&collection); err != nil {
		return models.Collection{}, err
	}
	return collection, nil
}

func (db *DB) GetCollectionsViaGroupID(gruID string) ([]models.Collection, error) {
	filter := bson.D{{Key: "groupid", Value: gruID}}
	cursor, err := db.Coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var collections []models.Collection
	if err = cursor.All(context.Background(), &collections); err != nil {
		return nil, err
	}
	return collections, nil
}

func (db *DB) AddTag(id primitive.ObjectID, tag string) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "tags", Value: tag}}}}
	result, err := db.Coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
