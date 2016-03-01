package testroute

import (
  "../../log"
  "../../route"
  "../../appcontext"
  "../../http"
  "../../controller"
)

type IndexControl struct {
  controller.Controller
}

func (*IndexControl) Index(){
  log.Info("Index")
}


func Test() {
   log.InitLog()
   app := &appcontext.AppContext{}
   rt := route.NewRouteProcess(app)
   ctr := &IndexControl{}
   rt.Add("/xiong/wen<pam:[0-9]+>/ping",ctr,"Index")
   rt.Add("/rrrr/ping/mmmm",ctr,"Index")
   rt.Add("/rrrr/ggg",ctr,"Index")
   rt.Add("/rrrr/<id:int>",ctr,"Index")
   rt.Add("/rrrr/sss/xxxx",ctr,"Index")
   //rt.DebugInfo()
   request := http.HttpRequest{Urlpath :"/xiong/wen997/ping?te=ag&rr=345" }
   request2 := http.HttpRequest{Urlpath :"/rrrr/447?te=ag&rr=345" }
   request3 := http.HttpRequest{Urlpath :"/rrrr/sss/xxxx?te=ag&rr=345" }
   var rti *route.RouteInfo
   rti = rt.ProcessRequest(&request)
   rti.CallMethod()
   rt.ProcessRequest(&request2)
   //rti.DebugInfo()
   rt.ProcessRequest(&request3)
   //rti.DebugInfo()
}