package webtemplate

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/form"
  "net/http"
  "reflect"
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
  form.BaseForm
}



func (this *HtmlTemplate) RenderHTMLTemplate(writer *http.ResponseWriter,filepathList []string,output interface{},errorInfo map[string][]string) error {
   log.Debug("RenderHTMLTemplate Start")
   tmpl := template.Must(template.ParseFiles(filepathList...))
   zeroValue := reflect.Value{}
   switch reflect.TypeOf(output).Kind() {
      case reflect.Ptr:
         errorField := reflect.ValueOf(output).Elem().FieldByName("Error")
         if errorField != zeroValue  && errorField.CanSet() {
            errorField.Set(reflect.ValueOf(errorInfo))
         }
      case reflect.Struct:
   }
   log.Debug(output)
   err := tmpl.Execute((*writer),output)
   if err != nil {
      log.ErrorArray("Render Page Error. Page:" , filepathList)
   }
   (*writer).Header().Set("Content-Type", "text/html")
   return err
}