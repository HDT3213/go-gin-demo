package core

import (
    "github.com/streadway/amqp"
    "encoding/binary"
    "log"
    "fmt"
)

var MQ *amqp.Connection

const (
    username = "go"
    password = "password"
    url =  "localhost:5672"
    vhost = "go"
    exchange = "go"
    queue = "general"
    routingKey = "general"
)

func SetupRabbitMQ() error {
    var err error
    MQ, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, url, vhost))
    if err != nil {
        return err
    }
    ch, err := MQ.Channel()
    defer ch.Close()
    if err != nil {
        return err
    }

    // prepare exchange
    err = ch.ExchangeDeclare(
        exchange,
        "direct",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    // prepare queue
    _, err = ch.QueueDeclare(
        queue,
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    err = ch.QueueBind(
        queue,
        routingKey,
        exchange,
        false,
        nil,
    )
    if err != nil {
        return err
    }
    return nil
}

func CloseRabbitMQ() {
    MQ.Close()
}

func GetChannel() *amqp.Channel {
    ch, err := MQ.Channel()
    if err != nil {
        panic(err)
    }
    return ch
}

func CloseChannel(ch *amqp.Channel) {
    ch.Close()
}

func Publish(msg *Msg) error {
    ch := GetChannel()
    defer CloseChannel(ch)

    bytes := make([]byte, 4 + len(msg.Payload))
    binary.LittleEndian.PutUint32(bytes, msg.Code)
    copy(bytes[4:], msg.Payload)

    err := ch.Publish(
        exchange,
        routingKey,
        false,
        false,
        amqp.Publishing {
            DeliveryMode: amqp.Persistent,
            ContentType:  "application/octet-stream",
            Body:         []byte(bytes),
        },
    )
    return err
}

func Consume(consumerMap map[uint32]func(*Msg)) error {
    ch := GetChannel()
    defer CloseChannel(ch)

    deliveries, err := ch.Consume(
       queue,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    forever := make(chan bool)

    go func() {
        for delivery := range deliveries {
            codeBytes := delivery.Body[:4]
            code := binary.LittleEndian.Uint32(codeBytes)
            consumer, ok := consumerMap[code]
            if !ok {
                log.Println(fmt.Sprintf("no consumer found for: %d", code))
                delivery.Ack(false)
                continue
            }
            consumer(&Msg{
                Code: code,
                Payload: delivery.Body[4:],
            })
            delivery.Ack(false)
        }
        forever <- true  // delivery has
    }()
    <-forever
    return nil
}
