package filter

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/binder"
)

func init() {
  log.InitLog()
}

func ParameterBinderFilter(app *appcontext.AppContext) bool {
   log.Debug("ParameterBinderFilter Start")
   binder.BinderParameter(app)
   return true
}