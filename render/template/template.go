package template

import (
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}

type templateType int

const (
    HTML_TEMPLATE   templateType = iota
    JSON_TEMPLATE 
    XML_TEMPLATE
    TEXT_TEMPLATE    
)

const (
   XML_TEMPLATE_EXTENSION = ".txml"
   JSON_TEMPLATE_EXTENSION = ".tjson"
   HTML_TEMPLATE_EXTENSION = ".thtml"
   TEXT_TEMPLATE_EXTENSION = ".ttxt"
)
    
  
type Template struct {
  Templatetype           templateType
}

func (this *Template) Render(OutData interface{}) error {
  return nil
}