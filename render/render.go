package render

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/render/webtemplate"
  "github.com/wpxiong/beargo/memorycash"
  "github.com/wpxiong/beargo/process"
  "github.com/wpxiong/beargo/appcontext"
  "os"
  "strings"
  "path/filepath"
  "net/http"
  "errors"
)

var USE_RENDER_PROCESS bool

const (
 DEFAULT_TEMPLATE_FOLDER = "/views/"
 INCLUDE_FOLDER = "include"
 LAYOUT_FOLDER = "layout"
 VIEW_FOLDER = "view"
 DEFAULT_INIT_SIZE  = 20
)

var rendermanager *RenderManager = nil
var filePath string 

func init() {
  USE_RENDER_PROCESS = true 
  log.InitLog()
}

type RenderInfo struct {
   MemorycashManager *memorycash.MemoryCashManager
   Writer            *http.ResponseWriter
   FinishSingal      chan int
   TemplateList      map[string]*webtemplate.Template
   TemplateCount     int
   OutPutData        interface{}
   UrlPath           string
}

func (this *RenderInfo) InitRenderInfo(app *appcontext.AppContext) {
   this.TemplateList = make(map[string]*webtemplate.Template)
   path := app.UrlPath
   if strings.HasPrefix(app.UrlPath,"/"){
     path = app.UrlPath[1:]
   }
   this.TemplateList[app.UrlPath] =  getManager().pagetemplateList[path]
}

func processRender(param interface{}) interface{} {
   log.Debug("processRender Start")
   var renderInfo *RenderInfo = param.(*RenderInfo)
   size := len(renderInfo.TemplateList)
   filepathList := make([]string,0,size)
   var firstTemplate *webtemplate.Template 
   for _,temp := range renderInfo.TemplateList {
      if firstTemplate == nil && temp != nil {
         firstTemplate = temp
      }
      if temp != nil {
        filepathList = append(filepathList, DEFAULT_TEMPLATE_FOLDER[1:]  +  temp.FilePath)
      }
   }
   var err error = errors.New("Template File is not Found")
   if firstTemplate != nil {
      err = firstTemplate.Render(renderInfo.Writer,filepathList,renderInfo.OutPutData)
   }
   if err != nil {
      log.Error("Render Page Failture")
   }
   if rendermanager.useProcess {
      renderInfo.FinishSingal <- 1
   }
   return true
}


func (this *RenderInfo) RenderProcess(data interface{}) interface{} {
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


type RenderManager struct {
  workprocess *process.WorkProcess
  useProcess  bool
  templateList map[string]*webtemplate.Template
  includetemplateList []*webtemplate.Template
  layouttemplateList []*webtemplate.Template
  pagetemplateList map[string]*webtemplate.Template
  templateFilePath string
}

func SetDefaultTemplateDir(workDir string){
   filePath = workDir + DEFAULT_TEMPLATE_FOLDER
}

func SetTemplateDir(folderPath string){
   filePath = folderPath
}

func getManager() *RenderManager {
   if rendermanager == nil {
      rendermanager = &RenderManager{workprocess : process.New(),useProcess:USE_RENDER_PROCESS,templateList:make(map[string]*webtemplate.Template,DEFAULT_INIT_SIZE),templateFilePath:filePath}
      rendermanager.includetemplateList = make([]*webtemplate.Template,0,DEFAULT_INIT_SIZE)
      rendermanager.layouttemplateList  = make([]*webtemplate.Template,0,DEFAULT_INIT_SIZE)
      rendermanager.pagetemplateList  = make(map[string]*webtemplate.Template,DEFAULT_INIT_SIZE)
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
     var info *RenderInfo = &RenderInfo{FinishSingal: make(chan int),UrlPath:urlPath}
     info.Writer = writer
     return info
   }else {
     return &RenderInfo{}
   }
}
    
func (this *RenderManager) parseTemplateFile(templateFile *webtemplate.Template) {
   filePath := filepath.Join(this.templateFilePath,templateFile.FilePath)
   var ext string = filepath.Ext(filePath)
   if strings.HasPrefix(templateFile.FilePath,INCLUDE_FOLDER) {
      this.includetemplateList = append(rendermanager.includetemplateList,templateFile)
   } else if strings.HasPrefix(templateFile.FilePath,LAYOUT_FOLDER){
      this.layouttemplateList = append(rendermanager.layouttemplateList,templateFile)
   } else if strings.HasPrefix(templateFile.FilePath,VIEW_FOLDER) {
     var startIndex int = len(VIEW_FOLDER)
     var path string = (templateFile.FilePath)[startIndex + 1:]
     path = path[: len(path) - len(ext)]
     this.pagetemplateList[path] = templateFile
   }
   switch ext {
      case webtemplate.HTML_TEMPLATE_EXTENSION:
        templateFile.Templatetype =  webtemplate.HTML_TEMPLATE
      case webtemplate.XML_TEMPLATE_EXTENSION:
        templateFile.Templatetype =  webtemplate.XML_TEMPLATE
      case webtemplate.JSON_TEMPLATE_EXTENSION:
        templateFile.Templatetype =  webtemplate.JSON_TEMPLATE
      case webtemplate.TEXT_TEMPLATE_EXTENSION:
        templateFile.Templatetype =  webtemplate.TEXT_TEMPLATE
   }
}

func (this *RenderManager) compileTemplate() error {
   var err error = nil
   err = filepath.Walk(this.templateFilePath, 
      func(path string, info os.FileInfo, err error) error {
         if !info.IsDir() {
            rel, err := filepath.Rel(this.templateFilePath, path)
            if err == nil {
               this.templateList[rel] = &webtemplate.Template{FilePath:rel}
               this.parseTemplateFile(this.templateList[rel])
            }
            return nil
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