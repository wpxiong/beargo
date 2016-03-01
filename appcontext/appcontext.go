package appcontext

import (
  "github.com/beargo/log"
)

func init() {
  log.InitLog()
}

type AppContext struct {
  Config string
  Port int
}
