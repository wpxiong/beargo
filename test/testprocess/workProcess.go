package main

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/process"
  "runtime"
  "time"
)

func processJob(interface{}) interface{} {
   time.Sleep(20)
   return 5
}

func AddJob(){
   for i:=1; i< 1000000; i++ {
      workjob := &process.WorkJob{}
      workjob.WorkProcess = processJob
      process.AddJob(workjob)
   }
}

func TestMoudle() {
   log.InitLog()
   runtime.GOMAXPROCS(runtime.NumCPU())
   pro := process.New()
   pro.Init_Default()
   AddJob()
   process.StopWork()
}

func main(){
  TestMoudle()
}
