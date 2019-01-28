package model

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func Setup() {
    var err error
    db, err = gorm.Open("mysql", "root:248536@tcp(127.0.0.1:3306)/go?charset=utf8&parseTime=True&loc=Asia%2fShanghai")
    if nil != err {
        panic(err)
    }
    if !db.HasTable(&User{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
            panic(err)
        }
    }
}

func Close() {
    if db != nil {
        db.Close()
    }
}