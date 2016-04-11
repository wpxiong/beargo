package session

import (
  "github.com/wpxiong/beargo/log"
  "encoding/base64"
  "encoding/gob"
  "bytes"
  "strings"
)

func init() {
  log.InitLog()
}


func SerializeSessionInfo(sessionInfo SessionInfo)  string {
    b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(sessionInfo)
    if err != nil {
      log.Error(err)
      return ""
    }
    return base64.StdEncoding.EncodeToString(b.Bytes())
}

func SerializeSession(session Session) string {
    sessionid := session.SessionId
    session.SessionId = ""
    var resStr string = sessionid + " "
    for key,val := range session.SessionValue {
       b := bytes.Buffer{}
       e := gob.NewEncoder(&b)
       err := e.Encode(val)
       if err != nil {
          log.Error(err)
       }else {
          resStr += key + " " +  base64.StdEncoding.EncodeToString(b.Bytes()) + " "
       }
    }
    return resStr
}

func DeserializeSessionInfo(sessionInfoStr string) *SessionInfo {
    sessionInfo := SessionInfo{}
    byteArray, err := base64.StdEncoding.DecodeString(sessionInfoStr)
    if err != nil { 
       log.Error(err)
       return nil
    }
    b := bytes.Buffer{}
    b.Write(byteArray)
    d := gob.NewDecoder(&b)
    err = d.Decode(&sessionInfo)
    if err != nil {
       log.Error(err)
       return nil
    }
    sessionInfo.IsSerialized = true
    return &sessionInfo
}

func DeserializeSession(sessionStr string) *Session {
    session := Session{}
    sessioncomponent :=  strings.Split(sessionStr," ")
    if len(sessioncomponent) < 2 {
       log.Error("DeserializeSession Error")
       return nil
    }
    sessionid := sessioncomponent[0]
    session.SessionSerializeInfo = make(map[string] []byte)
    for i := 1; i<len(sessioncomponent)-1 ; i=i+2 {
       key := sessioncomponent[i]
       valstr := sessioncomponent[i+1]
       byteArray, err := base64.StdEncoding.DecodeString(valstr)
       if err != nil { 
         log.Error(err)
       }else {
          session.SessionSerializeInfo[key] = byteArray
       }
    }
    session.SessionId = sessionid
    return &session
}


func  DeseriazeObject(bytearray []byte, obj interface{}) bool {
    b := bytes.Buffer{}
    b.Write(bytearray)
    d := gob.NewDecoder(&b)
    err := d.Decode(obj)
    if err != nil {
       log.Error(err)
       return false
    }
    return true
}

