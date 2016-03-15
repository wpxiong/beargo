package moudle

import (
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}

type DBTable struct {

}

func (this *DBTable) CreateDB (){

}

func (this *DBTable) DeleteDB (){

}


type DBTableInfo struct {
  DbName        string
  DbSchema      string
  DbStuct       DBTable
  FiledList     map[string] reflect.Type
  DbTableExist  bool
}


type Moudle stuct {
  DbDialect        string
  DbPort           int
  DbConnectionUrl  string
  DbUserName       string
  DbPassword       string
  DbTableInfo       map[string]DBTableInfo
}

func CreateModuleInstance(DbDialect string,DbPort int,DbConnectionUrl string, DbUserName string,DbPassword string) *Moudle {
   return &Moudle{DbDialect:DbDialect,DbPort:DbPort,DbConnectionUrl:DbConnectionUrl,DbUserName:DbUserName,DbPassword:DbPassword}
}