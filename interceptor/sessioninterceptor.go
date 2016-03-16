package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/session"
)

func init() {
  log.InitLog()
}

func Sessioninterceptor(app *appcontext.AppContext) bool {
   log.Debug("Sessioninterceptor Start")
   session.NewSession(app.Request.HttpRequest , *app.Writer.HttpResponseWriter)
   return true
}