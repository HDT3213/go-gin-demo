package zookeeper

import (
    "github.com/samuel/go-zookeeper/zk"
    "time"
)

type Settings struct {
    Servers []string `yaml:"hosts"`
    SessionTimeout time.Duration `yaml:"session-timeout"`
}

func Setup(settings *Settings) (*zk.Conn, error) {
    conn, _, err := zk.Connect(settings.Servers, settings.SessionTimeout)
    if err != nil {
        return nil, err
    }
    return conn, nil
}