package constvalue

import (
  "strings"
  "os"
)

const (

  BEFORE_FUNC = "Before"
  AFTER_FUNC = "After"
  
  DEFAULT_FUNC_NAME = "Index"
  
  RESOURCE_PATH_KEY = "resource_url"
  REQUEST_TIMEOUT_KEY  = "request.timeout"
  RESPONSE_TIMEOUT_KEY  = "response.timeout"
  RESOURCE_FOLDER = "./views/public/"
  SESSION_TIMEOUT_KEY = "session_time_out"
  SESSION_PROVIDER_KEY = "session_provider"
  
  SESSION_NAME = "GSESSION"
  
  BEFORE_interceptor_KEY  = "before_interceptor"
  AFTER_interceptor_KEY  = "after_interceptor"
  CASH_TYPE_KEY = "cash_type"
  CASH_MAXSIZE_KEY ="cash_max_size"
  
  DEFAULT_RESOURCE_PATH = "resource"
  DEFAULT_REQUEST_TIMEOUT = "300"
  DEFAULT_RESPONSE_TIMEOUT = "300"
  
  DEFAULT_SESSION_TIMEOUT = 3600 
  DEFAULT_CASH_MAXSIZE = 1000 
  DEFAULT_CASH_TYPE = "memory"
  DEFAULT_SESSION_PROVIDER = "MemorySessionProvider"
  DEFAULT_SESSIONID_SIZE = 32
  
  //interceptor Name
  ParameterParseinterceptor = "ParameterParseinterceptor"
  ParameterBinderinterceptor = "ParameterBinderinterceptor"
  RenderBindinterceptor = "RenderBindinterceptor"
  RenderOutPutinterceptor = "RenderOutPutinterceptor"
  Redirectinterceptor = "Redirectinterceptor"
  Sessioninterceptor = "Sessioninterceptor"
  Xsrfinterceptor = "Xsrfinterceptor"
  DBtransactionEndinterceptor = "DBtransactionEndinterceptor"
  DBtransactionStartinterceptor = "DBtransactionStartinterceptor"
  
  MemorySessionProvider = "MemorySessionProvider"
  
  
  
  ERROR_403 = "/error_403"
  ERROR_404 = "/error_404"
  ERROR_405 = "/error_405"
  ERROR_500 = "/error_500"
  
  ERROR_403_PATH_KEY = "403_error"
  ERROR_404_PATH_KEY = "404_error"
  ERROR_405_PATH_KEY = "405_error"
  ERROR_500_PATH_KEY = "500_error"

  
  REDIRECT_ERROR = "REDIRECT_ERROR"
  
  XSRF_TOKEN = "XSRF_TOKEN"
  
  
  DB_ID = "id"
  DB_COLUMN_NAME = "column_name"
  DB_NOT_NULL = "notnull"
  DB_LENGTH = "length"
  DB_SCALE = "scale"
  DB_UNIQUE_KEY = "unique_key"
  DB_DEFAULT_VALUE = "default_value"
  DB_AUTO_INCREMENT = "auto_increment"
  
  DB_RELATION_TYPE = "relation_type"
  DB_REFERENCED_COLUMN_NAME = "referenced_column_name"
  
  DB_RELATION_ONE_TO_MANY = "onetomany"
  DB_RELATION_MANY_TO_ONE = "manytoone"
  DB_RELATION_ONE_TO_ONE = "onetoone"
  
  DEFAULT_TIME_FORMATE = "2006-01-02 15:04:06"
  
  
  DEFAULT_MAX_DB_CONNECTION = 100
  DEFAULT_MIN_DB_CONNECTION = 5
  
)

var  DEFAULT_ERROR_403_PATH = "error/403"
var  DEFAULT_ERROR_404_PATH = "error/404"
var  DEFAULT_ERROR_405_PATH = "error/405"
var  DEFAULT_ERROR_500_PATH = "error/500"
  

var DEFULT_BEFORE_interceptor []string
var DEFULT_AFTER_interceptor []string

func init(){
  DEFULT_BEFORE_interceptor = []string{"ParameterParseinterceptor","ParameterBinderinterceptor","Sessioninterceptor","Xsrfinterceptor","DBtransactionStartinterceptor"}
  DEFULT_AFTER_interceptor = []string {"Redirectinterceptor","DBtransactionEndinterceptor","RenderBindinterceptor","RenderOutPutinterceptor"}
  
  DEFAULT_ERROR_403_PATH = strings.Replace(DEFAULT_ERROR_403_PATH, "/", string(os.PathSeparator), -1)
  DEFAULT_ERROR_404_PATH = strings.Replace(DEFAULT_ERROR_404_PATH, "/", string(os.PathSeparator), -1)
  DEFAULT_ERROR_405_PATH = strings.Replace(DEFAULT_ERROR_405_PATH, "/", string(os.PathSeparator), -1)
  DEFAULT_ERROR_500_PATH = strings.Replace(DEFAULT_ERROR_500_PATH, "/", string(os.PathSeparator), -1)
}

