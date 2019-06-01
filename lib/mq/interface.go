package mq

type ConsumerFunc func(*Msg)

type MQ interface {
    Close()
    Publish(msg *Msg) error
    Consume(consumerMap map[uint32]ConsumerFunc) error
}
