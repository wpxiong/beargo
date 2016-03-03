package util

import (
  "strconv"
  "strings"
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}

func StringToInt(valStr string) interface{}  {
   i, err := strconv.Atoi(valStr)
   if err != nil {
      return i
   }else {
      return nil
   }
}

func StringToFloat(valStr string) interface{} {
   i, err := strconv.ParseFloat(valStr, 64)
   if err != nil {
      return i
   }else {
      return nil
   }
}


func StringToDouble(valStr string) interface{} {
   i, err := strconv.ParseFloat(valStr, 64)
   if err != nil {
      return i
   }else {
      return nil
   }
}

func StringToBool(valStr string) interface{} {
  if strings.ToLower(valStr) == "true" {
     return true
  }else if strings.ToLower(valStr) == "false" {
     return false
  }else {
     return nil
  }
}
