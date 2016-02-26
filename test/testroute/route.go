package testroute

import (
  "../../log"
  "../../route"
  "../../appcontext"
)

func Test() {
   log.InitLog()
   app := &appcontext.AppContext{}
   route.NewRouteProcess(app)
   
}