package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"go_remote_control/base"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	ControllerIP   = "127.0.0.1" //控制端连接IP
	ControllerPort = "9900"      //控制端连接端口
)

var (
	ServerAddr      string
	DestinationAddr string
	onlineMap       = make(map[string]bool)
)

type Controller struct {
	base.Node
}

func main() {
	//tcp连接建立
	conn := getConn()
	controller := Controller{
		Node: base.Node{
			Conn:      conn,
			Addr:      conn.LocalAddr().String(),
			ReadChan:  make(chan base.Message),
			WriteChan: make(chan base.Message),
			Enc:       gob.NewEncoder(conn),
			Dec:       gob.NewDecoder(conn),
		},
	}
	//开启读写通道并读取控制台命令
	go controller.Read()
	go controller.Write()
	go controller.GetCommand()
	go controller.GetResult()
	go controller.keepAlive()

	select {}

}

func getConn() net.Conn {
	for {
		//conn, err := net.Dial("tcp", "42.192.69.243:6666")
		conn, err := net.Dial("tcp", ControllerIP+":"+ControllerPort)
		if err != nil {
			fmt.Println("client dial err=", err)
			continue
		}
		time.Sleep(10 * time.Second)
		return conn
	}
}

func (c *Controller) keepAlive() {
	for {
		msg := base.Message{
			Type:       "alive",
			CreateTime: time.Now(),
			ModifyTime: time.Now(),
			Src:        c.Addr,
			Dst:        ServerAddr,
			Content:    c.Addr + "alive",
			Log:        nil,
		}
		c.WriteChan <- msg
		time.Sleep(10 * time.Second)
	}
}

func (c *Controller) GetCommand() {
	reader := bufio.NewReader(os.Stdin) //os.Stdin 代表标准输入[终端]

	for {
		//从终端读取用户命令
		fmt.Print("Remote@", DestinationAddr, ">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("readString err=", err)
		}
		//如果用户输入的是 exit就退出
		line = strings.Trim(line, " \r\n")
		log.Println("发送命令=", line)
		if DestinationAddr == "" {
			fmt.Println("Please set destination!")
		}
		switch line {
		case "exit":
			fmt.Println("客户端退出")
			break
		case "setdst":
			for k, v := range onlineMap {
				fmt.Println(k, v)
			}
			dst, err := reader.ReadString('\n')
			dst = strings.Trim(dst, " \r\n")
			if err != nil {
				fmt.Println("readString err=", err)
			}
			if onlineMap[dst] == true {
				DestinationAddr = dst
			}
			continue
		case "resetdst":
			DestinationAddr = ""
		case "":
			continue
		}
		cmd := base.Message{
			Type:       "cmd",
			CreateTime: time.Now(),
			ModifyTime: time.Now(),
			Src:        c.Addr,
			Dst:        DestinationAddr,
			Content:    line,
			Log:        nil,
		}
		//发送命令
		c.WriteChan <- cmd

	}
}

func (c *Controller) GetResult() {
	//接收结果
	res := <-c.ReadChan
	if res.Type == "onlineList" {
		if list, ok := res.Content.([]string); ok == true {
			for _, v := range list {
				if v != "" {
					onlineMap[v] = true
				}
			}
			fmt.Println("OnLine:", list)
		}
	} else {
		fmt.Println(res.Content)
	}
}
