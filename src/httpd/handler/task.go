package handler

import (
	"context"
	"fmt"
	"net/http"
	"server/src/api/db"
	"server/src/api/db/models"
	"server/src/helper"
	"server/src/httpd/bodymodels"
	"server/src/httpd/handler/stream"
	"server/src/httpd/security"

	"github.com/gin-gonic/gin"
)

func AddTask(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var Task bodymodels.AddTaskMod

		if err := c.ShouldBindJSON(&Task); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		//! need to be checkt if admin and name already used
		_, err := db.AddTask(Task.Name, Task.Description, Task.Team, Task.Groupid, Task.CollID, Task.Tags)
		helper.CustomErrorApi(c, err)

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("task inserted", ""))
	}
}
func RemoveTask(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var Task bodymodels.RemoveTaskMod

		if err := c.ShouldBindJSON(&Task); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusConflict, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}

		helper.Prittylog("body", Task, "RemoveTask 1 ")
		//! need to be checkt if admin and name already used
		_, err := db.RemoveTask(Task.TeamId, Task.Tasks)

		helper.CustomErrorApi(c, err)

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("task removed", ""))
	}
}
func EditTask(db *db.DB, LN *stream.LunaNotifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		var Task models.Task
		if err := c.ShouldBindJSON(&Task); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusConflict, bodymodels.NewLunaResponse_ERROR("Invalid Payload"))
			return
		}

		UpdateRes, err := db.EditTask(Task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Failed to update Task"))
			return
		}
		if UpdateRes.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, bodymodels.NewLunaResponse_ERROR("Task not found"))
			return
		}
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("UpdatedTask", ""))
		LN.NotifieTaskUpdated(Task.GROUPID, Task.ID)
	}

}
func AssingForTask(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var Task bodymodels.EditAssingmentForTaskMod
		if err := c.ShouldBindJSON(&Task); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusConflict, bodymodels.NewLunaResponse_ERROR("Invalid Payload"))
			return
		}
		db.AssingForTask(Task.Task, Task.User)
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("itemInserted", ""))
	}
}
func RemoveAssingForTask(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type bodys struct {
			TASK models.TaskLink `json:"task"`
		}
		var body bodys
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusConflict, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		user, ok := security.DecodeUserFrom_C(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("User not found"))
			return
		}
		_, err := db.UnasingForTask(body.TASK, user.UserLink)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Failed to remove user From Task"))
			return
		}

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("Updated", ""))
	}
}
func GetTasksArray(ctx context.Context, db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqVal bodymodels.GetAllTasksFromTeamMod
		if err := c.ShouldBindJSON(&reqVal); err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		fmt.Println(reqVal)

		task, err := db.GetTaskArray(reqVal.Tasks, ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Failed to retrieve users"))
			return
		}
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(task, ""))
	}
}

// func UpdateTeamTask(db *db.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var Task bodymodels.UpdateTaskRequestMod
// 		if err := c.ShouldBindJSON(&Task); err != nil {
// 			c.JSON(http.StatusConflict)
// 			return
// 		}

// 		Tasks := db.GetTas

// 	}
// }
