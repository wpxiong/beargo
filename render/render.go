package render

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/render/template"
  "github.com/wpxiong/beargo/memorycash"
  "github.com/wpxiong/beargo/process"
  "github.com/wpxiong/beargo/appcontext"
  //"io/ioutil"
  "net/http"
)

func init() {
  log.InitLog()
}

type RenderInfo struct {
   MemorycashManager *memorycash.MemoryCashManager
   Writer            *http.ResponseWriter
   FinishSingal      chan int
   TemplateList      map[string]*template.Template
   TemplateCount     int
   OutPutData        interface{}
}

func processRender(param interface{}) interface{} {
   log.Debug("processRender Start")
   var renderInfo *RenderInfo = param.(*RenderInfo)
   for _,temp := range renderInfo.TemplateList {
      err := temp.Render(renderInfo.OutPutData)
      if err != nil {
         log.Error(err)
         break
      }
   }
   renderInfo.FinishSingal <- 1
   return true
}


func (this *RenderInfo) RenderProcess(interface{}) interface{} {
    log.Debug("RenderProcess Start")
    workjob := &process.WorkJob{ Parameter : this }
    workjob.WorkProcess = processRender
    process.AddJob(workjob)
    _ = <- this.FinishSingal
    return true
}


var rendermanager *RenderManager = nil

type RenderManager struct {
}

func getManager() *RenderManager {
   if rendermanager == nil {
      rendermanager = &RenderManager{}
   }
   return rendermanager
}

func (this *RenderManager) createRenderInfo(writer *http.ResponseWriter) *RenderInfo {
   var info *RenderInfo = &RenderInfo{FinishSingal: make(chan int)}
   return info
}


func CreateRenderInfo(app *appcontext.AppContext) *RenderInfo {
   return getManager().createRenderInfo(app.Writer.HttpResponseWriter)  
}