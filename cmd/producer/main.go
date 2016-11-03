package main

import (
	"encoding/binary"
	"encoding/gob"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/yafngzh/easyMq-go/msg"
)

func main() {
	log.Println("start client")
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close()
	for {
		encoder := gob.NewEncoder(conn)
		p := &msg.Msg{Content: make([]byte, 8)}
		binary.LittleEndian.PutUint64(p.Content, uint64(time.Now().Unix()))
		p.CreatTime = uint64(time.Now().Unix())
		p.PID = os.Getpid()
		copy(p.Topic[:], []byte("yafngzh"))
		p.Type = msg.MSG_TYPE_PUBLISH
		p.Content = []byte{0x1, 0x2, 0x3}
		p.Length = binary.Size(p)
		err := encoder.Encode(p)
		if err != nil {
			log.Println("err %v", err)
			break
		}
		log.Println("send %+v", p)

		decoder := gob.NewDecoder(conn)
		q := &msg.RespMsg{}
		err = decoder.Decode(q)
		if err != nil {
			if err == io.EOF {
				log.Println("远端关闭连接")
			}
			log.Println(err)
			break
		}
		log.Printf("receive %+v", q)
		time.Sleep(time.Second * 3)
		break
	}
	log.Println("done")
}
