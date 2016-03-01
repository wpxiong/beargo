package webhttp

import (
  "github.com/wpxiong/beargo/log"
  "net/http"
)

func init() {
  log.InitLog()
}

type HttpResponse struct {
  HttpResponseWriter *http.ResponseWriter
}


