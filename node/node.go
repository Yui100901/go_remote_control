package node

import (
	"encoding/gob"
	"go_remote_control/message"
	"net"
)

type Node struct {
	Conn      net.Conn             //节点连接对象
	Addr      string               //节点地址
	ReadChan  chan message.Message //节点读取通道
	WriteChan chan message.Message //节点写入通道
	Enc       *gob.Encoder         //节点序列化器
	Dec       *gob.Decoder         //节点反序列化器
	NodeType  string               //节点类型
}

func (n *Node) Write() {
	for {
		msg := <-n.WriteChan
		if err := n.Enc.Encode(msg); err != nil {
			return
		}
	}
}

func (n *Node) Read() {
	for {
		msg := message.Message{}
		if err := n.Dec.Decode(&msg); err != nil {
			return
		}
		n.ReadChan <- msg
	}
}
