package session

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/constvalue"
  "sync"
  "strconv"
)

func init() {
  log.InitLog()
}



type sessionManager struct {
  Sessionprovider   *SessionProvider
  CookieName  string
  SessionLifeTime  int64 //second
  SessionProviderMap  map[string] SessionProvider
  sessionAccess  *sync.Mutex
}

var sessionmanager *sessionManager = nil

func initSessionProvider (session_manager *sessionManager,provider string) {


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
     sessionmanager = &sessionManager{CookieName:constvalue.SESSION_NAME,SessionLifeTime:session_timeout,sessionAccess:&sync.Mutex{},SessionProviderMap:sessionProviderMap}
     initSessionProvider(sessionmanager,provider)
  }
}