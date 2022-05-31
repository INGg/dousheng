package repository

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() *gorm.DB {
	db = connectDB()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Video{})
	//db.AutoMigrate(Comment{})

	// 初始化总的数量
	db.Model(&User{}).Count(&UserCount)
	db.Model(&Video{}).Count(&VideoCount)

	return db
}

func connectDB() *gorm.DB {
	var err error
	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Error connecting to database : error=%v", err))
	}

	return db
}
