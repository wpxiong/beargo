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
