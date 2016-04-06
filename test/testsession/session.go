package main

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/session"
  "github.com/wpxiong/beargo/session/provider"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/constvalue"
  "time"
)


func TestSession() {
   log.InitLogWithLevel("Debug")
   config := appcontext.AppConfigContext{Port :9001,ConfigPath : "./setting.conf"}
   var appCon appcontext.AppContext = appcontext.AppContext{ ConfigContext :  &config}
   sessionProviderMap := make(map[string]session.SessionProvider)
   sessionProviderMap[constvalue.DEFAULT_SESSION_PROVIDER] = &provider.MemorySessionProvider{}
   session.CreateSessionManager(&appCon,sessionProviderMap)
   sessmanager := session.GetSessionManager()
   sess := sessmanager.CreateNewSession()
   sess.SaveSessionValue("test","xiongwenping")
   sess = sessmanager.CreateNewSession()
   session.StartSessionManager()
   sess.SaveSessionValue("test2",23)
   time.Sleep(time.Second * time.Duration(30))
   session.StopSessionManager()
}

func main(){
  TestSession()
}
