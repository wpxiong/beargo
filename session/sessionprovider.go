package session

import (
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}


type SessionProvider interface {
   CreateSession(sessionId string) (Session, error)
   DeleteSession(sessionId string) error
   RemoveAllSession(sessionIdList []string) error
   FindSessionById(sessionId string) bool
   StoreSessionById(sessionId string)  error
   StoreSession(sessionId []string)  error
   LoadSessionById(sessionId string) (Session ,error)
   LoadSession(sessionId []string)  (map[string]Session ,error)
}


