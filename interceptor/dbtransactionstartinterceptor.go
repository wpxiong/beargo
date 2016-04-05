package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
)

func init() {
  log.InitLog()
}

func DBtransactionStartinterceptor(app *appcontext.AppContext) bool {
   log.Debug("DBtransactioninterceptor Start")
   return true
}