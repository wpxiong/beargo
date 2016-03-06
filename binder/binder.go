package binder

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "reflect"
  "strconv"
  "strings"
  "errors"
)

func init() {
  log.InitLog()
}

var floatType = reflect.TypeOf(float64(0))


func getFloat(unk interface{}) (float64, error) {
    v := reflect.ValueOf(unk)
    v = reflect.Indirect(v)
    if !v.Type().ConvertibleTo(floatType) {
        return 0, errors.New("")
    }
    fv := v.Convert(floatType)
    return fv.Float(), nil
}

func GetParamValueFloat(paramValue interface{}) (float64,error) {
   switch paramValue.(type) {
     case float32:
       return getFloat(paramValue)
     case float64:
       return getFloat(paramValue)
     case string:
       val,err:= strconv.ParseFloat(paramValue.(string),64)
       if err == nil {
          return val,nil
       }else {
          return 0,err
       }
     default:
        return 0,errors.New("")
  }
}

func GetParamValueInt(paramValue interface{}) (int64,error) {
  switch paramValue.(type) {
     case int,int8,int16,int32,int64:
       return paramValue.(int64),nil
     case string:
       val,err:= strconv.Atoi(paramValue.(string))
       if err == nil {
          return int64(val),nil
       }else {
          return 0,err
       }
     default:
        return 0,errors.New("")
  }
}

func GetParamValueUint(paramValue interface{}) (uint64,error) {
  switch paramValue.(type) {
     case uint,uint8,uint16,uint32,uint64:
       return paramValue.(uint64),nil
     case string:
       val,err:= strconv.Atoi(paramValue.(string))
       if err == nil {
          return uint64(val),nil
       }else {
          return 0,err
       }
     default:
        return 0,errors.New("")
  }
}



func BinderByType(f reflect.Value, ft reflect.Type, param map[string] interface{},name string){
  switch ft.Kind() {
    case reflect.Map:
      f.Set(reflect.MakeMap(ft))
    case reflect.Slice:
      if param[strings.ToLower(name)] != nil {
         BinderSlice(strings.ToLower(name),&f,ft,param[strings.ToLower(name)])
      }
    case reflect.Chan:
      f.Set(reflect.MakeChan(ft, 0))
    case reflect.Struct:
      initializeStruct(ft, f,param)
    case reflect.Ptr:
      fv := reflect.New(ft.Elem())
      initializeStruct(ft.Elem(), fv.Elem(),param)
      f.Set(fv)
    case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
      if param[strings.ToLower(name)] != nil {
        BinderInt(&f,param[strings.ToLower(name)])
      }
    case reflect.Float32,reflect.Float64:
      if param[strings.ToLower(name)] != nil {
        BinderFloat(&f,param[strings.ToLower(name)])
      }
    case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
      if param[strings.ToLower(name)] != nil {
        BinderUint(&f,param[strings.ToLower(name)])
      }
    case reflect.String:
      if param[strings.ToLower(name)] != nil {
        log.Debug("66666")
        BinderString(&f,param[strings.ToLower(name)])
      }
    default:
  }
}

func initializeStruct(t reflect.Type, v reflect.Value,param map[string] interface{}) {
  for i := 0; i < v.NumField(); i++ {
    f := v.Field(i)
    ft := t.Field(i)
    BinderByType(f,ft.Type,param,ft.Name)
  }
}

func BinderSlice(name string ,field *reflect.Value ,filedtype reflect.Type, paramValue interface{}){
   
   switch paramValue.(type) {
      case []string :
        slice := reflect.MakeSlice(filedtype, len(paramValue.([]string)), len(paramValue.([]string)))
        //element := slice.Index(0)
        tf := reflect.TypeOf("dd")
        for k,str := range paramValue.([]string) {
           log.Debug("1111")
           arrayElement := slice.Index(k)
           log.Debug("2222")
           v := reflect.ValueOf(&arrayElement).Elem()
           log.Debug("3333")
           mapparam := make(map[string] interface{})
           mapparam[name] = str
           log.Debug("44444")
           log.Debug(tf)
           //BinderByType(v,tf,mapparam,name)
           log.Debug(v)
           log.Debug(tf)
           log.Debug(mapparam)
           log.Debug(name)
           log.Debug("555")
        }
        field.Set(slice)
      case string :
        slice := reflect.MakeSlice(filedtype, 1, 1)
        field.Set(slice)
      default:
   }
   
}

func BinderInt(field *reflect.Value , paramValue interface{}){
  intValue,err := GetParamValueInt(paramValue)
  intVal := int64(intValue)
  if err== nil && !field.OverflowInt(intVal) {
     field.SetInt(intVal)
  }
}

func BinderFloat(field *reflect.Value , paramValue interface{}){
  intValue,err := GetParamValueFloat(paramValue)
  intVal := float64(intValue)
  if err== nil && !field.OverflowFloat(intVal) {
     field.SetFloat(intVal)
  }
}

func BinderUint(field *reflect.Value , paramValue interface{}){
  intValue,err := GetParamValueUint(paramValue)
  intVal := uint64(intValue)
  if err== nil && !field.OverflowUint(intVal) {
     field.SetUint(intVal)
  }
}

func BinderString(field *reflect.Value , paramValue interface{}){
   field.SetString(paramValue.(string))
}

func Binder(field *reflect.Value , paramValue interface{}) {
  if !field.CanSet() {
      switch field.Kind() {
        case reflect.Int:
          BinderInt(field,paramValue)
        case reflect.Int8:
          BinderInt(field,paramValue)
        case reflect.Int16:
          BinderInt(field,paramValue)
        case reflect.Int64:
          BinderInt(field,paramValue)
        default:
          log.Debug("Default Binder")
      }  
  }
}

func BinderParameter(appcon *appcontext.AppContext){
   log.Debug("BinderParameter Start")
   v := reflect.New(appcon.FormType)
   initializeStruct(appcon.FormType, v.Elem(),appcon.Parameter)
   appcon.Form = v.Interface()
}