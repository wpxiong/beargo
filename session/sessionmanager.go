package session

import (
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}

type sessionManager struct {
  SessionProvider *sessionProvider
  
}