package process

import (
  "github.com/wpxiong/beargo/log"
  "time"
  "sync"
)

func init() {
  log.InitLog()
}

type ProcessStatus int

const (
    Create   ProcessStatus = iota
    Start 
    Stop
    Receive
    Finish
    TimeOut
)

const (
  DefaultMaxProcessCount = 10
  DefaultWorkListMaxCount = 2000
  DefaultTimeOut = 1000
)

type WorkJob struct {
   workReceiveTime time.Time
   workStartTime time.Time
   workEndTime time.Time
   workProcessTime time.Duration
   workDurationTime time.Duration
   Parameter interface{}
   Result interface{}
   WorkProcess func (interface{}) interface{}
   mutex *sync.Mutex
   status ProcessStatus
}


type WorkProcess struct {
   maxProcessCount int
   status ProcessStatus
   workListMaxCount int
   processTimeOut int
   processReceiveChan chan *WorkJob
   processFinishChan chan int
   stopChan chan int
   mutex *sync.Mutex
}

var workprocess *WorkProcess

func startReceiveWorkProcess(this *WorkProcess) {
   log.Debug("start workProcess")
   var job *WorkJob
   this.mutex.Lock()
   this.status = Start
   this.mutex.Unlock()
   for true {
     select {
       case job = <- this.processReceiveChan :
          log.Debug("receive new workJob")
          duration := time.Now().Sub(job.workReceiveTime)
          if (duration.Seconds() / 1000 < float64(this.processTimeOut)) {
             this.processWork(job)
             job.mutex.Lock()
             job.status = Finish
             job.mutex.Unlock()
          }else{
             job.mutex.Lock()
             job.status = TimeOut
             job.mutex.Unlock()
          }
       case _ = <- this.processFinishChan :
          log.Debug("stop workProcess")
          this.mutex.Lock()
          this.status = Stop
          this.mutex.Unlock()
          this.stopChan <- 1
          return
     }
   }
}


func (this *WorkProcess) Init_Default () {
   log.Debug("start WorkJob")
   this.Init(DefaultMaxProcessCount,DefaultWorkListMaxCount,DefaultTimeOut)
}

func (this *WorkProcess) Init (maxProcessCount int , workListMaxCount int , processTimeOut int ) {
   this.processTimeOut = processTimeOut
   this.maxProcessCount = maxProcessCount
   this.workListMaxCount = workListMaxCount
   this.processReceiveChan = make(chan *WorkJob , workListMaxCount)
   this.processFinishChan = make(chan int,maxProcessCount)
   this.stopChan = make(chan int,maxProcessCount)
   this.status = Create
   this.mutex = &sync.Mutex{}
   for i:=0; i < maxProcessCount ; i++  {
      go startReceiveWorkProcess(this)
   }
}


func (this *WorkProcess) stopProcessWork () {
   log.Debug("stop WorkJob")
   for i:=0; i < this.maxProcessCount ; i++  {
      this.processFinishChan <- 1
   }
   stopProcess := 0
   for stopProcess < this.maxProcessCount {
        _ = <- this.stopChan
       stopProcess ++ 
   }
   this.mutex.Lock()
   this.status = Finish
   this.mutex.Unlock()
   log.Debug("finish WorkJob")
}


func (this *WorkProcess) processWork(workJob *WorkJob){
   workJob.mutex.Lock()
   workJob.workStartTime = time.Now()
   workJob.mutex.Unlock()
   log.Debug("work process start")
   if workJob.WorkProcess != nil {
      result := workJob.WorkProcess(workJob.Parameter)
      workJob.mutex.Lock()
      workJob.Result = result
      workJob.mutex.Unlock()
   }
   log.Debug("workProcess end")
   workJob.mutex.Lock()
   workJob.workEndTime = time.Now()
   workJob.workProcessTime = workJob.workEndTime.Sub(workJob.workStartTime)
   workJob.workDurationTime = workJob.workEndTime.Sub(workJob.workReceiveTime)
   workJob.mutex.Unlock()
}

func (this *WorkProcess) addJob(workJob *WorkJob) {
  if workJob != nil {
     if workJob.mutex == nil {
        workJob.mutex = &sync.Mutex{}
     }
     workJob.status = Receive
     workJob.workReceiveTime = time.Now()
     this.processReceiveChan <- workJob
  }
}


func New() *WorkProcess {
   if workprocess == nil {
      workprocess = &WorkProcess{}
   }
   return workprocess;
}



func StopWork(){
  if workprocess != nil {
     workprocess.stopProcessWork()
  }
}

func AddJob(workJob *WorkJob){
  if workprocess != nil {
     workprocess.addJob(workJob)
  }
}
