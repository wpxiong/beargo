package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/util/httprequestutil"
)

func init() {
  log.InitLog()
}


func ParameterParseinterceptor(app *appcontext.AppContext) bool {
   log.Debug("ParameterParseinterceptor Start")
   httprequestutil.ProcessHttpRequestParam(app)
   return true
}