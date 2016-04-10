package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
)

func init() {
  log.InitLog()
}

func DBtransactionEndinterceptor(app *appcontext.AppContext) bool {
   log.Debug("DBtransactionEndinterceptor Start")
   if app.Trans != nil {
      errorlist := make([]error,len(app.Trans))
      var k int = 0
      for _,trans := range app.Trans {
         err := trans.Commit()
         errorlist[k] = err
         k+=1
      }
      result := true
      for _,errinfo := range errorlist {
         if errinfo != nil {
             result = false
         }
      }
      if !result {
         for key,trans := range app.Trans {
           err := trans.Rollback()
           if err != nil {
              log.ErrorArray(key,err)
           }
         }
      }
   }
   return true
}
