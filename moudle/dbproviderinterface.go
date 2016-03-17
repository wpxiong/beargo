package moudle

import (
  "github.com/wpxiong/beargo/log"
  "database/sql"
)

func init() {
  log.InitLog()
}

type DbProviderInterface interface {
   ConnectionDb(dburl string) error 
   Query(sql string) (*sql.Rows ,error)
   Insert(sql string) (sql.Result ,error)
   CreateTable(sql string) (sql.Result ,error)
}

