package handler

import (
	"context"
	"fmt"
	"net/http"
	email "server/src/api/Email"
	"server/src/api/db"
	"server/src/api/db/models"
	"server/src/helper"
	"server/src/httpd/bodymodels"
	"server/src/httpd/security"
	"server/src/httpd/security/caches"
	"strings"

	"github.com/gin-gonic/gin"
)

// * User Management

func GetUsers(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, err := db.GetAllUsers() // * alle nutzer aus datenbank laden

		if err != nil { // * error handeling
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}
		c.JSON(http.StatusOK, val) // * nuter zurück senden
	}
}
func GetUserSalt(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var email struct {
			Email string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&email); err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		salt, err := db.GetUserSalt(email.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Error while getting salt"))
			return
		}

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(salt, ""))
	}
}

func AddUser(db *db.DB, EKM *security.EmailTokenMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user bodymodels.AddUserMod

		if err := c.ShouldBindJSON(&user); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		is_used, err := db.AvalabileEmail(user.Email)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "DB not Found"})
		}

		if is_used {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already used"})
			return
		}

		if !EKM.ValidateEmail(user.Key, user.Email) {
			c.JSON(http.StatusConflict, gin.H{"error": "Wrong validation key"})
			return
		}

		email := strings.ToLower(user.Email)
		db.AddUser(user.Name, email, user.Password, user.Salt)
		c.JSON(http.StatusOK, user)
	}
}

func GetUserViaToken(db *db.DB) gin.HandlerFunc { // * is save (Behind middleware)
	return func(c *gin.Context) {
		u, _ := security.DecodeUserFrom_C(c)

		user, err := db.GetUser("_id", u.UserLink.ID)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password wrong"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

// * Email verification
func RequestEmailKey(e email.Email, EKM *security.EmailTokenMap, db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var msg bodymodels.RequestEmailKeyMod
		if err := c.ShouldBindJSON(&msg); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			helper.CustomError(err.Error())
			return
		}
		is_used, err := db.AvalabileEmail(msg.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error loading db"})
			return
		}

		if is_used {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already used"})
			return
		}

		key, err := helper.GenerateKey(6)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errorInKeyGenerator": err.Error()})
			helper.CustomError(err.Error())
			return
		}
		if EKM.Keys == nil {
			EKM.Keys = make(map[string]string, 0)
		}

		EKM.Keys[key] = msg.Email

		fmt.Println(key)
		e.SendEmail("<p> your key is </p> <h2>"+key+"</h2>", "Email varification", msg.Email)
		c.JSON(http.StatusOK, gin.H{"message": "email send", "email": msg.Email})
	}
}

// * Session handeling
func RequestSessionToken(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin bodymodels.RequestSessionTokenMod
		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			helper.CustomErrorApi(c, err)
			return
		} // * lade Post Models in Die Variable userLogin
		fmt.Println(userLogin)
		user, err := db.GetUser("email", strings.ToLower(userLogin.Email)) // * Finde Die Email
		if err != nil {                                                    //* error Handeling
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("Email or password wrong"))
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Println(user)
		if userLogin.Password != user.PASSWORD { // * vergleiche Email
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("Email or password wrong"))

			return
		}

		// * erstelle token und speichere sie in der TM ab
		token, err := security.CreateUserToken(models.UserLink{NAME: user.NAME, ID: user.ID}, user.GROUPID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Error while creating token"))
			return
		}

		//* Antwort an client
		c.SetSameSite(http.SameSiteLaxMode)
		// c.SetCookie("token", token, 3600, "/", "", false, true)
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(gin.H{"token": token}, ""))
	}
}
func RemoveSessionToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.SetCookie("token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Token removed"})

	}
}
func ValidateUserToken(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if caches.USERCACHE.IsNil() {
			fmt.Println("UserCache is nil")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Token available in cache"})
			c.Redirect(http.StatusFound, "http://localhost:8080/")
			return
		}

		token, err := helper.GetVerfToken(c)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Valid": false})
			return
		}

		if _, ok := caches.USERCACHE.Get(token); ok {

			c.JSON(http.StatusOK, gin.H{"Valid": true})

			return
		} else {

			c.JSON(http.StatusOK, gin.H{"Valid": false})

		}
	}
}

// * User team handeling
func GetTeamsFromUser(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var User caches.CacheUserData
		if err := security.ShouldDecodeFrom_C(c, &User); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		teams, err := db.UserTeams(User.UserLink)
		helper.CustomErrorApi(c, err)
		c.JSON(http.StatusOK, teams)
	}
}

func RemoveUser(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.UserLink
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			helper.CustomError(err.Error())
			return
		}
		val, err := db.RemoveUser(user)
		helper.CustomErrorApi(c, err)
		c.JSON(http.StatusOK, gin.H{"message": val})
	}
}

