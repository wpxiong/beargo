package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
)

func init() {
  log.InitLog()
}

func DBtransactionEndinterceptor(app *appcontext.AppContext) bool {
   log.Debug("DBtransactioninterceptor Start")
   return true
}