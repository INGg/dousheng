package repository

import (
	"demo1/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() *gorm.DB {
	db = connectDB()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Video{})
	db.AutoMigrate(&Favorite{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&Relation{})

	// 初始化总的数量
	//db.Model(&User{}).Count(&UserCount)
	//db.Model(&Video{}).Count(&VideoCount)

	return db
}

func connectDB() *gorm.DB {
	var err error
	dsn := config.DBConnectString()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Error connecting to database : error=%v", err))
	}

	return db
}
