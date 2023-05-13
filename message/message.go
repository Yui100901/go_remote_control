package message

import (
	"time"
)

// Message 自定义消息格式
type Message struct {
	//Type:sync,cmd,res
	Type       string
	CreateTime time.Time
	ModifyTime time.Time
	Src        string
	Dst        string
	Content    any
	Log        []string
}
