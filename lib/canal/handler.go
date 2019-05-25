package canal

import (
    "github.com/siddontang/go-mysql/canal"
    "github.com/go-gin-demo/lib/logger"
    "fmt"
    "github.com/siddontang/go-mysql/mysql"
    "github.com/siddontang/go-mysql/schema"
    "runtime/debug"
)

type Row map[string]interface{}
type ConsumerFunc func(formerRow Row, row Row)error

type EventHandler struct {
    canal.DummyEventHandler
    ConsumerMap map[string]ConsumerFunc
    CanalCli *Canal
}

func MakeEventHandler(consumerMap map[string]ConsumerFunc, canalCli *Canal) *EventHandler {
    return &EventHandler{
        ConsumerMap: consumerMap,
        CanalCli: canalCli,
    }
}

func (h *EventHandler) String() string {
    return "EventHandler"
}

func toMap(table *schema.Table, row []interface{}) map[string]interface{} {
    m := make(map[string]interface{})
    for i, column := range table.Columns {
        name := column.Name
        val := row[i]
        if column.RawType == "tinyint(1)" {
            // treat tinyint(1) as bool
            m[name] = val != 0
        } else if column.RawType == "text" {
            // canal treat `varchar` as `string`, treat `text` as `uint8[]`
            // convert uint8[] to string
            uint8slice, ok := val.([]uint8)
            if ok {
                bytes := make([]byte, len(uint8slice))
                for j, b := range uint8slice {
                    bytes[j] = byte(b)
                }
                m[name] = string(bytes)
            } else {
                m[name] = val
            }
        } else {
            m[name] = val
        }
    }
    return m
}

func (h *EventHandler) OnRow(e *canal.RowsEvent) error {
    defer func() {
        if err := recover(); err != nil {
            logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
        }
    }()

    key := e.Table.Schema + "." + e.Table.Name + ":" + e.Action
    consumer, ok := h.ConsumerMap[key]
    if !ok {
        return nil
    }
    var before, after Row
    if e.Action == canal.InsertAction {
        before = nil
        after = toMap(e.Table, e.Rows[0])
    } else if e.Action == canal.UpdateAction {
        before = toMap(e.Table, e.Rows[0])
        after = toMap(e.Table, e.Rows[1])
    } else {
        return nil
    }
    err :=  consumer(before, after)
    if err != nil {
        logger.Warn(fmt.Sprintf("consume failed, event: %v, err: %v\n%s", e, err, string(debug.Stack())))
    }
    return nil
}

func (h *EventHandler) OnPosSynced(pos mysql.Position, force bool) error {
    h.CanalCli.SavePos(&pos)
    return nil
}