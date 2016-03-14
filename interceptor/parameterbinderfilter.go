package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/binder"
)

func init() {
  log.InitLog()
}

func ParameterBinderinterceptor(app *appcontext.AppContext) bool {
   log.Debug("ParameterBinderinterceptor Start")
   binder.BinderParameter(app)
   return true
}