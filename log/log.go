package log

import (
  "fmt"
  "time"
)

func init() {
  
}


type LogLevel int

const (
    DebugLevel   LogLevel = iota
    InfoLevel 
    TraceLevel
    ErrorLevel
)

var log *logInfo
type logInfo struct {
   loglevel LogLevel
<<<<<<< HEAD
   timeFormate string
=======
>>>>>>> origin/master
}

func (this *logInfo) printArray(level LogLevel,message ... interface{}) {
   if this.loglevel <= level {
<<<<<<< HEAD
      fmt.Print(time.Now().Format(this.timeFormate) + " ")
=======
>>>>>>> origin/master
      fmt.Println(message)
   }
}


func (this *logInfo) printObject(level LogLevel,message interface{}) {
   if this.loglevel <= level {
<<<<<<< HEAD
     fmt.Print(time.Now().Format(this.timeFormate) + " ")
=======
>>>>>>> origin/master
     fmt.Println(message)
   }
}

func (this *logInfo) printArrayNoReturn(level LogLevel,message ... interface{} ) {
  if this.loglevel <= level {
<<<<<<< HEAD
    fmt.Print(time.Now().Format(this.timeFormate) + " ")
=======
>>>>>>> origin/master
    fmt.Print(message)
  }
}


func (this *logInfo) printObjectReturn(level LogLevel,message interface{}) {
  if this.loglevel <= level {
<<<<<<< HEAD
     fmt.Print(time.Now().Format(this.timeFormate) + " ")
=======
>>>>>>> origin/master
     fmt.Print(message)
  }
}

func InitLog() {
   if log == nil {
<<<<<<< HEAD
     log = &logInfo{loglevel : InfoLevel,timeFormate: "2006-01-02 15:04:05.999Z07" }
=======
     log = &logInfo{loglevel : InfoLevel}
>>>>>>> origin/master
   }
}

func InitLogWithLevel(level string) {
   var logLevel LogLevel = InfoLevel
   if level == "Debug" {
      logLevel = DebugLevel
   }else if level == "Error"  {
      logLevel = ErrorLevel
   }else if level == "Trace"  {
      logLevel = TraceLevel
   }
   if log == nil {
<<<<<<< HEAD
      log = &logInfo{loglevel : logLevel,timeFormate: "2006-01-02 15:04:05.999Z07" }
=======
      log = &logInfo{loglevel : logLevel}
>>>>>>> origin/master
   }else {
      log.loglevel = logLevel
   }
}

func Info(message interface{}) {
   switch message.(type) {
     case [] interface{}:
       log.printArray(InfoLevel,message)
     default:
       log.printObject(InfoLevel,message)
   }
}

func InfoNoReturn(message interface{}) {
   switch message.(type) {
     case [] interface{}:
       log.printArrayNoReturn(InfoLevel,message)
     default:
       log.printObjectReturn(InfoLevel,message)
   }
}

func Debug(message interface{}) {
   switch message.(type) {
     case [] interface{}:
       log.printArray(DebugLevel,message)
     default:
       log.printObject(DebugLevel,message)
   }
}

func ErrorArray(message ... interface{}) {
   log.printArray(ErrorLevel,message)
}

func Error(message interface{}) {
   switch message.(type) {
     case [] interface{}:
       log.printArray(ErrorLevel,message)
     default:
       log.printObject(ErrorLevel,message)
   }
}

func DebugNoReturn(message interface{}) {
   switch message.(type) {
     case [] interface{}:
       log.printArrayNoReturn(DebugLevel,message)
     default:
       log.printObjectReturn(DebugLevel,message)
   }
}




