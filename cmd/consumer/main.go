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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("start client")
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	p := &msg.Msg{Content: make([]byte, 8)}
	binary.LittleEndian.PutUint64(p.Content, uint64(time.Now().Unix()))
	p.CreatTime = uint64(time.Now().Unix())
	p.PID = os.Getpid()
	copy(p.Topic[:], []byte("yafngzh"))
	p.Length = binary.Size(p)
	p.Type = msg.MSG_TYPE_SUBSCRIBE
	err = encoder.Encode(p)
	if err != nil {
		log.Println("err %v", err)
		return
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
		return
	}
	log.Printf("receive %+v", q)
	for {
		qMsg := &msg.Msg{}
		err = decoder.Decode(qMsg)
		if err != nil {
			if err == io.EOF {
				log.Println("远端关闭连接 1")
			}
			log.Println(err)
			break
		}
		log.Printf("[receive %v] %v", conn.RemoteAddr(), qMsg)
		qResp := &msg.RespMsg{}
		err = encoder.Encode(qResp)
		if err != nil {
			if err == io.EOF {
				log.Println("远端关闭连接 2")
			}
			log.Println(err)
			break
		}
		time.Sleep(time.Second * 3)
	}
	log.Println("done")
}
