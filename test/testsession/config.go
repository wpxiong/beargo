package main

import (
  "github.com/wpxiong/beargo/interceptor"
  "github.com/wpxiong/beargo/session"
  "github.com/wpxiong/beargo/webapp"
  "github.com/wpxiong/beargo/render/template"
)

func InitConfig() webapp.ConfigMap{
  config := webapp.ConfigMap{}
  config.InterceptoFuncmap = InitFiltrConfig()
  config.Sessionprovidermap = InitSessionProviderConfig()
  config.Templatefuncmap = InitTemplatefuncConfig()
  return config
}

func InitFiltrConfig() map[string]interceptor.InterceptorFunc {
  funcMap := make(map[string]interceptor.InterceptorFunc)
  return funcMap
}


func InitSessionProviderConfig() map[string]session.SessionProvider {
  sessionProviderMap := make(map[string]session.SessionProvider)  
  return sessionProviderMap
}

func InitTemplatefuncConfig() template.TemplateFuncMap {
  templateFuncMap := make(template.TemplateFuncMap)  
  return templateFuncMap
}