package main

import (
  "github.com/beargo/log"
  "github.com/beargo/controller"
  "github.com/beargo/route"
  "fmt"
  "net/http"
)

func init() {
  log.InitLog()
}

type IndexControl struct {
  controller.Controller
}

func (*IndexControl) Index(rti *route.RouteInfo){
  var w http.ResponseWriter = *rti.Writer.HttpResponseWriter
  fmt.Fprintf(w,"Edit")
}
