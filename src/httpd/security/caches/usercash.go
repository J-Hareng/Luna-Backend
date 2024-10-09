package caches

import (
	"server/src/api/db/models"
	"time"

	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var USERCACHE UserCache = NewFromEnv()

type CacheUserData struct {
	UserLink models.UserLink
	GroupID  string
}

type UserCache struct {
	cache *cache.Cache
}

func NewFromEnv() UserCache {
	c := cache.New(24*time.Hour, cache.NoExpiration)
	cu := UserCache{cache: c}
	cu.Set("a", models.UserLink{ID: primitive.NewObjectID(), NAME: "a"}, "a")
	return cu
}
func (uc *UserCache) Set(UserKeyString string, UL models.UserLink, UGID string) bool {
	CUD := CacheUserData{UserLink: UL, GroupID: UGID}

	uc.cache.Set(UserKeyString, CUD, cache.DefaultExpiration)
	return true
}
func (uc *UserCache) Get(UserKeyString string) (CacheUserData, bool) {
	v, ok := uc.cache.Get(UserKeyString)
	return v.(CacheUserData), ok
}
func (uc *UserCache) Delete(k string) bool {
	uc.cache.Delete(k)
	return true
}
func (uc *UserCache) Len() int {
	items := uc.cache.Items()
	return len(items)
}

func (uc *UserCache) IsNil() bool {
	return uc.cache == nil
}
func (uc *UserCache) Edit(k string, Entrykey string, Value any) (CacheUserData, bool) {
	u_raw, expTime, ok := uc.cache.GetWithExpiration(k)
	if !ok {
		return CacheUserData{}, false
	}
	u := u_raw.(CacheUserData)

	switch Entrykey {
	case "GroupID":
		u.GroupID = Value.(string) // replace GroupIDType with the actual type
	default:
		return CacheUserData{}, false
	}
	// replace with actual method to get TTL
	remainingDuration := time.Until(expTime)
	uc.cache.Replace(k, u, remainingDuration)
	return u, true
}
