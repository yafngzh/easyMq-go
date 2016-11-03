package mqserver

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/yafngzh/easyMq-go/mqclnt"
	"github.com/yafngzh/easyMq-go/msg"
)

//Mq 统一接口
//订阅，转发
type Mq interface {
}

//Mqserver 处理主要消息的功能
type Mqserver struct {
	Lock       sync.RWMutex
	Topic2clnt map[[msg.TOPIC_LEN]byte][]mqclnt.Mqclnt
}

func NewMqserver() *Mqserver {
	s := &Mqserver{}
	s.Topic2clnt = make(map[[msg.TOPIC_LEN]byte][]mqclnt.Mqclnt)
	return s
}

func (s *Mqserver) subscribe(topic [msg.TOPIC_LEN]byte, clnt mqclnt.Mqclnt) error {
	var (
		clntSlice []mqclnt.Mqclnt
		exists    bool
	)
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if clntSlice, exists = s.Topic2clnt[topic]; exists != true {
		//添加对象
		clntSlice = append(clntSlice, clnt)
		s.Topic2clnt[topic] = clntSlice
		log.Printf("add clnt %s %d", clnt.IP, clnt.Port)
	}
	//查找对应的clnt
	for i := range clntSlice {
		if clntSlice[i] == clnt {
			//do nothing
		} else {
			//
			clntSlice = append(clntSlice, clnt)
			log.Printf("add clnt %s %d", clnt.IP, clnt.Port)
		}
	}
	return nil
}

func (s *Mqserver) HandleMsg(m *msg.Msg, clnt mqclnt.Mqclnt) (*msg.RespMsg, error) {
	var (
		err  error
		resp = &msg.RespMsg{}
	)
	log.Printf("receive Msg %+v, clnt %+v\n", m, clnt)
	if m.Type == 1 {
		err = s.subscribe(m.Topic, clnt)
	} else if m.Type == 2 {
		err = s.transfer(m.Topic, m.Content, clnt)
	} else {
		err = fmt.Errorf("Error type %+v", m)
	}

	if err != nil {
		resp.Resp = msg.RESP_ERROR
	} else {
		resp.Resp = msg.RESP_OK
	}

	return resp, err
}
func (s *Mqserver) transfer(topic [msg.TOPIC_LEN]byte, msgContent []byte, clnt mqclnt.Mqclnt) error {
	var (
		exists    bool
		clntSlice []mqclnt.Mqclnt
		req       = &msg.Msg{}
	)

	s.Lock.Lock()
	defer s.Lock.Unlock()
	if len(msgContent) == 0 {
		return fmt.Errorf("msgContent is null")
	}
	if clntSlice, exists = s.Topic2clnt[topic]; exists != true {
		return nil
	}
	//查找对应的clnt
	for i := range clntSlice {
		encoder := gob.NewEncoder(clntSlice[i].Conn)
		log.Printf("map is %+v, conn is %+v\n", s.Topic2clnt, clntSlice[i].Conn)
		req.Content = msgContent[:]
		req.Type = msg.MSG_TYPE_PUBLISH
		err := encoder.Encode(req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[write %v] bytes %v", clntSlice[i].Conn.RemoteAddr(), err)

		decoder := gob.NewDecoder(clntSlice[i].Conn)
		q := &msg.RespMsg{}
		err = decoder.Decode(q)
		if err != nil {
			if err == io.EOF {
				log.Println("远端关闭连接")
			}
			log.Printf("receive consumer msg resp err %v", err)
			break
		}
		log.Printf("[receive %v] %+v", clntSlice[i].Conn.RemoteAddr(), q)

	}
	return nil
}
