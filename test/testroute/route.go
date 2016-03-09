package main

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/route"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/webhttp"
  "github.com/wpxiong/beargo/controller"
)

type IndexControl struct {
  controller.Controller
}

func (*IndexControl) Index(ctx *appcontext.AppContext){
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
   rt.Add("/xiong/<id:int>",ctr,"Index")
   rt.Add("/rrrr/sss/xxxx",ctr,"Index")
   //rt.DebugInfo()
   request := webhttp.HttpRequest{Urlpath :"/xiong/wen997/ping?te=ag&rr=345" }
   request2 := webhttp.HttpRequest{Urlpath :"/rrrr/447?te=ag&rr=345" }
   request3 := webhttp.HttpRequest{Urlpath :"/rrrr/sss/xxxx?te=ag&rr=345" }
   request4 := webhttp.HttpRequest{Urlpath :"/xiong/445?te=ag&rr=345" }
   rt.ProcessRequest(&request)
   rt.ProcessRequest(&request2)
   //rti.DebugInfo()
   rt.ProcessRequest(&request3)
   rt.ProcessRequest(&request4)
   //rti.DebugInfo()
}


func init(){
  log.InitLogWithLevel("Debug")
}


func main() {
  log.InitLogWithLevel("Debug")
  Test()
}