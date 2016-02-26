package log

import (
  "fmt"
  "strings"
)

func init() {
  
}

var log *logInfo
type logInfo struct {

}

func (this *logInfo) info(message string) {
   //fmt.Println(message)
}


func InitLog() {
   if log == nil {
     log = &logInfo{}
     fmt.Println("start log")
   }
}

func Info(message... string) {
   mess := strings.Join(message," ")
   log.info(mess)
}




