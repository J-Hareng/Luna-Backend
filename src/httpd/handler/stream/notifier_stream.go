package stream

import (
	"context"
	"fmt"
	"server/src/api/db"
	"server/src/httpd/security"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LunaNotifier struct {
	NewNotifiMsg   chan NotifiType
	UserNotifiMsg  map[UserNotifierRef]chan interface{}
	UserDisconnect chan UserNotifierRef
	DbRev          *db.DB
}

func NewLunaNotifier(db *db.DB) *LunaNotifier {
	return &LunaNotifier{
		NewNotifiMsg:  make(chan NotifiType),
		UserNotifiMsg: make(map[UserNotifierRef]chan interface{}),
		DbRev:         db,
	}
}
func (LN *LunaNotifier) run(ctx context.Context) {

	defer close(LN.NewNotifiMsg)
	//Keep Alive
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case <-time.After(5 * time.Second):

				LN.NotifieClients(BROADCAST, KeepAliveOperation, NONE, NONE)
			}
		}
	}()

	for {
		select {
		case msg := <-LN.NewNotifiMsg:

			for user, ch := range LN.UserNotifiMsg {
				if user.groupID == msg.grupID || msg.grupID == "BROADCAST" {
					ch <- msg.value
				}
			}

		case user := <-LN.UserDisconnect:

			println("User Disconnected Deleting Now")
			close(LN.UserNotifiMsg[user])
			delete(LN.UserNotifiMsg, user)

			return
		case <-ctx.Done():
			for _, ch := range LN.UserNotifiMsg {
				close(ch)
			}
			return
		}
	}
}
func (LN *LunaNotifier) Start(ctx context.Context) {
	go LN.run(ctx)
}

func (LN *LunaNotifier) NotifieClients(groupID string, operation string, typ string, body interface{}) {

	LN.NewNotifiMsg <- NotifiType{groupID, NotifierBody{operation, typ, body}}

}

func (LN *LunaNotifier) AddUserToSteam(ctx context.Context, c *gin.Context) {

	ch := make(chan interface{})

	user, ok := security.DecodeUserFrom_C(c)
	if !ok {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	userRef := UserNotifierRef{user.UserLink.ID.String(), user.GroupID}
	LN.UserNotifiMsg[userRef] = ch

	// Set the headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.SSEvent("connected", "Connected to the server")
	// Flush the headers
	c.Writer.Flush()

	// Goroutine to handle sending events to the client
	func() {
		defer func() {
			LN.UserDisconnect <- userRef
		}()

		for {
			select {
			case msg := <-ch:
				c.SSEvent("message", msg)
				fmt.Println("Message sent")
				c.Writer.Flush()

			case <-c.Request.Context().Done():

				println("User Disconnected from the server")
				return

			case <-ctx.Done():
				print("Server is shutting down")
				return
			}
		}
	}()
}

func ConnectToNotifierStream(ctx context.Context, LN *LunaNotifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		LN.AddUserToSteam(ctx, c)
	}
}

// * Notifier Macros
func (LN *LunaNotifier) NotifieAllClients(t string, b interface{}) {
	LN.NotifieClients(BROADCAST, NONE, t, b)
}

func (LN *LunaNotifier) NotifieTaskAdded(gId string, taskID primitive.ObjectID) {
	b, err := LN.DbRev.GetTask(taskID)
	if err != nil {
		fmt.Println(err)
		return
	}
	LN.NotifieClients(gId, AddOperation, TaskObject, b)
}
func (LN *LunaNotifier) NotifieTaskDeleted(gId string, taskID primitive.ObjectID) {
	LN.NotifieClients(gId, DeleteOperation, IdObjectID, taskID)
}
func (LN *LunaNotifier) NotifieTaskUpdated(gId string, taskID primitive.ObjectID) {
	b, err := LN.DbRev.GetTask(taskID)
	if err != nil {
		fmt.Println(err)
		return
	}
	LN.NotifieClients(gId, UpdateOperation, TaskObject, b)
}
