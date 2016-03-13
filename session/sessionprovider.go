package session

import (
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}


type SessionProvider interface {
   InitProvider()
   CreateSession(sessionId string) (Session, error)
   DeleteSession(sessionId string) error
   FindSessionById(sessionId string) bool
   LoadSessionById(sessionId string) (Session ,error)
}


