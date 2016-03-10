package session

import (
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}


type SessionProvider interface {
   CreateSession(sessionId string) (Session, error)
   LoadSession(sessionId string) (Session, error)
   DeleteSession(sessionId string) error
   RemoveAllSession(timeout int64) error
}


