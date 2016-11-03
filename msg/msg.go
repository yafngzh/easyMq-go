package msg

const TYPE_LEN = 16
const TOPIC_LEN = 16

const (
	RESP_OK            = "ok"
	RESP_ERROR         = "error"
	MSG_TYPE_SUBSCRIBE = 1
	MSG_TYPE_PUBLISH   = 2
)

type Msg struct {
	Length    int //消息包总长度
	PID       int
	CreatTime uint64 //创建的纳秒数

	Type    uint32
	Topic   [TOPIC_LEN]byte
	Content []byte
}

type RespMsg struct {
	Resp string
}
