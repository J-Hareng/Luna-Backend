package handler

import (
	"context"
	"fmt"
	"server/src/httpd/security"
	"time"

	"github.com/gin-gonic/gin"
)

type UserNotifierRef struct {
	userID  string
	groupID string
}
type NotifiType struct {
	grupID string
	value  NotifierBody
}
type NotifierBody struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type LunaNotifier struct {
	NewNotifiMsg chan NotifiType

	UserNotifiMsg map[UserNotifierRef]chan interface{}

	UserDisconnect chan UserNotifierRef
}

func NewLunaNotifier() *LunaNotifier {
	return &LunaNotifier{
		NewNotifiMsg:  make(chan NotifiType),
		UserNotifiMsg: make(map[UserNotifierRef]chan interface{}),
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

			case <-time.After(30 * time.Second):

				LN.NotifieClients("BROADCAST", "KeepAlive", "")
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
func (LN *LunaNotifier) NotifieClients(groupID string, t string, b interface{}) {

	LN.NewNotifiMsg <- NotifiType{groupID, NotifierBody{t, b}}

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
