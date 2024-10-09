package handler

import (
	"fmt"
	"net/http"
	"server/src/api/db"
	"server/src/httpd/bodymodels"

	"github.com/gin-gonic/gin"
)

func AddCollection(DB *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var coll bodymodels.AddCollectionMod
		err := c.BindJSON(&coll)
		if err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		_, err = DB.AddCollection(coll.Name, coll.Groupid)
		if err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR("Failed to add collection to DB"))
			return
		}
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("Collection added to DB", ""))
	}
}

func GetCollections(DB *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqVal bodymodels.GetAllCollectionsMod
		err := c.BindJSON(&reqVal)
		if err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		collections, err := DB.GetCollectionsViaGroupID(reqVal.Groupid)

		// Fix untill we have a better solution
		fmt.Println(collections)

		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Failed to get collections From DB"))
			return
		}
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(collections, ""))
	}
}

func RemoveCollection(DB *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var coll bodymodels.RemoveCollectionMod
		err := c.BindJSON(&coll)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := DB.RemoveCollection(coll.CollId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func AddTagToCollection(DB *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tag bodymodels.AddTagMod
		err := c.BindJSON(&tag)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Tag : ")
		fmt.Println(tag)
		result, err := DB.AddTag(tag.CollId, tag.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
