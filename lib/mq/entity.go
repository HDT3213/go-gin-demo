package mq

type Msg struct {
    Code    uint32 // for type
    Payload []byte
}