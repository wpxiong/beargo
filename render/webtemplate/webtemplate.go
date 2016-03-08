package webtemplate

import (
  "github.com/wpxiong/beargo/log"
  "net/http"
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
  Templatetype      templateType
  FilePath          string
  FileContent       string
  HasLoad           bool      
}


func (this *Template) Render(writer *http.ResponseWriter,filepathList []string,output interface{}) error {
  var err error = nil
  switch this.Templatetype {
    case HTML_TEMPLATE:
      htmltemp := &HtmlTemplate{(*this)}
      err = htmltemp.RenderHTMLTemplate(writer,filepathList,output)
    case JSON_TEMPLATE:
    case XML_TEMPLATE:
    case TEXT_TEMPLATE:
    default:
      log.Error("There is no Template File")
  }
  return err
}