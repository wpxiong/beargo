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

const (
   SESSION_SCAN_TIME = 5 //second
)

type sessionManager struct {
  Sessionprovider   SessionProvider
  CookieName  string
  SessionLifeTime  int64 //second
  SessionProviderMap  map[string] SessionProvider
  sessionAccess  *sync.Mutex
}

var sessionmanager *sessionManager = nil

func initSessionProvider (session_manager *sessionManager,provider SessionProvider) {
  session_manager.Sessionprovider = provider
  provider.InitProvider(session_manager.SessionLifeTime)
}

func GetSessionManager() *sessionManager{
   return sessionmanager
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

func StartSessionManager() {
   sessionmanager.Sessionprovider.DeserializeSession()
   go sessionmanager.startSessionManagerListener()
}

func (this* sessionManager ) startSessionManagerListener() {
   for true {
      this.Sessionprovider.ClearSession(this.sessionAccess)
      time.Sleep(time.Second * time.Duration(SESSION_SCAN_TIME))
   }
}


func StopSessionManager() {
   sessionmanager.Sessionprovider.SerializeSession()
}


func (this *sessionManager ) getSession(sessionId string) Session {
   var res bool = true
   var sess Session
   var err error
   this.sessionAccess.Lock()
   res = this.Sessionprovider.FindSessionById(sessionId)
   this.sessionAccess.Unlock()
   if !res {
     sess = this.CreateNewSession()
   }else {
     sess,err = this.Sessionprovider.LoadSessionById(sessionId)
     if err != nil {
        log.Error("Load Session Error")
        sess = this.CreateNewSession()
     }
   }
   return sess
}



func  NewSession(r *http.Request , w http.ResponseWriter)  Session{
  var sess Session
  cookie, err := r.Cookie(sessionmanager.CookieName)
  if err != nil || cookie.Value == "" {
      sess = sessionmanager.newSession()
      cookie := http.Cookie{Name: sessionmanager.CookieName, Value: url.QueryEscape(sess.SessionId), Path: "/", HttpOnly: true, MaxAge: int(sessionmanager.SessionLifeTime)}
      http.SetCookie(w, &cookie)
  } else {
      sid, _ := url.QueryUnescape(cookie.Value)
      sess = sessionmanager.getSession(sid)
      //update cookie
      if sess.SessionId != sid {
          cookie.Path = "/"
          cookie.MaxAge = int(sessionmanager.SessionLifeTime)
          cookie.HttpOnly = true
          cookie.Value = url.QueryEscape(sess.SessionId)
          http.SetCookie(w, cookie)
      }
  }
  return sess
}


func (this *sessionManager ) newSession() Session {
   sess := this.CreateNewSession()
   return sess
}

func (this *sessionManager ) CreateNewSession() Session {
   var sessionIdString string
   var res bool = true
   this.sessionAccess.Lock()
   for res  {
     sessionIdString = generateId()
     res = this.Sessionprovider.FindSessionById(sessionIdString)
   }
   sess,err := this.Sessionprovider.CreateSession(sessionIdString)
   this.sessionAccess.Unlock()
   if err != nil {
     log.Error("Create Session Error")
     panic("Create Session Error")
   }
   return sess
}


func GetInvalidateTime(timeOut int64 ) time.Time {
   return time.Now().Add(time.Second * time.Duration(timeOut))
}