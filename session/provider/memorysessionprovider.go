package provider

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/session"
)

func init() {
  log.InitLog()
}


type MemorySessionProvider struct {

}

func (provider *MemorySessionProvider) CreateSession(sessionId string) (session.Session, error) {
   return session.Session{},nil
}


func (provider *MemorySessionProvider)  DeleteSession(sessionId string) error {
   return nil
}

func (provider *MemorySessionProvider)  RemoveAllSession(sessionIdList []string) error {
  return nil
}

func (provider *MemorySessionProvider)  FindSessionById(sessionId string) bool {
   return true
}

func (provider *MemorySessionProvider)  StoreSessionById(sessionId string) error {
   return nil
}

func (provider *MemorySessionProvider)  StoreSession(sessionId []string)  error{
   return nil
}

func (provider *MemorySessionProvider)  LoadSessionById(sessionId string) (session.Session,error){
   return session.Session{},nil
}

func (provider *MemorySessionProvider)  LoadSession(sessionId []string) (map[string]session.Session ,error) {
   return make(map[string]session.Session),nil
}
