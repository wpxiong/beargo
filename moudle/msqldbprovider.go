package moudle

import (
  "github.com/wpxiong/beargo/log"
  "database/sql"
  "strconv"
  "reflect"
  _ "github.com/go-sql-driver/mysql"
)

func init() {
  log.InitLog()
}

type MysqlDBProvider struct {
   db *sql.DB
}

func (this *MysqlDBProvider ) ConnectionDb(dburl string) error {
  var err error = nil
  this.db, err = sql.Open("mysql", dburl)
  if err != nil {
    log.Error("Mysql Connection Error ")
  }
  return err
}

func (this *MysqlDBProvider ) Query(sql string) (*sql.Rows ,error){
   log.Info(sql)
   return this.db.Query(sql)
}

func (this *MysqlDBProvider ) Insert(sql string) (sql.Result ,error){
   log.Info(sql)
   return this.db.Exec(sql)
}

func (this *MysqlDBProvider ) CreateTable(sql string) (sql.Result ,error){
   log.Info(sql)
   return this.db.Exec(sql)
}


func (this *MysqlDBProvider ) DropTable(tableName string) (sql.Result ,error){
   var sql string = "drop table if exists " + tableName + ";"
   log.Info(sql)
   return this.db.Exec(sql)
}

func (this *MysqlDBProvider ) GetDBIntType() string {
  return "INT"
}

func (this *MysqlDBProvider )  GetDBInt8Type() string {
   return "TINYINT"
}

func (this *MysqlDBProvider ) GetDBInt16Type() string {
   return "SMALLINT"
}

func (this *MysqlDBProvider )  GetDBInt32Type() string {
  return "SMALLINT"
}

func (this *MysqlDBProvider ) GetDBInt64Type() string {
  return "BIGINT"
}

func (this *MysqlDBProvider )  GetDBUintType() string {
  return "INT"
}

func (this *MysqlDBProvider )  GetDBUint8Type() string {
   return "TINYINT"
}

func (this *MysqlDBProvider )  GetDBUint16Type() string {
  return "SMALLINT"
}

func (this *MysqlDBProvider )  GetDBUint32Type() string {
    return "INT"
}

func (this *MysqlDBProvider )  GetDBUint64Type() string {
  return "BIGINT"
}

func (this *MysqlDBProvider )  GetDBFloat32Type() string {
  return "FLOAT"
}

func (this *MysqlDBProvider )  GetDBFloat64Type() string {
  return "DOUBLE"
}

func (this *MysqlDBProvider )  GetDBComplex64Type() string {
  return "VARCHAR(64)"
}

func (this *MysqlDBProvider )  GetDBComplex128Type() string {
  return "VARCHAR(128)"
}

func (this *MysqlDBProvider )  GetDBBoolType() string {
  return "CHAR(1)"
}

func (this *MysqlDBProvider )  GetDBStringType() string {
  return "VARCHAR"
}

func (this *MysqlDBProvider )  GetDBTimeType() string {
  return "TIMESTAMP"
}

func (this *MysqlDBProvider )  CreateDefaultValue(defaultValue interface{}) string {
  if defaultValue == nil {
     return ""
  }
  t := reflect.TypeOf(defaultValue)
  switch t.Name() {
     case "int64":
       val := defaultValue.(int64) 
       return " DEFAULT " + strconv.FormatInt(val,10) + " "
     case "uint64":
       val := defaultValue.(uint64)
       return " DEFAULT " + strconv.FormatUint(val,10) + " "
     case "float64":
       val := defaultValue.(float64) 
       return " DEFAULT " + strconv.FormatFloat(val, 'f', -1, 64) + " "
     case "complex64":   
       val := defaultValue.(complex64)
       return " DEFAULT '" + strconv.FormatFloat(float64(real(val)), 'f', -1, 32) + "," +  strconv.FormatFloat(float64(imag(val)), 'f', -1, 32) +  "' "
     case "complex128":
       val := defaultValue.(complex128) 
       return " DEFAULT '" + strconv.FormatFloat(real(val), 'f', -1, 64) + "," +  strconv.FormatFloat(imag(val), 'f', -1, 64) +  "' "
     case "bool":
       val := defaultValue.(bool) 
       if val {
          return " DEFAULT '" + "1" + "' "
       }else {
          return " DEFAULT '" + "0" + "' "
       }
     case "string":
       val := defaultValue.(string)
       if val == "" {
          return ""
       }else {
          return " DEFAULT '" + defaultValue.(string) + "' "
       }
     default :
       return ""
  }
  return ""
}


func (this *MysqlDBProvider ) CreateSqlTypeByLength(sqlType string,length int, scale int) string {
   switch sqlType {
      case "INT","TINYINT","SMALLINT","BIGINT","VARCHAR" :
        if length != 0 {
          return sqlType + "(" + strconv.Itoa(length)  + ")"
        }else {
          return sqlType
        }
      case "FLOAT","DOUBLE":
        if length != 0 && scale != 0 {
          return sqlType + "(" + strconv.Itoa(length)  + ","  +  strconv.Itoa(scale) + ")"
        } else if scale != 0 {
          return sqlType + "(" + strconv.Itoa(scale)  + ","  +  strconv.Itoa(scale) + ")"
        } else if length != 0 {
          return sqlType + "(" + strconv.Itoa(length)  + ","  +  strconv.Itoa(length) + ")"
        } else {
          return sqlType
        }
      default:
        return sqlType
   }
}

