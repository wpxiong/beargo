package appcontext

import (
  "../log"
)

func init() {
  log.InitLog()
}

type AppContext struct {
   config string
}
