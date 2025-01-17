package db

import (
	"context"
	"errors"
	"fmt"
	"server/src/api/db/models"
	"server/src/helper"
	"server/src/httpd/security/caches"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) AddUser(name string, email string, password string) (*mongo.InsertOneResult, error) {
	newUser := models.CreateUser(name, email, password, "")
	result, err := db.User.InsertOne(context.TODO(), newUser)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (db *DB) ChangeGrupID(newGruID string, id primitive.ObjectID, t string) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	caches.USERCACHE.Edit(t, "GroupID", newGruID)

	result, err := db.User.UpdateOne(context.TODO(), filter,
		bson.D{{Key: "$set", Value: bson.D{{Key: "groupid", Value: newGruID}}}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

type valid_val_GetUser interface {
	any
}

func (db *DB) GetUser(key string, val valid_val_GetUser) (models.User, error) {
	filter := bson.D{{Key: key, Value: val}}
	cursor, err := db.User.Find(context.TODO(), filter)
	if err != nil {
		return models.User{}, err
	}
	var users []models.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		fmt.Println(err)
		return models.User{}, err
	}
	fmt.Println(users)
	if users == nil {
		return models.User{}, errors.New("no user found")
	}
	return users[0], nil
}

// get all user
func (db *DB) GetAllUsers() ([]models.User, error) {
	filter := bson.D{} // Empty filter to select all documents
	cursor, err := db.User.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	fmt.Println(users)
	return users, nil
}

//TODO check if user name is avalibe

func (db *DB) AvalabileEmail(email string) (bool, error) {
	usercoll := db.User
	filter := bson.D{{Key: "email", Value: email}}
	cursor, err := usercoll.Find(context.TODO(), filter)
	if err != nil {

		return false, err
	}

	var users []models.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		println(err)
		return false, err

	}
	if len(users) != 0 {
		return true, nil
	}
	return false, nil
}

func (db *DB) UserTeams(user models.UserLink) ([]models.Team, error) {
	filter := bson.M{"users": bson.M{"$elemMatch": user}}
	cursor, err := db.Team.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var teams []models.Team
	if err := cursor.All(context.TODO(), &teams); err != nil {
		println(err)
		return nil, err
	}
	fmt.Println(teams)
	return teams, nil

}

func (db *DB) RemoveUser(user models.UserLink) (*mongo.DeleteResult, error) {

	db.RemoveFromMultibleTeamArray("users", user, helper.InsertOption[string]{Property: ""})
	//! TASK AND TEAMS NEED TO BE REMOVED

	filter := bson.D{{Key: "_id", Value: user.ID}}
	result, err := db.User.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *DB) UpdateUserArrays(user models.UserLink, key string, prop string, value any) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: prop, Value: bson.D{{Key: key, Value: value}}}}
	result, err := db.User.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount <= 0 {
		return nil, errors.New("no documents with these properties found")
	}
	return result, nil
}
func (db *DB) GetAllUsersInTeam(groupid string) ([]models.OtherUsers, error) {
	filter := bson.D{{Key: "groupid", Value: groupid}}
	res, err := db.User.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []models.OtherUsers
	if err := res.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil
}
func (db *DB) UpdateUserTasks(task models.Task, user models.UserLink) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "tasks.$.name", Value: task.NAME}}}}
	result, err := db.User.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
