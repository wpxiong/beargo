package session

import (
  "github.com/wpxiong/beargo/log"
  "time"
)

func init() {
  log.InitLog()
}

type SessionInfo struct {
  sessionId  string
  sessionInvalidateTime  time.Time
  isInMemory             bool
}

type Session struct {
  sessionId  string
  sessionValue           map[string] interface{}
}


func (session *SessionInfo) UpdateSession(timeOut int){
  session.sessionInvalidateTime = time.Now().Add(time.Second * time.Duration(timeOut))
}

func (session *SessionInfo) InitSession(timeOut int) {
  session.sessionInvalidateTime = time.Now().Add(time.Second * time.Duration(timeOut))
  session.isInMemory = true
}



