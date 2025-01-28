package stream

var (
	// Operations
	DeleteOperation    = "DELETE"
	UpdateOperation    = "UPDATE"
	AddOperation       = "ADD"
	SyncOperation      = "SYNC"
	KeepAliveOperation = "KeepAlive"

	// Types
	TaskObject = "Task"
	TeamObject = "Team"
	UserObject = "User"
	FileObject = "FILE"
	IdObjectID = "ID"

	NONE      = "NONE"
	BROADCAST = "BROADCAST"
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
	Operation string      `json:"operation"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
}
