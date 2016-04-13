package moudle

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/constvalue"
  "database/sql"
  "strconv"
  "strings"
  "reflect"
  "time"
  _ "github.com/lib/pq"
)

func init() {
  log.InitLog()
}

type PostgresqlDBProvider struct {
   db *sql.DB
   
}

func (this *PostgresqlDBProvider )  SetMinConnection(max int) {
   this.db.SetMaxIdleConns(max)
}

func (this *PostgresqlDBProvider )  SetMaxConnection(max int) {
   this.db.SetMaxOpenConns(max)
}

func (this *PostgresqlDBProvider ) Begin() (*sql.Tx,error) {
   return this.db.Begin()
}

func (this *PostgresqlDBProvider )  Commit(tx *sql.Tx) error {
   return tx.Commit()
}

func (this *PostgresqlDBProvider )  Close() error {
   return this.db.Close()
}

func (this *PostgresqlDBProvider ) ConnectionDb(dburl string) error {
  var err error = nil
  this.db, err = sql.Open("postgres", dburl)
  if err != nil {
    log.Error("Postgresql Connection Error ")
  }
  return err
}

func (this *PostgresqlDBProvider ) Query(sql string,ts *Trans) (*sql.Rows ,error){
   log.Info(sql)
   if ts == nil {
     return this.db.Query(sql)
   }else {
     return ts.tx.Query(sql)
   }
}

func (this *PostgresqlDBProvider ) Insert(sql string,ts *Trans) (sql.Result ,error){
   log.Info(sql)
   if ts == nil {
     return this.db.Exec(sql)
   }else {
     return ts.tx.Exec(sql)
   }
}

func (this *PostgresqlDBProvider ) CreateTable(tableName string,sqlstr string,primaryKey []string) (sql.Result ,error) {
   if len(primaryKey) != 0 {
      sqlstr = sqlstr[0:len(sqlstr)-2] +  ",\n PRIMARY KEY("  + strings.Join(primaryKey,",") + ") );"  
   }
   log.Info(sqlstr)
   return this.db.Exec(sqlstr)
}

func  (this *PostgresqlDBProvider ) CreateForeignKey(tableName string ,  keyColumn string, refrenceTableName string, referenceColumnName string) (sql.Result ,error) {
   foreignkeyname := keyColumn + "_" + refrenceTableName + "_" + referenceColumnName
   sqlstr := "alter table `" + tableName + "` add constraint `" + foreignkeyname  + "` foreign key (`"  + keyColumn + "`) references `" + refrenceTableName + "` (`" + referenceColumnName + "`) on delete cascade on update cascade "
   log.Info(sqlstr)
   return this.db.Exec(sqlstr)
}

func  (this *PostgresqlDBProvider ) CreatePrimaryKey(tableName string,keyList []string)  (sql.Result ,error) {
   if len(keyList) != 0 {
      sqlstr := "ALTER TABLE " + tableName  + " ADD PRIMARY KEY("  + strings.Join(keyList,",") + ")" 
      log.Info(sqlstr)
      return this.db.Exec(sqlstr)  
   }else {
      return nil,nil
   }
}

func (this *PostgresqlDBProvider ) ExecuteSQL(sql string,ts *Trans) (sql.Result ,error){
   log.Info(sql)
   if ts == nil {
     return this.db.Exec(sql)
   }else {
     return ts.tx.Exec(sql)
   }
}

func (this *PostgresqlDBProvider ) Update(sql string,ts *Trans) (sql.Result ,error){
   log.Info(sql)
   if ts == nil {
     return this.db.Exec(sql)
   }else {
     return ts.tx.Exec(sql)
   }
}


func (this *PostgresqlDBProvider ) DropTable(tableName string) (sql.Result ,error){
   var sql string = "drop table if exists " + tableName + ";"
   log.Info(sql)
   return this.db.Exec(sql)
}

func (this *PostgresqlDBProvider ) GetDBIntType() string {
  return "INTEGER"
}

func (this *PostgresqlDBProvider )  GetDBInt8Type() string {
   return "SMALLINT"
}

func (this *PostgresqlDBProvider ) GetDBInt16Type() string {
   return "SMALLINT"
}

func (this *PostgresqlDBProvider )  GetDBInt32Type() string {
  return "INTEGER"
}

func (this *PostgresqlDBProvider ) GetDBInt64Type() string {
  return "BIGINT"
}

func (this *PostgresqlDBProvider )  GetDBUintType() string {
  return "INT"
}

