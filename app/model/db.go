package model

import (
    "github.com/HDT3213/go-gin-demo/app/entity"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

func Migrate(db *gorm.DB) error {
    if !db.HasTable(&entity.User{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&entity.User{}).Error; err != nil {
           return err
        }
        db.Model(&entity.User{}).AddIndex("idx_username", "username")
    }
    if !db.HasTable(&entity.Post{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&entity.Post{}).Error; err != nil {
            return err
        }
        db.Model(&entity.Post{}).AddIndex("idx_user", "uid", "valid", "created_at")
    }
    if !db.HasTable(&entity.Follow{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&entity.Follow{}).Error; err != nil {
            return err
        }
        db.Model(&entity.Follow{}).AddIndex("idx_user", "uid", "valid", "created_at")
    }
    return nil
}
