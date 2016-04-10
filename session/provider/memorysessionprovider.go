package provider

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/session"
  "errors"
  "sync"
  "time"
  "os"
  "io"
  "bufio"
  "encoding/gob"
)

var SESSION_FILE string
var SESSION_INFO_FILE string
var SESSION_FOLDER string
func init() {
  log.InitLog()
  SESSION_FILE = string(os.PathSeparator) + "session" + string(os.PathSeparator) + "session.data"
  SESSION_INFO_FILE = string(os.PathSeparator) + "session" + string(os.PathSeparator) + "sessioninfo.data"
  SESSION_FOLDER = string(os.PathSeparator) + "session"
}

const (
  MAX_SESSION_COUNT = 1000
)


type MemorySessionProvider struct {
   Sessionmap map[string] *session.Session
   SessionInfoMap map[string] *session.SessionInfo
   SessionLifeTime int64
}

func (provider *MemorySessionProvider) InitProvider(sessionLifeTime int64){
   provider.Sessionmap = make(map[string]*session.Session)
   provider.SessionInfoMap = make(map[string] *session.SessionInfo)
   provider.SessionLifeTime = sessionLifeTime
}

func (provider *MemorySessionProvider) ClearSession(sessionAccess *sync.Mutex){
   for key,value := range provider.SessionInfoMap {
      if value.SessionInvalidateTime.Before(time.Now()) {
         sessionAccess.Lock()
         delete(provider.SessionInfoMap,key)
         provider.DeleteSession(key)
         sessionAccess.Unlock()
      }
   }
}

func (provider *MemorySessionProvider) CreateSession(id string) (session.Session, error) {
   ses := &session.Session{}
   ses.InitSession(id)
   provider.Sessionmap[id] = ses
   sessiontimeout := session.GetInvalidateTime(provider.SessionLifeTime)
   provider.SessionInfoMap[ses.SessionId] = &session.SessionInfo{SessionId:id,IsSerialized:false,SessionInvalidateTime:sessiontimeout}
   return *ses,nil
}


func (provider *MemorySessionProvider)  DeleteSession(sessionId string) error {
   delete(provider.Sessionmap,sessionId)
   delete(provider.SessionInfoMap,sessionId)
   return nil
}

func (provider *MemorySessionProvider)  FindSessionById(sessionId string) bool {
   if _, ok := provider.Sessionmap[sessionId]; ok{
     return  true
   }else {
     return false
   }
}

func (provider *MemorySessionProvider)  LoadSessionById(sessionId string) (session.Session,error){
  if val, ok := provider.Sessionmap[sessionId]; ok{
    sessionInfo := provider.SessionInfoMap[val.SessionId]
    sessionInfo.SessionInvalidateTime = session.GetInvalidateTime(provider.SessionLifeTime)
    return  *val,nil
  }else {
    return session.Session{} ,errors.New("Load Session Error")
  }
}

func (provider *MemorySessionProvider) SerializeSession(sessionTypeInfo  map[string] interface{} ){
   for _,reflectValue := range sessionTypeInfo {
     gob.Register(reflectValue)
   }
   pwd, _ := os.Getwd()
   if err := os.MkdirAll(pwd + SESSION_FOLDER, 0777); err != nil {
       log.Error(err)
       return
   }
   var f *os.File
   var err error
   
   f, err = os.Create(pwd + SESSION_FILE)
   if err != nil {
      log.Error(err)
      return 
   }
   w := bufio.NewWriter(f)
   for key,val := range provider.Sessionmap {
     _,errinfo := w.WriteString(session.SerializeSession(*val)+"\n")
     if errinfo != nil {
        log.ErrorArray("serialize session [" + key + "] failture",errinfo)
     }
   }
   w.Flush()
   
   f, err = os.Create(pwd + SESSION_INFO_FILE)
   if err != nil {
      log.Error(err)
      return 
   }
   w = bufio.NewWriter(f)
   for key,val := range provider.SessionInfoMap {
     _,errinfo := w.WriteString(session.SerializeSessionInfo(*val)+"\n")
     if errinfo != nil {
        log.ErrorArray("serialize sessioninfo [" + key + "] failture",errinfo)
     }
   }
   w.Flush()
}

func (provider *MemorySessionProvider) DeserializeSession(sessionTypeInfo  map[string] interface{}){
  for _,reflectValue := range sessionTypeInfo {
     gob.Register(reflectValue)
  }
  pwd, _ := os.Getwd()
  file1, err := os.Open(pwd + SESSION_FILE)
  if err == nil {
      reader := bufio.NewReaderSize(file1, 4096)
      defer file1.Close()
      for {
        line, _, errinfo := reader.ReadLine()
        log.Debug(line)
        if errinfo == io.EOF {
            break
        } else if errinfo != nil {
           log.Error(errinfo)
        }else {
          var sess *session.Session = session.DeserializeSession(string(line))
          if sess != nil {
             provider.Sessionmap[sess.SessionId] = sess
          }
        }
      }
  }
  file2, err2 := os.Open(pwd + SESSION_INFO_FILE)
  if err2 == nil {
      defer file2.Close()
      reader := bufio.NewReaderSize(file2, 4096)
      for {
         line, _, errinfo := reader.ReadLine()
         if errinfo == io.EOF {
            break
         } else if errinfo != nil {
            log.Error(errinfo)
         }else {
            var sessinfo *session.SessionInfo = session.DeserializeSessionInfo(string(line))
            if sessinfo != nil {
               provider.SessionInfoMap[sessinfo.SessionId] = sessinfo
            }
         }
      }
  }
}

