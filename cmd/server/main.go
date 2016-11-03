package main

import (
	"encoding/gob"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/yafngzh/easyMq-go/mqclnt"
	"github.com/yafngzh/easyMq-go/mqserver"
	"github.com/yafngzh/easyMq-go/msg"
)

var (
	svr *mqserver.Mqserver
)

func handleConn(c net.Conn) {
	var exitClose = true
	defer func() {
		if exitClose {
			log.Printf("c exit,ip %v, port %+v", c.RemoteAddr)
			c.Close()
		}
	}()
	for {
		var (
			clnt mqclnt.Mqclnt
		)
		t1 := time.Now().UnixNano()
		dec := gob.NewDecoder(c)
		curMsg := &msg.Msg{}
		err := dec.Decode(curMsg)
		if err != nil {
			if err == io.EOF {
				log.Printf("[read %v] 对方关闭连接\n", c.RemoteAddr())
			} else {
				log.Printf("read error %v from %v, now exit\n", err, c.RemoteAddr())
			}
			break
		}
		//log.Printf("receive msg %+v\n", curMsg)
		clnt.CreatTime = curMsg.CreatTime
		addr := c.RemoteAddr().String()
		addrSlice := strings.Split(addr, ":")
		if len(addrSlice) != 2 {
			log.Printf("格式错误 %s \n", addr)
			break
		}
		clnt.IP = addrSlice[0]
		clnt.Port, err = strconv.Atoi(addrSlice[1])
		if err != nil {
			log.Printf("%v\n", err)
			break
		}
		clnt.PID = curMsg.PID
		clnt.Conn = c

		resp, err := svr.HandleMsg(curMsg, clnt)
		if err != nil {
			log.Fatalf("handle failed %v", err)
			break
		}
		encoder := gob.NewEncoder(c)
		err = encoder.Encode(resp)
		if err != nil {
			log.Println("send error %v", err)
		}

		t2 := time.Now().UnixNano()
		log.Printf("time elapsed %d ms\n", (t2-t1)/1000)
		if curMsg.Type == msg.MSG_TYPE_SUBSCRIBE {
			exitClose = false
			break
		}
	}
}
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	svr = mqserver.NewMqserver()
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Printf("Listen error %v", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Accept error, %v", err)
			continue
		}
		go handleConn(c)
	}
}
