package moudle

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/util/dbutil"
  "reflect"
  "strings"
)

func init() {
  log.InitLog()
}


type DbDialectType int

const (
    MYSQL   DbDialectType = iota
    SQLITE
    POSTGRESQL
)

type DBTableInfo struct {
  DbName        string
  DbSchema      string
  DbStuct       interface{}
  FiledList     map[string] reflect.Type
  FiledNameMap  map[string] string
  DbTableExist  bool
}

type Moudle struct {
  DbDialect        DbDialectType
  DbName           string
  DbConnectionUrl  string
  DbUserName       string
  DbPassword       string
  DbTableInfo      map[string]DBTableInfo
  DbProiver        DbProviderInterface
  connectionStatus bool
}

type DBTableInterface interface {
   GetDbTableName() string
}

type DBTableInterface interface {
   GetDbTableName() string
}


type DBTable struct {

}

func (this *DBTable) GetDbTableName() string {
   return reflect.TypeOf(*this).Name()
}


func CreateModuleInstance(DbDialect DbDialectType,DbName string,DbConnectionUrl string, DbUserName string,DbPassword string) *Moudle {
   module :=  &Moudle{DbDialect:DbDialect,DbName:DbName,DbConnectionUrl:DbConnectionUrl,DbUserName:DbUserName,DbPassword:DbPassword}
   module.initModuleInstance()
   return module
}

func (this *Moudle) initModuleInstance(){
   this.DbTableInfo = make(map[string]DBTableInfo)
   var connectionUrl string  
   switch this.DbDialect {
      case MYSQL :
        connectionUrl = this.DbUserName  + ":" + this.DbPassword +   "@" + this.DbConnectionUrl + "/" + this.DbName;
        this.DbProiver = &MysqlDBProvider{}
      case POSTGRESQL :
   }
   err := this.DbProiver.ConnectionDb(connectionUrl)
   if err != nil {
      log.Error("DB Connection Error!")
   }else {
      this.connectionStatus = true
   }
}



func (this *Moudle) AddTable(dbtable interface{}){
  if !this.connectionStatus {
     return 
  }else {
     dbInfo := DBTableInfo{}
     dbname := strings.ToLower(reflect.TypeOf(dbtable).Name())
     fieldNum := reflect.TypeOf(dbtable).NumField()
     dbInfo.FiledList = make(map[string]reflect.Type)
     dbInfo.FiledNameMap = make(map[string]string)
     for i:=0;i<fieldNum;i++{
         field := reflect.TypeOf(dbtable).Field(i)
         dbInfo.FiledList[field.Name] = field.Type
         dbInfo.FiledNameMap[strings.ToLower(field.Name)] = field.Name
     }
     dbInfo.DbName = dbname
     dbInfo.DbSchema = ""
     dbInfo.DbStuct = dbtable
     this.DbTableInfo[dbname] = dbInfo
     dbutil.GetCreateTableSql()
     log.Debug(reflect.TypeOf(dbtable).Elem().Kind()) 
     this.DbTableInfo[dbtable.GetDbTableName()] = dbInfo
     log.Debug(this.DbTableInfo)
  }

}