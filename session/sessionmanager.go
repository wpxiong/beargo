package session

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/constvalue"
  "sync"
  "time"
  "strconv"
)

func init() {
  log.InitLog()
}



type sessionManager struct {
  Sessionprovider   SessionProvider
  CookieName  string
  SessionLifeTime  int64 //second
  SessionProviderMap  map[string] SessionProvider
  sessionAccess  *sync.Mutex
  sessionInfoMap map[string]SessionInfo   
}

var sessionmanager *sessionManager = nil

func initSessionProvider (session_manager *sessionManager,provider SessionProvider) {
  session_manager.Sessionprovider = provider
}

func CreateSessionManager(context *appcontext.AppContext,sessionProviderMap map[string]SessionProvider ){
  if sessionmanager == nil {
     var timeoutStr string = context.GetConfigValue(constvalue.SESSION_TIMEOUT_KEY,"").(string)
     var session_timeout int64 = constvalue.DEFAULT_SESSION_TIMEOUT
     if timeoutStr != ""{
        timeout,err :=  strconv.Atoi(timeoutStr)
        if err != nil {
           session_timeout = int64(timeout)
        }
     }
     var provider string = context.GetConfigValue(constvalue.SESSION_PROVIDER_KEY,constvalue.DEFAULT_SESSION_PROVIDER).(string)
     sessionmanager = &sessionManager{sessionInfoMap:make(map[string]SessionInfo), CookieName:constvalue.SESSION_NAME,SessionLifeTime:session_timeout,sessionAccess:&sync.Mutex{},SessionProviderMap:sessionProviderMap}
     if sessionProviderMap[provider] != nil {
        initSessionProvider(sessionmanager,sessionProviderMap[provider])
     }else {
        initSessionProvider(sessionmanager,sessionProviderMap[provider])
     }
  }
}


func generateId() string {
  return ""
}


func (this *sessionManager ) GetSession(sessionId string) Session {
   this.sessionAccess.Lock()
   defer this.sessionAccess.Unlock()
   var res bool = true
   var sess Session
   var err error
   res = this.Sessionprovider.FindSessionById(sessionId)
   if !res {
     sess = this.createNewSession()
   }else {
     sess,err = this.Sessionprovider.LoadSessionById(sessionId)
     if err != nil {
        log.Error("Load Session Error")
        panic("Create NewSession Error")
     }
   }
   sessionInfo := this.sessionInfoMap[sess.sessionId]
   sessionInfo.sessionInvalidateTime = getInvalidateTime(this.SessionLifeTime)
   return sess
}

func getInvalidateTime(timeOut int64 ) time.Time {
   return time.Now().Add(time.Second * time.Duration(timeOut))
}


func (this *sessionManager ) NewSession() Session {
   this.sessionAccess.Lock()
   defer this.sessionAccess.Unlock()
   sess := this.createNewSession()
   return sess
}

func (this *sessionManager ) createNewSession() Session {
   var sessionIdString string
   var res bool = true
   for res  {
     sessionIdString = generateId()
     res = this.Sessionprovider.FindSessionById(sessionIdString)
   }
   sess,err := this.Sessionprovider.CreateSession(sessionIdString)
   if err != nil {
     log.Error("Create Session Error")
     panic("Create NewSession Error")
   }
   sessiontimeout := getInvalidateTime(this.SessionLifeTime)
   this.sessionInfoMap[sessionIdString] = SessionInfo{sessionId:sessionIdString,isInMemory:true,sessionInvalidateTime:sessiontimeout}
   return sess
}