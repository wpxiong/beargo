package log

import (
  "fmt"
)

func init() {
  
}

var log *logInfo
type logInfo struct {

}

func (this *logInfo) infoArray(message ... interface{}) {
   fmt.Println(message)
}


func (this *logInfo) info(message interface{}) {
   fmt.Println(message)
}

func (this *logInfo) infoArrayNoReturn(message ... interface{}) {
   fmt.Print(message)
}


func (this *logInfo) infoNoReturn(message interface{}) {
   fmt.Print(message)
}

func InitLog() {
   if log == nil {
     log = &logInfo{}
     fmt.Println("start log")
   }
}

func Info(message interface{}) {
   switch message.(type) {
     case [] interface{}:
       log.infoArray(message)
     default:
       log.info(message)
   }
}

func InfoNoReturn(message interface{}) {
   switch message.(type) {
     case [] interface{}:
       log.infoArrayNoReturn(message)
     default:
       log.infoNoReturn(message)
   }
}





