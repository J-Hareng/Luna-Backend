package db

import (
	"context"
	"fmt"
	"server/src/api/db/models"
	"server/src/helper"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) AddTask(name string, des string, team models.TeamLink, gruID string, collID string, tags []string) (*mongo.InsertOneResult, error) {
	newTask := models.CreateTask(name, des, team, gruID, collID, tags)
	result, err := db.Task.InsertOne(context.TODO(), newTask)
	if err != nil {
		return nil, err
	}
	var insertedID primitive.ObjectID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		insertedID = oid
	} else {
		// Handle the case where InsertedID is not of type primitive.ObjectID
		return nil, fmt.Errorf("unexpected type for InsertedID: %T", result.InsertedID)
	}
	link := models.TaskLink{
		NAME: name,
		ID:   insertedID,
	}
	updateres, err := db.InsertInInTeamArray(team.ID, "tasks", link, helper.InsertOption[string]{Property: ""})
	if err != nil {
		fmt.Println(err)
		fmt.Println(updateres)
		return nil, err
	}
	fmt.Println(updateres)
	return result, err
}
func (db *DB) RemoveTasks(tasks []models.TaskLink) (*mongo.DeleteResult, error) {
	idArray := make([]primitive.ObjectID, len(tasks))
	for i, task := range tasks {
		idArray[i] = task.ID
	}
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: idArray}}}}
	result, err := db.Task.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (db *DB) EditTask(task models.Task) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: task.ID}}
	update := bson.D{{Key: "$set", Value: task}}
	result, err := db.Task.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *DB) RemoveTask(id primitive.ObjectID, tasks models.TaskLink) (*mongo.DeleteResult, error) {

	filter := bson.D{{Key: "_id", Value: tasks.ID}}
	result, err := db.Task.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	res, err := db.RemoveFromTeamArray(id, "tasks", tasks, helper.InsertOption[string]{Property: ""})
	if err != nil {
		return nil, err
	}

	fmt.Println(res)
	return result, nil
}
func (db *DB) AssingForTask(task models.TaskLink, user models.UserLink) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: task.ID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "inprogress", Value: user}}}}
	result, err := db.Task.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	db.UpdateUserArrays(user, "tasks", "$push", task)
	return result, nil
}
func (db *DB) UnasingForTask(task models.TaskLink, user models.UserLink) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: task.ID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "inprogress", Value: user}}}}
	result, err := db.Task.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	db.UpdateUserArrays(user, "tasks", "$pull", task)
	return result, nil
}

func (db *DB) EditTasksTags(TagNew string, tagOld string, groupID string) (*[]models.Task, error) {
	filter := bson.D{
		{Key: "groupid", Value: groupID},
		{Key: "tags", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "$eq", Value: tagOld}}}}},
	}

	filterForFindingIDs := bson.D{
		{Key: "groupid", Value: groupID},
		{Key: "tags", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "$eq", Value: TagNew}}}}},
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "tags.$", Value: TagNew}}}}
	result, err := db.Task.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, nil
	}
	UpdatedTasks, err := db.FindTaskWithFilter(filterForFindingIDs)

	if err != nil {
		return nil, err
	}
	return UpdatedTasks, nil

}

func (db *DB) GetTaskArray(taskArray []models.TaskLink, ctx context.Context) ([]models.Task, error) {
	if len(taskArray) == 0 {
		return nil, fmt.Errorf("no tasks given")
	}

	retunTaskArray := make([]models.Task, len(taskArray))

	// wg := sync.WaitGroup{}

	TaskChan := make(chan models.Task, len(taskArray))
	errChan := make(chan error, len(taskArray))
	defer close(TaskChan)
	defer close(errChan)

	fmt.Println(taskArray)
	for _, task := range taskArray {

		go func(task models.TaskLink, errchan chan error, taskchan chan models.Task, wgl *sync.WaitGroup) {

			filter := bson.D{{Key: "_id", Value: task.ID}}

			var taskFound models.Task
			err := db.Task.FindOne(ctx, filter).Decode(&taskFound)

			if err != nil {
				errchan <- err
				return
			}
			taskchan <- taskFound
		}(task, errChan, TaskChan, nil)
	}

	for i := 0; i < len(taskArray); i++ {
		select {
		case task := <-TaskChan:
			fmt.Println(task)
			retunTaskArray = append(retunTaskArray, task)

		case err := <-errChan:

			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return retunTaskArray, nil
}
func (db *DB) GetUsersTask(userId primitive.ObjectID) (map[string]models.Task, error) {

	filter := bson.M{"Users": bson.M{"$elemMatch": bson.M{"_id": userId}}}
	cursor, err := db.Task.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	tasksResponseMod := make(map[string]models.Task)
	for cursor.Next(context.TODO()) {
		var t models.Task
		err := cursor.Decode(&t)

		if err != nil {
			return nil, err
		}
		tasksResponseMod[t.ID.Hex()] = t
	}
	return tasksResponseMod, nil

}
func (db *DB) GetAllTasksInGroup(GroupID string) ([]models.Task, error) {
	filter := bson.M{"groupid": GroupID}
	cursor, err := db.Task.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var task []models.Task
	cursor.All(context.TODO(), &task)
	return task, nil
}

func (db *DB) FindTaskWithFilter(filter primitive.D) (*[]models.Task, error) {
	var task models.Task

	cursor, err := db.Task.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}
	var tasks []models.Task
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return &tasks, nil
}
