package binder

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
)

func init() {
  log.InitLog()
}


func BinderParameter(appcon *appcontext.AppContext){
  log.Debug("start")
}