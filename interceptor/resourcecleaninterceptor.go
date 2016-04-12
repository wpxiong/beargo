package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
)

func init() {
  log.InitLog()
}


func dbRollBack(app *appcontext.AppContext) {
   log.Debug("DBtransaction RollBack")
   if  app.Trans != nil {
       for key,trans := range app.Trans {
          err := trans.Rollback()
          if err != nil {
            log.ErrorArray(key,err)
          }
       }
   }
   app.Trans = nil
}

func ResourceCleaninterceptor(app *appcontext.AppContext) bool {
   log.Debug("ResourceCleaninterceptor Start")
   dbRollBack(app)
   app.DestoryAppContext()
   return true
}
