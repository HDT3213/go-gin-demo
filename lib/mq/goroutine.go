package mq

import "github.com/pkg/errors"

type GoRoutineMQ struct {
    ConsumerMap map[uint32]ConsumerFunc
}

func SetupGoRoutineMQ(consumerMap map[uint32]ConsumerFunc) (*GoRoutineMQ, error) {
    return &GoRoutineMQ{
        ConsumerMap: consumerMap,
    }, nil
}

func (mq *GoRoutineMQ)Close() {
}

func (mq *GoRoutineMQ)Publish(msg *Msg) error {
    consumer, ok := mq.ConsumerMap[msg.Code]
    if !ok {
        return errors.New("consumer not found")
    }
    go func() {
        consumer(msg)
    }()
    return nil
}

func (mq *GoRoutineMQ)Consume(consumerMap map[uint32]ConsumerFunc) error {
    return errors.New("operation not supported")
}