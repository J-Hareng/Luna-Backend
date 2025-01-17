package httpd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	email "server/src/api/Email"
	"server/src/api/db"
	filemanagement "server/src/api/file_management"
	"server/src/httpd/handler"
	"server/src/httpd/security"

	// "github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"

	"github.com/gin-gonic/gin"
)

type Server struct {
	db *db.DB

	engine *gin.Engine
}

func Init(ctx context.Context, DB *db.DB, E email.Email, EKM *security.EmailTokenMap, GKM *security.SelectGroupTokenMap, S3 *filemanagement.S3Client) Server {

	r := gin.Default() // * Initialisire End-Punkt

	// * CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "localhost:4200/")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// * WEB
	r.Use(static.Serve("/", static.LocalFile("client/browser", true)))
	r.Use(static.Serve("/Dashboard", static.LocalFile("client/browser", true)))

	// * without middleware
	r.GET("/validateSessionToken", handler.ValidateUserToken(DB))

	r.POST("/addUser", handler.AddUser(DB, EKM))
	r.POST("/reqSessionToken", handler.RequestSessionToken(DB))

	r.POST("/reqEmailKey", handler.RequestEmailKey(E, EKM, DB))

	// * With middleware
	secureGroup := r.Group("/api")             // * gruppe erstellen
	secureGroup.Use(security.ConditionToken()) // * middleware einbinden
	{
		// * testAuth
		secureGroup.GET("/testAuth", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello from the server"})
		})

		// * grupmanagement
		secureGroup.POST("/addGruopKey", handler.AddGruopKey(DB, GKM))
		secureGroup.POST("/joinGroupViaKey", handler.JoinGruopViaKey(DB, GKM))
		secureGroup.POST("/setGupID", handler.SetGupID(DB))
		secureGroup.GET("/removeGrupID", handler.RemoveGrupID(DB))

		// * user
		secureGroup.GET("/getUserViaToken", handler.GetUserViaToken(DB))
		secureGroup.GET("/getUsers", handler.GetUsers(DB))
		// secureGroup.GET("/removeSessionToken", handler.RemoveSessionToken())
		secureGroup.GET("/getAllUserData", handler.GetAllUserData(ctx, DB))

		secureGroup.POST("/RemoveUser", handler.RemoveUser(DB))
		secureGroup.GET("/getUserTeams", handler.GetTeamsFromUser(DB))
		// * teams
		secureGroup.POST("/editTeam", handler.EditTeam(DB))
		secureGroup.GET("/getTeams", handler.GetTeams(DB))
		secureGroup.POST("/RemoveTeam", handler.RemoveTeam(DB))
		secureGroup.POST("/addTeam", handler.AddTeam(DB))

		secureGroup.POST("/AddUserToTeam", handler.AddUserToTeam(DB))
		secureGroup.POST("/RemoveUserFromTeam", handler.RemoveUserFromTeam(DB))

		//* tasks
		secureGroup.POST("/AddTask", handler.AddTask(DB))
		secureGroup.POST("/RemoveTask", handler.RemoveTask(DB))
		secureGroup.POST("/ApplyForTask", handler.AssingForTask(DB))
		secureGroup.POST("/RevokeTaskApplication", handler.RemoveAssingForTask(DB))

		secureGroup.POST("/GetTaskArray", handler.GetTasksArray(ctx, DB))
		secureGroup.POST("/EditTask", handler.EditTask(DB))

		//* Collections
		secureGroup.POST("/AddCollection", handler.AddCollection(DB))
		secureGroup.POST("/GetCollections", handler.GetCollections(DB))
		secureGroup.POST("/RemoveCollection", handler.RemoveCollection(DB))
		secureGroup.POST("/AddTagToCollection", handler.AddTagToCollection(DB))
		secureGroup.POST("/EditTagInCollection", handler.EditTagName(DB))

		//* Files
		secureGroup.POST("/UploadFile", handler.UploadFile(DB, S3))
		secureGroup.POST("/GetFile", handler.GetFile(DB, S3))
		secureGroup.POST("/GetFiles", handler.GetMultibleFiles(DB, S3))
		secureGroup.POST("/DownloadFile", handler.DownloadFile(DB, S3))
		secureGroup.POST("/RemoveFile", handler.RemoveFile(DB, S3))

	}
	return Server{
		engine: r,
		db:     DB,
	}
}

func (s Server) Run() {
	port := os.Getenv("PORT")
	fmt.Println("Running now")
	s.engine.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for
}
