package canal


import (
    "github.com/siddontang/go-mysql/canal"
    "github.com/go-gin-demo/lib/logger"
)

type ConsumerFunc func(formerRow interface{}, row interface{})error

type Settings struct {
    Addr string `yaml:"addr"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    DB string `yaml:"db"`
    Tables []string `yaml:"tables"`
}

type Canal struct {
    Client *canal.Canal
    Settings *Settings
}

type EventHandler struct {
    canal.DummyEventHandler
    ConsumerMap map[string]ConsumerFunc
}

func (h *EventHandler) OnRow(e *canal.RowsEvent) error {
    logger.Info("%s %v\n", e.Action, e.Rows)
    tableName := e.Table.Schema + "." + e.Table.Name
    consumer, ok := h.ConsumerMap[tableName]
    if !ok {
        return nil
    }
    return consumer(e.Rows[0], e.Rows[1])
}

func (h *EventHandler) String() string {
    return "EventHandler"
}

func SetupCanal(settings *Settings) (*Canal, error) {
    cfg := canal.NewDefaultConfig()
    cfg.Addr = settings.Addr
    cfg.User = settings.Username
    cfg.Password = settings.Password
    cfg.Dump.TableDB = settings.DB
    cfg.Dump.Tables = settings.Tables
    client, err := canal.NewCanal(cfg)
    if err != nil {
        return nil, err
    }
    return &Canal{
        Client: client,
        Settings: settings,
    }, nil
}

func (c *Canal)Listen(ConsumerMap map[string]ConsumerFunc) {
    c.Client.SetEventHandler(&EventHandler{})
    c.Client.Run()
}