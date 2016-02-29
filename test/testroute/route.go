package testroute

import (
  "../../log"
  "../../route"
  "../../appcontext"
  "../../http"
)

func Test() {
   log.InitLog()
   app := &appcontext.AppContext{}
   rt := route.NewRouteProcess(app)
   rt.Add("/xiong/<wen:[0-9]+>/ping")
   rt.Add("/rrrr/ping/")
   rt.Add("/rrrr/ggg/")
   rt.Add("/rrrr/sss/")
   //rt.DebugInfo()
   request := http.HttpRequest{Urlpath :"/xiong/wen997/ping?te=ag&rr=345" }
   rt.ProcessRequest(&request)
   
}