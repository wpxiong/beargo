package session

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/constvalue"
  "sync"
  "time"
  "strconv"
  "crypto/rand"
  "encoding/hex"
  "net/http"
  "net/url"
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
  provider.InitProvider()
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
  buffer := make([]byte, constvalue.DEFAULT_SESSIONID_SIZE)
  if _, err := rand.Read(buffer); err != nil {
    panic(err)
  }
  text := hex.EncodeToString(buffer)
  return text
}


func (this *sessionManager ) getSession(sessionId string) Session {
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
        sess = this.createNewSession()
     }
   }
   sessionInfo := this.sessionInfoMap[sess.sessionId]
   sessionInfo.sessionInvalidateTime = getInvalidateTime(this.SessionLifeTime)
   return sess
}

func getInvalidateTime(timeOut int64 ) time.Time {
   return time.Now().Add(time.Second * time.Duration(timeOut))
}

func  NewSession(r *http.Request , w http.ResponseWriter)  Session{
  var sess Session
  cookie, err := r.Cookie(sessionmanager.CookieName)
  if err != nil || cookie.Value == "" {
      sess = sessionmanager.newSession()
      cookie := http.Cookie{Name: sessionmanager.CookieName, Value: url.QueryEscape(sess.sessionId), Path: "/", HttpOnly: true, MaxAge: int(sessionmanager.SessionLifeTime)}
      http.SetCookie(w, &cookie)
  } else {
      sid, _ := url.QueryUnescape(cookie.Value)
      sess = sessionmanager.getSession(sid)
      if sess.sessionId != sid {
          cookie.Path = "/"
          cookie.MaxAge = int(sessionmanager.SessionLifeTime)
          cookie.HttpOnly = true
          cookie.Value = url.QueryEscape(sess.sessionId)
          http.SetCookie(w, cookie)
      }
  }
  return sess
}


func (this *sessionManager ) newSession() Session {
   this.sessionAccess.Lock()
   defer this.sessionAccess.Unlock()
   sess := this.createNewSession()
   return sess
}

func (this *sessionManager ) createNewSession() Session {
   log.Debug("createNewSession Start")
   var sessionIdString string
   var res bool = true
   for res  {
     sessionIdString = generateId()
     res = this.Sessionprovider.FindSessionById(sessionIdString)
   }
   sess,err := this.Sessionprovider.CreateSession(sessionIdString)
   if err != nil {
     log.Error("Create Session Error")
     sess = sessionmanager.newSession()
   }
   sessiontimeout := getInvalidateTime(this.SessionLifeTime)
   this.sessionInfoMap[sessionIdString] = SessionInfo{sessionId:sessionIdString,isInMemory:true,sessionInvalidateTime:sessiontimeout}
   return sess
}