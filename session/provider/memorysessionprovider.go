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

func (provider *MemorySessionProvider)  LoadSession(sessionId string) (session.Session, error) {
   return session.Session{},nil
}

func (provider *MemorySessionProvider)  DeleteSession(sessionId string) error {
   return nil
}

func (provider *MemorySessionProvider)  RemoveAllSession(timeout int64)  error {
  return nil
}
