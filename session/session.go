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


func (this *Session) InitSession(id string){
  this.sessionId = id
  this.sessionValue = make(map[string] interface{})
}

func (this *SessionInfo) UpdateSession(timeOut int){
  this.sessionInvalidateTime = time.Now().Add(time.Second * time.Duration(timeOut))
}

func (this *SessionInfo) InitSession(timeOut int) {
  this.sessionInvalidateTime = time.Now().Add(time.Second * time.Duration(timeOut))
  this.isInMemory = true
}

func (this *Session) SaveSessionValue(valueId string, obj interface{}){
  this.sessionValue[valueId] = obj
}

func (this *Session) GetSessionValue(valueId string)  interface{} {
  return this.sessionValue[valueId] 
}

