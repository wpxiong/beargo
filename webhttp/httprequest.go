package webhttp

import (
  "../log"
  "net/http"
)

func init() {
  log.InitLog()
}

type HttpRequest struct {
  Urlpath string
  HttpRequest *http.Request
}


