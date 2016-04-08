package session

import (
  "github.com/wpxiong/beargo/log"
  "time"
)

func init() {
  log.InitLog()
}

type SessionInfo struct {
  SessionId  string
  SessionInvalidateTime  time.Time
  IsSerialized             bool
}

type Session struct {
  SessionId  string
  SessionValue           map[string] interface{}
}


func (this *Session) InitSession(id string){
  this.SessionId = id
  this.SessionValue = make(map[string] interface{})
}

func (this *SessionInfo) UpdateSession(timeOut int){
  this.SessionInvalidateTime = time.Now().Add(time.Second * time.Duration(timeOut))
}

func (this *SessionInfo) InitSession(timeOut int) {
  this.SessionInvalidateTime = time.Now().Add(time.Second * time.Duration(timeOut))
  this.IsSerialized = false
}

func (this *Session) SaveSessionValue(valueId string, obj interface{}){
  this.SessionValue[valueId] = obj
}

func (this *Session) GetSessionValue(valueId string)  (interface{},bool) {
    val,ok := this.SessionValue[valueId]
    return val,ok 
}

