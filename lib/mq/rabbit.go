package mq

import (
    "github.com/streadway/amqp"
    "encoding/binary"
    "fmt"
    "github.com/go-gin-demo/lib/logger"
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

type RabbitMQ struct {
    Conn *amqp.Connection
    Settings *Settings
}

func SetupRabbitMQ(settings *Settings) (*RabbitMQ, error) {
    var err error
    conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s",
        settings.Username,
        settings.Password,
        settings.Host,
        settings.VHost))
    if err != nil {
        return nil, err
    }
    ch, err := conn.Channel()
    defer ch.Close()
    if err != nil {
        return nil, err
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
        return nil, err
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
        return nil, err
    }

    err = ch.QueueBind(
        settings.Queue,
        settings.RoutingKey,
        settings.Exchange,
        false,
        nil,
    )
    if err != nil {
        return nil, err
    }
    return &RabbitMQ{
        Conn: conn,
        Settings: settings,
    }, nil
}

func (mq *RabbitMQ)Close() {
    mq.Conn.Close()
}


func (mq *RabbitMQ)Publish(msg *Msg) error {
    ch, err := mq.Conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()

    bytes := make([]byte, 4 + len(msg.Payload))
    binary.LittleEndian.PutUint32(bytes, msg.Code)
    copy(bytes[4:], msg.Payload)

    err = ch.Publish(
        mq.Settings.Exchange,
        mq.Settings.RoutingKey,
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

func (mq *RabbitMQ)Consume(consumerMap map[uint32]ConsumerFunc) error {
    ch, err := mq.Conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()

    deliveries, err := ch.Consume(
       mq.Settings.Queue,
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
