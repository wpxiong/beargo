package main

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/constvalue"
  "net/rpc"
  "os"
  "strconv"
)


func init() {
  log.InitLogWithLevel("Debug")
}


type Args struct {
    Command string
}


func main() {
  if len(os.Args) > 1 {
     command := os.Args[1]
     log.Debug(command)
     
     client, err := rpc.Dial("tcp",  constvalue.DEFAULT_MANAGER_HOST +  ":" + strconv.Itoa(constvalue.DEFAULT_MANAGER_PORT))
     if err != nil {
        log.ErrorArray("dialing:", err)
     }
     
     args := Args{Command:"stop"}
     var reply int
     err = client.Call("CommandInterface.SendCommand", args, &reply)
     if err != nil {
        log.ErrorArray("arith error:", err)
     }else {
        log.Info("result: " + strconv.Itoa(reply))
     }
  }
}