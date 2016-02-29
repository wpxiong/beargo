package testroute

import (
  "../../log"
  "../../route"
  "../../appcontext"
)

func Test() {
   log.InitLog()
   app := &appcontext.AppContext{}
   rt := route.NewRouteProcess(app)
   rt.Add("/xiong/<wen:[0-9]+>/ping")
   rt.Add("/rrrr/ping/")
   rt.Add("/rrrr/ggg/")
   rt.Add("/rrrr/sss/")
   rt.DebugInfo()
}