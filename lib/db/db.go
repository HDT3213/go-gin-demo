package db

import (
    "github.com/jinzhu/gorm"
    "fmt"
)

type Settings struct {
    Dialect string `yaml:"dialect"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Host string `yaml:"host"`
    DB string `yaml:"db"`
    Charset string `yaml:"charset"`
    Location string `yaml:"location"`
}


func Setup(settings *Settings) (*gorm.DB, error) {
    var err error
    db, err := gorm.Open(settings.Dialect, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=%s",
        settings.Username,
        settings.Password,
        settings.Host,
        settings.DB,
        settings.Charset,
        settings.Location))

    if nil != err {
        panic(err)
    }
    db.LogMode(true)
    return db, nil
}
