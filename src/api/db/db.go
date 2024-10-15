package db

import (
	"context"
	"server/src/api/db/models"
	"server/src/helper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client

	Coll *mongo.Collection
	User *mongo.Collection
	Team *mongo.Collection
	Task *mongo.Collection
}

func New() (*DB, error) {
	u := helper.GetEnvVar("MONGOURI")

	// * Verbinde Mit Der Datenbank
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(u).SetServerAPIOptions(serverAPI)

	// * Setze die maxilae Ladedauer auf 10 Sekunden
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	c, e := mongo.Connect(ctx, opts)
	if e != nil {
		return nil, e
	}

	// * Verbinde Mit den Collectiones
	userCollection := c.Database("LunaDB").Collection("User")
	teamCollection := c.Database("LunaDB").Collection("Team")
	taskCollection := c.Database("LunaDB").Collection("Task")
	CollCollection := c.Database("LunaDB").Collection("Coll")

	return &DB{
		client: c,
		User:   userCollection,
		Team:   teamCollection,
		Task:   taskCollection,
		Coll:   CollCollection,
	}, nil
}
func (db *DB) CheckLastTeam(grupid string) (string, error) {
	cursor, err := db.User.Find(context.TODO(), bson.D{{Key: "groupid", Value: grupid}})
	if err != nil {
		return "", err
	}
	var user []models.User
	if err := cursor.All(context.TODO(), &user); err != nil {
		println(err)
		return "", err
	}
	if user != nil {
		return "", nil
	} else {
		errChan := make(chan error, 2)
		defer close(errChan)

		go func() {
			colls, err := db.GetCollectionsViaGroupID(grupid)
			if err != nil {
				errChan <- err
				return
			}
			for _, coll := range colls {
				_, err := db.RemoveCollection(*coll.ID)
				if err != nil {
					errChan <- err
					return
				}
			}
			errChan <- nil
		}()
		go func() {
			teams, err := db.GetAllTeams(grupid)
			if err != nil {
				errChan <- err
				return
			}
			for _, team := range teams {
				_, err := db.RemoveTeam(models.TeamLink{ID: team.ID})
				if err != nil {
					errChan <- err
					return
				}
			}
			errChan <- nil
		}()

		for i := 0; i < 2; i++ {
			if err := <-errChan; err != nil {
				return "", err
			}
		}

		return "deleted: " + grupid, nil
	}

}
