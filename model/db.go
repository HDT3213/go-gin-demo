package model

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/go-gin-demo/entity"
)

var DB *gorm.DB

func setupDB() {
    var err error
    DB, err = gorm.Open("mysql", "go:password@tcp(127.0.0.1:3306)/go?charset=utf8&parseTime=True&loc=Asia%2fShanghai")
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

func Setup() {
    setupDB()
    setupRedis()
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