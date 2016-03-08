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


func (this *HtmlTemplate) RenderHTMLTemplate(writer *http.ResponseWriter,filepathList []string,output interface{}) error {
   log.Debug("RenderHTMLTemplate Start")
   tmpl := template.Must(template.ParseFiles(filepathList...))
   err := tmpl.Execute((*writer), output)
   (*writer).Header().Set("Content-Type", "text/html")
   if err != nil {
      log.ErrorArray("Render Page Error. Page:" , filepathList)
   }
   return err
}