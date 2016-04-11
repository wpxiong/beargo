package webapp

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/constvalue"
  "net"
  "net/rpc"
  "os"
  "strconv"
)


const (
   STOP_COMMAND = "stop"
   RESTART_COMMAND = "restart"
)
 
func init() {
  log.InitLog()
}

type Args struct {
    Command string
}

type CommandInterface struct {
   webapp  *WebApplication
}


func (t *CommandInterface) SendCommand(cmd *Args,reply *int) error {
    switch cmd.Command {
       case STOP_COMMAND:
         t.webapp.Stop()
       case RESTART_COMMAND:
         t.webapp.ReStart()
       default:
       
    }
    *reply = 200
    return nil
}


func checkError(err error) {
   if err != nil {
       log.Error(err)
       os.Exit(1)
   }
}

func startCommanListener(web *WebApplication){
   command := new(CommandInterface)
   command.webapp = web
   rpc.Register(command)
   var port string = web.AppContext.GetConfigValue(constvalue.MANAGER_PORT,"").(string)
   val,err:= strconv.Atoi(port)
   var managerport int 
   if err == nil {
      managerport = val
   }else {
      managerport = constvalue.DEFAULT_MANAGER_PORT
   }
   var host string = web.AppContext.GetConfigValue(constvalue.MANAGER_HOST,constvalue.DEFAULT_MANAGER_HOST).(string)
   tcpAddr, err := net.ResolveTCPAddr("tcp", host + ":" + strconv.Itoa(managerport))
   checkError(err)
   listener, err := net.ListenTCP("tcp", tcpAddr)
   checkError(err)
   for {
      conn, err := listener.Accept()
      if err != nil {
         continue
      }
      rpc.ServeConn(conn)
   }

}