func AddGruopKey(db *db.DB, GMK *security.SelectGroupTokenMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		var val bodymodels.AddGrupIDTokenMod
		if err := c.ShouldBindJSON(&val); err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR(err.Error()))
			helper.CustomError(err.Error())
			return
		}
		user, err := db.GetUser("_id", val.User.ID)
		helper.CustomErrorApi(c, err)

		token := GMK.AddToken(user.GROUPID)
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(token, ""))
	}
}
func JoinGruopViaKey(db *db.DB, GKM *security.SelectGroupTokenMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		var val bodymodels.GetGrupIDTokenMod
		if err := c.ShouldBindJSON(&val); err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR("Error While Binding JSON"))
			helper.CustomError(err.Error())
			return
		}

		cUser, ok := security.DecodeUserFrom_C(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("No Token available"))
			return
		}

		usertoken, err := helper.GetVerfToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("No Token available"))
			return
		}

		groupId, status := GKM.GetGrupID(val.Key)
		if status == 0 {
			c.JSON(http.StatusNotFound, bodymodels.NewLunaResponse_ERROR("Key not found"))
			return
		} else {
			db.ChangeGrupID(groupId, cUser.UserLink.ID, usertoken)
		}

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(groupId, ""))
	}
}
func SetGupID(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var val bodymodels.AddGrupIDMod
		if err := c.ShouldBindJSON(&val); err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR(err.Error()))
			helper.CustomError(err.Error())
			return
		}

		cUser, ok := security.DecodeUserFrom_C(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("No Token available"))
			return
		}

		token, err := helper.GetVerfToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("No Token available"))
			return
		}

		key, err := helper.GenerateKey(6)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Error While Generating Key"))
			helper.CustomError(err.Error())
			return
		}
		NewGrupID := val.GroupidName + key
		_, err = db.ChangeGrupID(NewGrupID, cUser.UserLink.ID, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Error While Changing GrupID"))
			helper.CustomError(err.Error())
			return
		}
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(NewGrupID, ""))
	}
}
func RemoveGrupID(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userOutdated, ok := security.DecodeUserFrom_C(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("No Token available"))
			return
		}

		usertoken, err := helper.GetVerfToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("No Token available"))
			return
		}
		user, err := db.GetUser("_id", userOutdated.UserLink.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Failed to retrieve users"))
			return
		}

		userLink := models.UserLink{ID: user.ID, NAME: user.NAME}
		if user.TASKS != nil {
			for _, tasks := range user.TASKS {
				_, err = db.UnasingForTask(tasks, userLink)
				if err != nil {
					c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Error While Unasing For Task"))
					return
				}
			}
		}
		_, err = db.RemoveUserFromAllTeamsInGroup(userLink, user.GROUPID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Error While Removeing FromTeams"))
			return
		}
		_, err = db.ChangeGrupID("", user.ID, usertoken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Error While Changing GrupID"))
			return
		}
		_, err = db.CheckLastTeam(user.GROUPID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("Error While Checking if the team needs to be removed"))
			return
		}
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK("Removed", ""))
	}
}

func GetAllUserData(ctx context.Context, db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("------------------start get all user data------------------")
		// var wg sync.WaitGroup

		var u caches.CacheUserData
		if err := security.ShouldDecodeFrom_C(c, &u); err != nil {
			c.JSON(http.StatusUnauthorized, bodymodels.NewLunaResponse_ERROR("User not found"))
			return
		}
		User, err := db.GetUser("_id", u.UserLink.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR("User not found in db"))
			return
		}
		if User.GROUPID == "" {
			c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(gin.H{"user": User}, "NoID"))
			return
		}
		errorChannel := make(chan error, 3)
		defer close(errorChannel)

		collectionChannel := make(chan []models.Collection)
		collectionChannelOK := make(chan bool)
		defer close(collectionChannel)
		defer close(collectionChannelOK)

		teamChannel := make(chan []models.Team)
		teamChannelOK := make(chan bool)
		defer close(teamChannel)
		defer close(teamChannelOK)

		teamTaskChannel := make(chan []models.Task)
		teamTaskChannelOK := make(chan bool)
		defer close(teamTaskChannel)
		defer close(teamTaskChannelOK)

		UsersChannel := make(chan []models.PublicUser)
		UsersChannelOK := make(chan bool)
		defer close(UsersChannel)
		defer close(UsersChannelOK)

		// wg.Add(3)
		go func() {
			// defer wg.Done()
			teamTask, err := db.GetAllTasksInGroup(u.GroupID)
			if err != nil {
				errorChannel <- err
				return
			}
			if len(teamTask) == 0 {
				teamTaskChannelOK <- false
				return
			}
			teamTaskChannelOK <- true
			teamTaskChannel <- teamTask
		}()

		go func() {
			// defer wg.Done()
			collection, err := db.GetCollectionsViaGroupID(User.GROUPID)
			if err != nil {
				errorChannel <- err
				return
			}
			if len(collection) == 0 {
				collectionChannelOK <- false
				return
			}
			collectionChannelOK <- true
			collectionChannel <- collection

		}()

		go func() {
			// defer wg.Done()
			teams, err := db.GetAllTeams(u.GroupID)
			if err != nil {
				errorChannel <- err
				return
			}
			if len(teams) == 0 {
				teamChannelOK <- false
				return
			}
			teamChannelOK <- true
			teamChannel <- teams

		}()
		go func() {
			users, err := db.GetAllUsersInTeam(User.GROUPID)
			if err != nil {
				errorChannel <- err
				return
			}
			if len(users) == 0 {
				UsersChannelOK <- false
				return
			}
			UsersChannelOK <- true
			UsersChannel <- users
		}()

		var collection []models.Collection = []models.Collection{}
		var teams []models.Team = []models.Team{}
		var teamTask []models.Task = []models.Task{}
		var users []models.PublicUser = []models.PublicUser{}

		for i := 0; i < 4; i++ {
			select {
			case err := <-errorChannel:
				fmt.Printf("Error: %v", err)
				c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR(err.Error()))
				i = 10
				return
			case ok := <-collectionChannelOK:
				if ok {
					collection = <-collectionChannel
				}
			case ok := <-teamChannelOK:
				if ok {
					teams = <-teamChannel
				}

			case ok := <-teamTaskChannelOK:
				if ok {
					teamTask = <-teamTaskChannel
				}
			case ok := <-UsersChannelOK:
				if ok {
					users = <-UsersChannel
				}
			}
		}
		fmt.Println("Done")
		body := bodymodels.GetAllUserDataMod{
			User:        User,
			Teams:       teams,
			Collections: collection,
			UserTasks:   teamTask,
			OtherUsers:  users,
		}

		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(body, ""))
	}
}
