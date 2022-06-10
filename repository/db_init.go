package repository

import (
	"demo1/config"
	"demo1/model/entity"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() *gorm.DB {
	db = connectDB()
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Video{})
	db.AutoMigrate(&entity.Favorite{})
	db.AutoMigrate(&entity.Comment{})
	db.AutoMigrate(&entity.Relation{})

	// 初始化总的数量
	//db.Model(&User{}).Count(&UserCount)
	//db.Model(&Video{}).Count(&VideoCount)

	return db
}

// 缓存预编译语句
//
func connectDB() *gorm.DB {
	var err error
	dsn := config.DBConnectString()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Error connecting to database : error=%v", err))
	}

	return db
}
