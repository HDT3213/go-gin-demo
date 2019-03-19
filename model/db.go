package model

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/go-gin-demo/entity"
)

var db *gorm.DB

func setupDB() {
    var err error
    db, err = gorm.Open("mysql", "root:248536@tcp(127.0.0.1:3306)/go?charset=utf8&parseTime=True&loc=Asia%2fShanghai")
    if nil != err {
        panic(err)
    }
    if !db.HasTable(&entity.User{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&entity.User{}).Error; err != nil {
            panic(err)
        }
    }
    db.LogMode(true)
}

func Setup() {
    setupDB()
    setupRedis()
}

func closeDB() {
    if db != nil {
        db.Close()
    }
}

func Close() {
    closeDB()
    closeCache()
}