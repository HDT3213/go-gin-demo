package mq

type MQ interface {
    Close()
    Publish(msg *Msg) error
    Consume(consumerMap map[uint32]func(*Msg)) error
}
