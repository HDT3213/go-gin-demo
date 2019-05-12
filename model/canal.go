package model

import (
    "github.com/siddontang/go-mysql/canal"
    "github.com/go-gin-demo/utils/logger"
)

type CanalSettings struct {
    Addr string `yaml:"addr"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    DB string `yaml:"db"`
    Tables []string `yaml:"tables"`
}

var canalClient *canal.Canal

type EventHandler struct {
    canal.DummyEventHandler
}

func (h *EventHandler) OnRow(e *canal.RowsEvent) error {
    logger.Info("%s %v\n", e.Action, e.Rows)
    return nil
}

func (h *EventHandler) String() string {
    return "EventHandler"
}

func SetupCanal(settings *CanalSettings) {
    cfg := canal.NewDefaultConfig()
    cfg.Addr = settings.Addr
    cfg.User = settings.Username
    cfg.Password = settings.Password
    cfg.Dump.TableDB = settings.DB
    cfg.Dump.Tables = settings.Tables
    var err error
    canalClient, err = canal.NewCanal(cfg)
    if err != nil {
        panic(err)
    }
}

func Listen() {
    canalClient.SetEventHandler(&EventHandler{})
    canalClient.Run()
}

func CloseCanal() {
    canalClient.Close()
}