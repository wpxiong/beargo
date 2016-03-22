package moudle

import (
  "github.com/wpxiong/beargo/log"
  "database/sql"
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
   return this.db.Query(sql)
}

func (this *MysqlDBProvider ) Insert(sql string) (sql.Result ,error){
   return this.db.Exec(sql)
}

func (this *MysqlDBProvider ) CreateTable(sql string) (sql.Result ,error){
   return this.db.Exec(sql)
}


func (this *MysqlDBProvider ) GetDBIntType() string {
  return "INT "
}

func (this *MysqlDBProvider )  GetDBInt8Type() string {
   return "TINYINT "
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
  return "INT "
}

func (this *MysqlDBProvider )  GetDBUint8Type() string {
   return "TINYINT "
}

func (this *MysqlDBProvider )  GetDBUint16Type() string {
  return "SMALLINT"
}

func (this *MysqlDBProvider )  GetDBUint32Type() string {
    return "INT "
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
  return "char(1)"
}

