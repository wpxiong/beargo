package template

import (
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}

type Template struct {
 
}

func (this *Template) Render(OutData interface{}) error {
  return nil
}