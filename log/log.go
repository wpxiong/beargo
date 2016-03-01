package log

import (
  "fmt"
)

func init() {
  
}

var log *logInfo
type logInfo struct {
   debug bool
}

func (this *logInfo) infoArray(message ... interface{}) {
   if this.debug {
      fmt.Println(message)
   }
}


func (this *logInfo) info(message interface{}) {
   if this.debug {
     fmt.Println(message)
   }
}

func (this *logInfo) infoArrayNoReturn(message ... interface{}) {
  if this.debug {
    fmt.Print(message)
  }
}


func (this *logInfo) infoNoReturn(message interface{}) {
  if this.debug {
     fmt.Print(message)
  }
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





