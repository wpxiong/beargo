package filter

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/util/httprequestutil"
)

func init() {
  log.InitLog()
}


func ParameterParseFilter(app *appcontext.AppContext) bool {
   log.Debug("ParameterParseFilter Start")
   httprequestutil.ProcessHttpRequestParam(app)
   return true
}