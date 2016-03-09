package constvalue


const (
  RESOURCE_PATH_KEY = "resource_url"
  REQUEST_TIMEOUT_KEY  = "request.timeout"
  RESPONSE_TIMEOUT_KEY  = "response.timeout"
  RESOURCE_FOLDER = "./views/public/"
  
  BEFORE_FILTER_KEY  = "before_filter"
  AFTER_FILTER_KEY  = "after_filter"
  
  DEFAULT_RESOURCE_PATH = "resource"
  DEFAULT_REQUEST_TIMEOUT = "300"
  DEFAULT_RESPONSE_TIMEOUT = "300"
)

var DEFULT_BEFORE_FILTER []string
var DEFULT_AFTER_FILTER []string

func init(){
  DEFULT_BEFORE_FILTER = []string{"ParameterParseFilter","ParameterBinderFilter"}
  DEFULT_AFTER_FILTER = []string {"RenderBindFilter","RenderOutPutFilter"}
}

