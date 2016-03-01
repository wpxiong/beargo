package webhttp

import (
  "github.com/beargo/log"
  "net/http"
)

func init() {
  log.InitLog()
}

type HttpResponse struct {
  HttpResponseWriter *http.ResponseWriter
}


