package provider

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/session"
  "errors"
)

func init() {
  log.InitLog()
}

const (
  MAX_SESSION_COUNT = 1000
)


type MemorySessionProvider struct {
   Sessionmap map[string] *session.Session
}

func (provider *MemorySessionProvider) InitProvider(){
   provider.Sessionmap = make(map[string]*session.Session)
}

func (provider *MemorySessionProvider) CreateSession(id string) (session.Session, error) {
   ses := &session.Session{}
   ses.InitSession(id)
   provider.Sessionmap[id] = ses
   return *ses,nil
}


func (provider *MemorySessionProvider)  DeleteSession(sessionId string) error {
   delete(provider.Sessionmap,sessionId)
   return nil
}

func (provider *MemorySessionProvider)  FindSessionById(sessionId string) bool {
   if _, ok := provider.Sessionmap[sessionId]; ok{
     return  true
   }else {
     return false
   }
}

func (provider *MemorySessionProvider)  LoadSessionById(sessionId string) (session.Session,error){
  if val, ok := provider.Sessionmap[sessionId]; ok{
    return  *val,nil
  }else {
    return session.Session{} ,errors.New("Load Session Error")
  }
}

