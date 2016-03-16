package moudle

import (
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}

type DbProviderInterface interface {
   ConnectionDb(dburl string) error 
}

