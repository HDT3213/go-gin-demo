package canal


import (
    "fmt"
    "github.com/HDT3213/go-gin-demo/lib/logger"
    "github.com/HDT3213/go-gin-demo/lib/zookeeper"
    "github.com/samuel/go-zookeeper/zk"
    "github.com/siddontang/go-mysql/canal"
    "github.com/siddontang/go-mysql/mysql"
    "github.com/vmihailenco/msgpack"
)

type Settings struct {
    Addr string `yaml:"addr"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    DB string `yaml:"db"`
    Tables []string `yaml:"tables"`
    Zookeeper zookeeper.Settings
    ZkPath string `yaml:"path"`
}

type Canal struct {
    Client *canal.Canal
    Settings *Settings
    ZkCli *zk.Conn
}

func Setup(settings *Settings) (*Canal, error) {
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

    zk, err := zookeeper.Setup(&settings.Zookeeper)
    if err != nil {
        return nil, err
    }

    return &Canal{
        Client: client,
        Settings: settings,
        ZkCli: zk,
    }, nil
}

func (c *Canal)Listen(consumerMap map[string]ConsumerFunc) error {
    pos, err := c.GetPos()
    if err != nil {
        return err
    }
    if pos == nil {
        position := c.Client.SyncedPosition()
        pos = &position
    }
    c.Client.SetEventHandler(MakeEventHandler(consumerMap, c))
    c.Client.RunFrom(*pos)
    return nil
}

var (
    acls = zk.WorldACL(zk.PermAll)
)

func (c *Canal)SavePos(pos *mysql.Position) error {
    posBytes, err := msgpack.Marshal(pos)
    if err != nil {
        logger.Warn(fmt.Sprintf("marshal failed, err: %v", err))
    }
    exists, stat, err := c.ZkCli.Exists(c.Settings.ZkPath)
    if err != nil {
        logger.Warn(fmt.Sprintf("zk err: %v", err))
        return err
    }
    if exists {
        stat, err = c.ZkCli.Set(c.Settings.ZkPath, posBytes, stat.Version)
        if err != nil {
            logger.Warn(fmt.Sprintf("zk err: %v", err))
            return err
        }
    } else {
        _, err := c.ZkCli.Create(c.Settings.ZkPath, posBytes, 0, acls)
        if err != nil {
            logger.Warn(fmt.Sprintf("zk err: %v", err))
            return err
        }
    }
    return nil
}

func (c *Canal)GetPos() (*mysql.Position, error) {
    exists, _, err := c.ZkCli.Exists(c.Settings.ZkPath)
    if err != nil {
        logger.Warn(fmt.Sprintf("zk err: %v", err))
        return nil, err
    }
    if !exists {
        return nil, nil
    }
    data, _, err := c.ZkCli.Get(c.Settings.ZkPath)
    if err != nil {
        logger.Warn(fmt.Sprintf("zk err: %v", err))
        return nil, err
    }
    var pos mysql.Position
    err = msgpack.Unmarshal(data, &pos)
    if err != nil {
        logger.Warn(fmt.Sprintf("unmarshal err: %v", err))
        return nil, err
    }
    return &pos, nil
}

func (c *Canal)Close() {
    c.Client.Close()
    c.ZkCli.Close()
}