func (this *PostgresqlDBProvider )  GetDBUint8Type() string {
   return "SMALLINT"
}

func (this *PostgresqlDBProvider )  GetDBUint16Type() string {
  return "INTEGER"
}

func (this *PostgresqlDBProvider )  GetDBUint32Type() string {
    return "BIGINT"
}

func (this *PostgresqlDBProvider )  GetDBUint64Type() string {
  return "BIGINT"
}

func (this *PostgresqlDBProvider )  GetDBFloat32Type() string {
  return "DECIMAL"
}

func (this *PostgresqlDBProvider )  GetDBFloat64Type() string {
  return "DOUBLE PRECISION"
}

func (this *PostgresqlDBProvider )  GetDBComplex64Type() string {
  return "VARCHAR(64)"
}

func (this *PostgresqlDBProvider )  GetDBComplex128Type() string {
  return "VARCHAR(128)"
}

func (this *PostgresqlDBProvider )  GetDBBoolType() string {
  return "CHAR(1)"
}

func (this *PostgresqlDBProvider ) GetInsertDBComplex64Sql(val complex128 ) string {
  return "'" + strconv.FormatFloat(real(val),'f', -1, 32) +"," + strconv.FormatFloat(imag(val),'f', -1, 32) + "'"
}

func (this *PostgresqlDBProvider ) GetInsertDBComplex128Sql(val complex128 ) string {
   return "'" + strconv.FormatFloat(real(val),'f', -1, 64) +"," + strconv.FormatFloat(imag(val),'f', -1, 64) + "'"
}

func (this *PostgresqlDBProvider )   GetInsertDBTimeSql(ti time.Time) string {
   return "'" + ti.Format(constvalue.DEFAULT_TIME_FORMATE) + "'"
}

func (this *PostgresqlDBProvider )   AppendScanComplexField(list *[]interface{}) {
    var complexField string 
    (*list) = append(*list,&complexField)
}

func (this *PostgresqlDBProvider )  PrepareExecuteSQL(sql string ,parameter []interface{},ts *Trans) {
 log.Debug(sql)
 if ts == nil {
    if smt,err := this.db.Prepare(sql);err == nil {
       if _, err1 := smt.Exec(parameter...);err1 != nil {
          panic(err)
       }
    }else {
       panic(err)
    }
 }else {
    if smt,err := this.db.Prepare(sql);err == nil {
       smt = ts.tx.Stmt(smt)
       if _, err1 := smt.Exec(parameter...);err1 != nil {
          panic(err)
       }
    }else {
       panic(err)
    }
 }
}

   
func (this *PostgresqlDBProvider )  TableExistsInDB(tableName string) (bool,error) {
   var sql string = "show tables like '" + tableName + "';"
   log.Info(sql)
   rows,err := this.db.Query(sql)
   if err != nil {
      return false,err
   }else {
      return rows.Next(),err
   }
}
   
func (this *PostgresqlDBProvider )  GetDBStringType(length int ) string {
  if length < 65535 {
     return "VARCHAR"
  }else if length < 16777215 {
     return "TEXT"
  }else {
     return "TEXT"
  }
}

func (this *PostgresqlDBProvider )  GetDBTimeType() string {
  return "TIMESTAMP"
}

func (this *PostgresqlDBProvider )  GetDBByteArrayType(length int) string {
  return "BYTEA"
}

func (this *PostgresqlDBProvider )  LimitSql( limit int ) string {
   return " LIMIT "  + strconv.Itoa(limit)
}

func (this *PostgresqlDBProvider )  CreateDefaultValue(defaultValue interface{}) string {
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


func (this *PostgresqlDBProvider ) CreateSqlTypeByLength(auto_increment bool ,sqlType string,length int, scale int) string {
   result := sqlType
   switch sqlType {
      case "INT","TINYINT","SMALLINT","BIGINT","VARCHAR" :
        if length != 0 {
          result = sqlType + "(" + strconv.Itoa(length)  + ")"
        }
      case "FLOAT","DOUBLE":
        if length != 0 && scale != 0 {
          result = sqlType + "(" + strconv.Itoa(length)  + ","  +  strconv.Itoa(scale) + ")"
        } else if scale != 0 {
          result = sqlType + "(" + strconv.Itoa(scale)  + ","  +  strconv.Itoa(scale) + ")"
        } else if length != 0 {
          result = sqlType + "(" + strconv.Itoa(length)  + ","  +  strconv.Itoa(length) + ")"
        }
      default:
        result = sqlType
   }
   if auto_increment {
     result = " SERIAL "
   }
   return result
}

