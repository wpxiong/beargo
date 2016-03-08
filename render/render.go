package render

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/render/template"
  "github.com/wpxiong/beargo/memorycash"
  "github.com/wpxiong/beargo/process"
  "github.com/wpxiong/beargo/appcontext"
  "os"
  "path/filepath"
  "net/http"
)

var USE_RENDER_PROCESS bool

func init() {
  USE_RENDER_PROCESS = true 
  log.InitLog()
}

type RenderInfo struct {
   MemorycashManager *memorycash.MemoryCashManager
   Writer            *http.ResponseWriter
   FinishSingal      chan int
   TemplateList      map[string]*template.Template
   TemplateCount     int
   OutPutData        interface{}
   UrlPath           string
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
   if rendermanager.useProcess {
      renderInfo.FinishSingal <- 1
   }
   return true
}


func (this *RenderInfo) RenderProcess(interface{}) interface{} {
    log.Debug("RenderProcess Start")
    if rendermanager.useProcess {
       workjob := &process.WorkJob{ Parameter : this }
       workjob.WorkProcess = processRender
       process.AddJob(rendermanager.workprocess,workjob)
       _ = <- this.FinishSingal
       return true
    }else {
       return processRender(this)
    }
}

const (
 DEFAULT_TEMPLATE_FOLDER = "/views/"
)

var rendermanager *RenderManager = nil
var filePath string 

type RenderManager struct {
  workprocess *process.WorkProcess
  useProcess  bool
  templateList map[string]*RenderInfo
  templateFilePath string
}

func SetDefaultTemplateDir(workDir string){
   filePath = workDir + DEFAULT_TEMPLATE_FOLDER
   log.Debug("Template Dir: " + filePath)
}

func SetTemplateDir(folderPath string){
   filePath = folderPath
}

func getManager() *RenderManager {
   if rendermanager == nil {
      rendermanager = &RenderManager{workprocess : process.New(),useProcess:USE_RENDER_PROCESS,templateList:make(map[string]*RenderInfo),templateFilePath:filePath}
      if rendermanager.useProcess {
         rendermanager.workprocess.Init_Default()  
      }else {
         rendermanager.workprocess = nil
      }
   }
   return rendermanager
}

func StartTemplateManager(){
   getManager()
}

func (this *RenderManager) createRenderInfo(writer *http.ResponseWriter,urlPath string) *RenderInfo {
   if this.useProcess {
     log.Debug("URL PATH: " + urlPath)
     var info *RenderInfo = &RenderInfo{FinishSingal: make(chan int),UrlPath:urlPath}
     return info
   }else {
     return &RenderInfo{}
   }
}

func (this *RenderManager) compileTemplate() error {
   var err error = nil
   err = filepath.Walk(this.templateFilePath, 
      func(path string, info os.FileInfo, err error) error {
         if !info.IsDir() {
            rel, err := filepath.Rel(this.templateFilePath, path)
            if err != nil {
               log.Debug(rel)
            }
            return nil
         }
         rel, err := filepath.Rel(this.templateFilePath, path)
         if err != nil {
            log.Debug(rel)
         }
         return nil
    })
   if err != nil {
      log.Error(err)
   }
   return err
}

func CompileTemplate() error {
   log.Debug("CompileTemplate Start")
   if err := getManager().compileTemplate();err != nil {
      log.Error(err)
      return err
   }
   return nil
}


func CreateRenderInfo(app *appcontext.AppContext) *RenderInfo {
   return getManager().createRenderInfo(app.Writer.HttpResponseWriter,app.UrlPath)  
}