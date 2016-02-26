package main

import (
  "./log"
  "./process"
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

func main() {
   log.InitLog()
   runtime.GOMAXPROCS(runtime.NumCPU())
   pro := process.New()
   pro.Init_Default()
   AddJob()
   process.StopWork()
}