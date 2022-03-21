package queue

//IQueue 消息队列
type IQueue interface {
	Send(key string, value Message) error
	Pop(key string) (string, error)
	Count(key string) (int64, error)
}

//IMQCMessage  队列消息
type IMQCMessage interface {
	Ack() error
	Nack() error
	GetMessage() Message
}

type Message interface {
	Header() map[string]string
	Body() map[string]string
}

type ConsumeCallback func(IMQCMessage)

//IMQC consumer接口
type IMQC interface {
	Connect() error
	Consume(queue string, callback ConsumeCallback) (err error)
	Unconsume(queue string)
	Close()
}

//IMQP 消息生产
type IMQP interface {
	Push(key string, value Message) error
	Pop(key string) (string, error)
	Count(key string) (int64, error)
	Close() error
}

//IComponentQueue Component Queue
type IComponentQueue interface {
	GetQueue(name string) (q IQueue)
}