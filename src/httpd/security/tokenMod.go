package security

import (
	"server/src/helper"
	"time"
)

// for UserSessions
// type UserToken struct {
// 	USERID    primitive.ObjectID
// 	USERNAME  string
// 	RAW_TOKEN *jwt.Token
// 	TOKEN     string
// 	GROUPID   string
// }
// type TokenMap struct {
// 	USERS        []UserToken
// 	TOKENTOINDEX map[string]*UserToken
// 	USERINDEX    map[string]*UserToken
// }

// func (tm *TokenMap) AddUserByUsername(user *UserToken) {

// 	if tm.TOKENTOINDEX == nil {
// 		tm.TOKENTOINDEX = make(map[string]*UserToken)
// 	}

// 	if tm.USERINDEX == nil {
// 		tm.USERINDEX = make(map[string]*UserToken)
// 	}
// 	tm.TOKENTOINDEX[user.TOKEN] = user
// 	tm.USERINDEX[user.USERNAME] = user

// }

// * for Email varification only
type EmailTokenMap struct {
	Keys map[string]string
}

func (EKM *EmailTokenMap) ValidateEmail(key string, email string) bool {
	if EKM.Keys[key] == email {
		delete(EKM.Keys, key)
		return true
	}
	return false
}

type SelectGroupTokenMap struct {
	grupID map[string]Gtoken
}

type Gtoken struct {
	GruID     string
	Timestamp time.Time
}

func (GKM *SelectGroupTokenMap) AddToken(GrupID string) string {
	if GKM.grupID == nil {
		GKM.grupID = make(map[string]Gtoken)
	}
	for k, v := range GKM.grupID {
		if v.GruID == GrupID {
			return k
		}
	}

	var GenerateKey = helper.RandomString(6)
	GKM.grupID[GenerateKey] = Gtoken{GrupID, time.Now()}
	return (GenerateKey)
}
func (GKM *SelectGroupTokenMap) GetGrupID(Key string) (string, int) {
	if GKM.grupID == nil {
		return "", 0
	}
	if GKM.grupID != nil {
		if value, ok := GKM.grupID[Key]; ok {
			return value.GruID, 1
		} else {
			return "", 0
		}
	} else {
		return "", 0
	}
}
