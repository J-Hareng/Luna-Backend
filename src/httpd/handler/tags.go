package handler

import (
	"fmt"
	"net/http"
	"server/src/api/db"
	"server/src/api/db/models"
	"server/src/httpd/bodymodels"

	"github.com/gin-gonic/gin"
)

func EditTagName(DB *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req bodymodels.EditTagMod
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}

		res, err := DB.RemoveTagFromCollection(req.OldName, req.Name, req.CollId)
		if err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR(err.Error()))
			return
		}
		if res == nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR("Failed to remove tag from collection"))
		}

		UpdatedTasks, errU := DB.EditTasksTags(req.Name, req.OldName, res.GrupID)
		if errU != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR("Failed to edit tag in tasks"))
		}
		var response bodymodels.UpdatedTasksMod
		if UpdatedTasks == nil {
			response = bodymodels.UpdatedTasksMod{TaskLs: []models.TaskLink{}}
		} else {
			for _, task := range *UpdatedTasks {
				response.TaskLs = append(response.TaskLs, models.TaskLink{ID: task.ID, NAME: task.NAME})
			}
		}
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(response, ""))
	}
}

func AddTagToCollection(DB *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tag bodymodels.AddTagMod
		err := c.BindJSON(&tag)
		if err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		result, err := DB.AddTag(tag.CollId, tag.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR("Failed to add tag to collection"))
			fmt.Println(err)
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR("Failed to add tag to collection"))
			return
		}

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("Added tag to collection", ""))
	}
}
