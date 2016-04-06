package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/moudle"
)

func init() {
  log.InitLog()
}

func DBtransactionStartinterceptor(app *appcontext.AppContext) bool {
   log.Debug("DBtransactionStartinterceptor Start")
   if app.DBSession != nil {
      trans := make(map[string]* moudle.Trans,len(app.DBSession))
      for key,dbsession := range app.DBSession {
         trans[key] = dbsession.Begin()
      }
      app.Trans = trans
   }
   return true
}