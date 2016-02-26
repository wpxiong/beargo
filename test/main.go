package main

import (
 "./testprocess"
 "./testroute"
 "../log"
)


func init(){
  log.InitLog()
}


func main() {
  testprocess.Test()
  testroute.Test()
}