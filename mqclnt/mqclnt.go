package mqclnt

import "net"

type Mqclnt struct {
	CreatTime uint64
	IP        string
	Port      int
	PID       int
	Conn      net.Conn
}

func (clnt *Mqclnt) Equal(clnt_another *Mqclnt) bool {
	if clnt.IP == clnt_another.IP && clnt.Port == clnt_another.Port && clnt.CreatTime == clnt_another.CreatTime {
		return true
	}
	return false
}
