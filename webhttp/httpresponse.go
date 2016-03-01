package webhttp

import (
  "../log"
  "net/http"
)

func init() {
  log.InitLog()
}

type HttpResponse struct {
  HttpResponseWriter *http.ResponseWriter
}


