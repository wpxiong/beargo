package memorycash

import (
  "github.com/wpxiong/beargo/log"
  "errors"
)

func init() {
  log.InitLog()
}

var memorycashManager *MemoryCashManager 

const (
  DEFAULT_CASH_SIZE = 1000
)

type MemoryCash interface {
   AddStringObject(key string, val string) error
   GetStringObject(key string) (string,error)
   GetObject(key string) (interface{},error)
   AddObject(key string, val interface{}) error
   DeleteObject(key string) error
   InitMemoryCash(size int) error
}


type MemoryCashManager struct {
   memoryCash MemoryCash
}

func CreateMemoryCashManager(memorycashType string){
  if memorycashManager == nil {
     memorycashManager = &MemoryCashManager{}
     initMemoryCashManager(memorycashType)
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

func initMemoryCashManager(memorycashType string){
   switch memorycashType {
      case "mapmemorycash":
        memorycashManager.memoryCash = &MapMemoryCash{}
        if err := memorycashManager.memoryCash.InitMemoryCash(DEFAULT_CASH_SIZE);err != nil {
            memorycashManager.memoryCash = nil
            log.Error(err)
        }
      default:
        memorycashManager.memoryCash = nil
   }
}
