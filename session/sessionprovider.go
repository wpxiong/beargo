package session

import (
  "github.com/wpxiong/beargo/log"
  "sync"
)

func init() {
  log.InitLog()
}


type SessionProvider interface {
   InitProvider(SessionLifeTime int64)
   CreateSession(sessionId string) (Session, error)
   DeleteSession(sessionId string) error
   FindSessionById(sessionId string) bool
   LoadSessionById(sessionId string) (Session ,error)
   SerializeSession(sessionTypeInfo  map[string] interface{})
   DeserializeSession(sessionTypeInfo  map[string] interface{})
   ClearSession(sessionAccess  *sync.Mutex)
}


