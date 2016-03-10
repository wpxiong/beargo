package memorycash

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/constvalue"
  "github.com/wpxiong/beargo/appcontext"
  "errors"
  "strconv"
)

func init() {
  log.InitLog()
}

var memorycashManager *MemoryCashManager 


type MemoryCash interface {
   AddStringObject(key string, val string) error
   GetStringObject(key string) (string,error)
   GetObject(key string) (interface{},error)
   AddObject(key string, val interface{}) error
   DeleteObject(key string) error
   InitMemoryCash(size int64) error
}


type MemoryCashManager struct {
   memoryCash MemoryCash
}

func CreateMemoryCashManager(context *appcontext.AppContext) {
  if memorycashManager == nil {
     var memorycashType string
     memorycashType = context.GetConfigValue(constvalue.CASH_TYPE_KEY,constvalue.DEFAULT_CASH_TYPE).(string)
     maxcashsize := context.GetConfigValue(constvalue.CASH_MAXSIZE_KEY,"").(string)
     var maxsize int64 = constvalue.DEFAULT_CASH_MAXSIZE
     if maxcashsize != ""{
        size,err :=  strconv.Atoi(maxcashsize)
        if err != nil {
           maxsize = int64(size)
        }
     }
     memorycashManager = &MemoryCashManager{}
     initMemoryCashManager(memorycashType,context,maxsize)
  }
}


func AddStringObject(key string, val string) error {
   if memorycashManager == nil || memorycashManager.memoryCash == nil {
       return errors.New("memory cash is not be created")
   }else {
       return memorycashManager.memoryCash.AddStringObject(key,val)
   }
}

func GetStringObject(key string) (string,error) {
   if memorycashManager == nil || memorycashManager.memoryCash == nil {
       return "",errors.New("memory cash is not be created")
   }else {
       return memorycashManager.memoryCash.GetStringObject(key)
   }
}

func GetObject(key string) (interface{},error) {
   if memorycashManager == nil || memorycashManager.memoryCash == nil {
       return nil,errors.New("memory cash is not be created")
   }else {
       return memorycashManager.memoryCash.GetStringObject(key)
   }
}

func AddObject(key string, val interface{}) error {
   if memorycashManager == nil || memorycashManager.memoryCash == nil {
       return errors.New("memory cash is not be created")
   }else {
       return memorycashManager.memoryCash.AddObject(key,val)
   }
}

func DeleteObject(key string) error {
   if memorycashManager == nil || memorycashManager.memoryCash == nil {
       return errors.New("memory cash is not be created")
   }else {
       return memorycashManager.memoryCash.DeleteObject(key)
   }
}

func initMemoryCashManager(memorycashType string,context *appcontext.AppContext,cashSize int64){
   switch memorycashType {
      case "mapmemorycash":
        memorycashManager.memoryCash = &MapMemoryCash{}
        if err := memorycashManager.memoryCash.InitMemoryCash(cashSize);err != nil {
            memorycashManager.memoryCash = nil
            log.Error(err)
        }
      default:
        memorycashManager.memoryCash = nil
   }
}
