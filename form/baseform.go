package form

import (
  "github.com/wpxiong/beargo/log"
)


func init() {
  log.InitLog()
}


type BaseForm struct {
   Error map[string][]string
}

