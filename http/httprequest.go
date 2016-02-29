package http

import (
  "../log"
)

func init() {
  log.InitLog()
}

type HttpRequest struct {
  Urlpath string
}


