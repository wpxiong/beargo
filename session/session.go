package session

import (
  "github.com/wpxiong/beargo/log"
  "time"
  "reflect"
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
  SessionSerializeInfo   map[string] []byte
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

func (this *Session) DeleteSessionValue(valueId string){
  delete(this.SessionValue,valueId)
}

func (this *Session) GetSessionValue(valueId string,obj interface{})  bool {
    val,ok := this.SessionValue[valueId]
    if ok {
      reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(val))
      return true
    }
    return GetSessionManager().LoadSessionBySessionId(valueId,this.SessionSerializeInfo[valueId],obj,this)
}

