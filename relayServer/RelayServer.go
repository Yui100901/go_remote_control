package main

import (
	"encoding/gob"
	"fmt"
	"go_remote_control/base"
	"log"
	"net"
	"sync"
	"time"
)

const (
	ControllerPort    = "9900" //控制端连接端口
	CommandServerPort = "9901" //被控制端连接端口
)

var (
	controllerMap    = sync.Map{}
	commandServerMap = sync.Map{}
)

type OnlineNode struct {
	base.Node
}

func (on *OnlineNode) handle() {
	for {
		msg := <-on.ReadChan
		switch msg.Type {
		case "cmd":
			if msg.Dst != "" {
				dstNode, _ := commandServerMap.Load(msg.Dst)
				if dn, ok := dstNode.(OnlineNode); ok == true {
					dn.WriteChan <- msg
				}
			} else {
				var content []string
				commandServerMap.Range(func(k, v any) bool {
					if n, ok := v.(OnlineNode); ok == true {
						content = append(content, n.Addr)
					} else {
						fmt.Print("false")
					}
					return true
				})
				res := base.Message{
					Type:       "onlineList",
					CreateTime: time.Now(),
					ModifyTime: time.Now(),
					Src:        msg.Dst,
					Dst:        msg.Src,
					Content:    content,
					Log:        nil,
				}
				on.WriteChan <- res
				log.Println(res)
			}
		case "res":
			dstNode, _ := controllerMap.Load(msg.Dst)
			if dn, ok := dstNode.(OnlineNode); ok == true {
				dn.WriteChan <- msg
			}
		case "alive":
			//log.Println(msg.Content)
		}
	}
}

func main() {

	log.Println("服务器开始监听....")
	//tcp协议6666端口监听

	go startListen(ControllerPort)
	go startListen(CommandServerPort)
	select {}

}

func startListen(port string) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Println("监听失败 err=", err)
		return
	}
	defer listener.Close()
	//循环等待客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept() err=", err)
			continue
		}
		log.Printf("Accept() suc con=%v 客户端ip=%v\n", conn, conn.RemoteAddr().String())
		onlineNode := OnlineNode{
			Node: base.Node{
				Conn:      conn,
				Addr:      conn.RemoteAddr().String(),
				ReadChan:  make(chan base.Message),
				WriteChan: make(chan base.Message),
				Enc:       gob.NewEncoder(conn),
				Dec:       gob.NewDecoder(conn),
				NodeType:  "",
			},
		}
		fmt.Println(port)
		switch port {
		case "9900":
			controllerMap.Store(onlineNode.Addr, onlineNode)
		case "9901":
			commandServerMap.Store(onlineNode.Addr, onlineNode)
		}
		go onlineNode.Read()
		go onlineNode.Write()
		go onlineNode.handle()
	}
}