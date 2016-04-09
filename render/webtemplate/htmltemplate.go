package webtemplate

import (
  "github.com/wpxiong/beargo/log"
  "net/http"
  "html/template"
)

func init() {
  log.InitLog()
}

type HtmlTemplate struct {
   Template
}

type  Indform struct{
  Name   string
  Password  string
}



func (this *HtmlTemplate) RenderHTMLTemplate(writer *http.ResponseWriter,filepathList []string,output interface{},errorInfo map[string][]string,useLayout bool,layoutName string) error {
   log.Debug("RenderHTMLTemplate Start")
   tmpl := template.Must(template.ParseFiles(filepathList...))
   var err error
   if useLayout == true {
      err = tmpl.ExecuteTemplate((*writer),layoutName, output)
   }else {
      err = tmpl.Execute((*writer),output)
   }
   if err != nil {
      log.ErrorArray("Render Page Error. Page:" , filepathList)
   }
   (*writer).Header().Set("Content-Type", "text/html")
   return err
}