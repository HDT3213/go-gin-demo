package model

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/go-gin-demo/entity"
    "fmt"
)

type DBSettings struct {
    Dialect string `yaml:"dialect"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Host string `yaml:"host"`
    DB string `yaml:"db"`
    Charset string `yaml:"charset"`
    Location string `yaml:"location"`
}

var DB *gorm.DB

func setupDB(settings *DBSettings) {
    var err error
    DB, err = gorm.Open(settings.Dialect, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=%s",
        settings.Username,
        settings.Password,
        settings.Host,
        settings.DB,
        settings.Charset,
        settings.Location))

    if nil != err {
        panic(err)
    }
    if !DB.HasTable(&entity.User{}) {
        if err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&entity.User{}).Error; err != nil {
            panic(err)
        }
        DB.Model(&entity.User{}).AddIndex("idx_username", "username")
    }
    if !DB.HasTable(&entity.Post{}) {
        if err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&entity.Post{}).Error; err != nil {
            panic(err)
        }
        DB.Model(&entity.Post{}).AddIndex("idx_user", "uid", "valid", "created_at")
    }
    if !DB.HasTable(&entity.Follow{}) {
        if err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&entity.Follow{}).Error; err != nil {
            panic(err)
        }
        DB.Model(&entity.Follow{}).AddIndex("idx_user", "uid", "valid", "created_at")
    }
    DB.LogMode(true)
}

func Setup(dbSettings *DBSettings, redisSettings *RedisSettings) {
    setupDB(dbSettings)
    setupRedis(redisSettings)
}

func closeDB() {
    if DB != nil {
        DB.Close()
    }
}

func Close() {
    closeDB()
    closeCache()
}