package main

import (
 "./testroute"
 "../log"
)


func init(){
  log.InitLog()
}


func main() {
  testroute.Test()
}