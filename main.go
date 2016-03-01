package main

import (
  "./log"
  "./webapp"
  "./appcontext"
  "runtime"
)

func main() {
   log.InitLog()
   runtime.GOMAXPROCS(runtime.NumCPU())
   var appCon appcontext.AppContext = appcontext.AppContext{ Config : "" , Port : 9001 }
   app := webapp.New(&appCon)
   indexCtrl := &IndexControl{}
   app.AddRoute("/test/<pam:[0-9]+>",indexCtrl,"Index")
   app.Start()
}