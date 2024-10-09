package handler

import (
	"fmt"
	"net/http"
	"server/src/api/db"
	"server/src/helper"
	"server/src/httpd/bodymodels"
	"server/src/httpd/security"
	"server/src/httpd/security/caches"

	"github.com/gin-gonic/gin"
)

func AddTeam(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var team bodymodels.AddTeamMod

		if err := c.ShouldBindJSON(&team); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}

		_, err := db.AddTeam(team.Name, team.User, team.Des, team.Groupid)
		helper.CustomErrorApi(c, err)
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("Team Added", ""))
	}
}

func GetTeams(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var User caches.CacheUserData
		if err := security.ShouldDecodeFrom_C(c, &User); err != nil {
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		teams, err := db.GetAllTeams(User.GroupID)
		if err != nil {
			// Handle the error appropriately
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Failed to retrieve users"))
			return
		}
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(teams, ""))
	}
}

func AddUserToTeam(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqVal bodymodels.AddUserToTeamMod

		if err := c.ShouldBindJSON(&reqVal); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		_, err := db.AddUserToTeam(reqVal.Team, reqVal.User)
		helper.CustomErrorApi(c, err)

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("User Added", ""))

	}

}
func RemoveUserFromTeam(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqVal bodymodels.AddUserToTeamMod

		if err := c.ShouldBindJSON(&reqVal); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}

		_, err := db.RemoveUserFromTeam(reqVal.Team, reqVal.User)
		helper.CustomErrorApi(c, err)

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("User Removed", ""))
	}
}
func RemoveTeam(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var team bodymodels.RemoveTeamMod
		if err := c.ShouldBindJSON(&team); err != nil {
			c.JSON(http.StatusConflict, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		_, err := db.RemoveTeam(team.Team)
		helper.CustomErrorApi(c, err)
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("Team Deleted", ""))
	}
}

func EditTeam(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqVal bodymodels.EditTeamMod

		if err := c.ShouldBindJSON(&reqVal); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}

		_, err := db.EditTeam(reqVal.Team)
		helper.CustomErrorApi(c, err)

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("Updated", ""))
	}
}
