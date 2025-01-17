package db

import (
	"context"
	"errors"
	"fmt"
	"server/src/api/db/models"
	"server/src/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) AddTeam(Name string, User models.UserLink, Des string, gruID string) (*mongo.InsertOneResult, error) {
	newTeam := models.CreateTeam(Name, User, Des, gruID)
	result, err := db.Team.InsertOne(context.TODO(), newTeam)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (db *DB) GetTeam(TeamID primitive.ObjectID) ([]models.Team, error) {
	filter := bson.D{{Key: "_id", Value: TeamID}} // Empty filter to select all documents
	cursor, err := db.Team.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var team []models.Team
	if err := cursor.All(context.TODO(), &team); err != nil {
		return nil, err
	}
	return team, nil
}
func (db *DB) RemoveTeam(teaml models.TeamLink) ([]*mongo.DeleteResult, error) {
	filter := bson.D{{Key: "_id", Value: teaml.ID}}
	team, err := db.GetTeam(teaml.ID)
	if err != nil || len(team) == 0 {
		return nil, err
	}
	result, err := db.Team.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	val, err := db.RemoveTasks(team[0].TASKS)
	if err != nil {
		return nil, err
	}

	//! NEED TO DO THE SAME THING WITH POSTS
	resultArray := []*mongo.DeleteResult{result, val}
	return resultArray, nil
}
func (db *DB) RemoveUserFromTeam(team models.TeamLink, user models.UserLink) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: team.ID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "users", Value: user}}}}
	result, err := db.Team.UpdateOne(context.TODO(), filter, update)
	fmt.Println(result)
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount <= 0 {
		return nil, errors.New("no documents with these properties found")
	}
	return result, nil
}
func (db *DB) GetAllTeams(grupID string) ([]models.Team, error) {
	filter := bson.D{{Key: "groupid", Value: grupID}} // Empty filter to select all documents
	cursor, err := db.Team.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var teams []models.Team
	if err := cursor.All(context.TODO(), &teams); err != nil {
		return nil, err
	}
	return teams, nil
}
func (db *DB) RemoveUserFromAllTeamsInGroup(user models.UserLink, grupID string) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "groupid", Value: grupID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "users", Value: user}}}}
	result, err := db.Team.UpdateMany(context.TODO(), filter, update)
	fmt.Println(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (db *DB) AddUserToTeam(team models.TeamLink, user models.UserLink) (*mongo.UpdateResult, error) {
	fmt.Println(user)
	filter := bson.D{{Key: "_id", Value: team.ID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "users", Value: user}}}}
	result, err := db.Team.UpdateOne(context.TODO(), filter, update)
	fmt.Println(result)
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount <= 0 {
		return nil, errors.New("no documents with these properties found")
	}
	return result, nil
}
func (db *DB) EditTeam(team models.Team) (*mongo.UpdateResult, error) {
	fmt.Println(team.ID)
	fmt.Println(team.NAME)

	filter := bson.D{{Key: "_id", Value: team.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: team.NAME}, {Key: "des", Value: team.DESCRIPTION}}}}
	result, err := db.Team.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount <= 0 {
		return nil, errors.New("no documents with these properties found")
	}
	return result, nil
}

func (db *DB) InsertInInTeamArray(id primitive.ObjectID, key string, val any, property helper.InsertOption[string]) (*mongo.UpdateResult, error) {
	prop := "$push"
	if property.Property != "" {
		prop = property.Property
	}
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: prop, Value: bson.D{{Key: key, Value: val}}}}
	result, err := db.Team.UpdateOne(context.TODO(), filter, update)
	fmt.Println(result)
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount <= 0 {
		return nil, errors.New("no documents with these properties found")
	}
	return result, nil
}
func (db *DB) RemoveFromTeamArray(id primitive.ObjectID, key string, val any, property helper.InsertOption[string]) (*mongo.UpdateResult, error) {
	prop := "$pull"
	if property.Property != "" {
		prop = property.Property
	}
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: prop, Value: bson.D{{Key: key, Value: val}}}}
	result, err := db.Team.UpdateOne(context.TODO(), filter, update)
	fmt.Println(result)
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount <= 0 {
		return nil, errors.New("no documents with these properties found")
	}
	return result, nil
}

func (db *DB) RemoveFromMultibleTeamArray(key string, val any, property helper.InsertOption[string]) (*mongo.UpdateResult, error) {
	prop := "$pull"
	if property.Property != "" {
		prop = property.Property
	}
	filter := bson.D{{}}
	update := bson.D{{Key: prop, Value: bson.D{{Key: key, Value: val}}}}
	result, err := db.Team.UpdateMany(context.TODO(), filter, update)
	fmt.Println(result)
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount <= 0 {
		return nil, errors.New("no documents with these properties found")
	}
	return result, nil
}

func (db *DB) UpdateTaskLinksName(task models.Task) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: task.TEAM.ID}, {Key: "tasks._id", Value: task.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "tasks.$.name", Value: task.NAME}}}}
	result, err := db.Team.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount <= 0 {
		return nil, errors.New("no documents with these properties found")
	}
	return result, nil
}
