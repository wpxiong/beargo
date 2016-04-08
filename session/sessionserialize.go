package session

import (
  "github.com/wpxiong/beargo/log"
  "encoding/base64"
  "encoding/gob"
  "bytes"
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
    b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(session)
    if err != nil {
      log.Error(err)
      return ""
    }
    return base64.StdEncoding.EncodeToString(b.Bytes())
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
    return &sessionInfo
}

func DeserializeSession(sessionStr string) *Session {
    session := Session{}
    byteArray, err := base64.StdEncoding.DecodeString(sessionStr)
    if err != nil { 
       log.Error(err)
       return nil
    }
    b := bytes.Buffer{}
    b.Write(byteArray)
    d := gob.NewDecoder(&b)
    err = d.Decode(&session)
    if err != nil {
       log.Error(err)
       return nil
    }
    return &session
}
