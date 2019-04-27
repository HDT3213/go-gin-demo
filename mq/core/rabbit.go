package core

import (
    "github.com/streadway/amqp"
    "encoding/binary"
    "fmt"
    "github.com/go-gin-demo/utils/logger"
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

type Settings struct {
    Host string `yaml:"host"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    VHost string `yaml:"vhost"`
    Exchange string `yaml:"exchange"`
    Queue string `yaml:"queue"`
    RoutingKey string `yaml:"routing-key"`
}

var settings *Settings

func SetupRabbitMQ(rabbitSettings *Settings) error {
    settings = rabbitSettings
    var err error
    MQ, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s",
        settings.Username,
        settings.Password,
        settings.Host,
        settings.VHost))
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
        settings.Exchange,
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
        settings.Queue,
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
        settings.Queue,
        settings.RoutingKey,
        settings.Exchange,
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
        settings.Exchange,
        settings.RoutingKey,
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
       settings.Queue,
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
                logger.Warn(fmt.Sprintf("no consumer found for: %d", code))
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
