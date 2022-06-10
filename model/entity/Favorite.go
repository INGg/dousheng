package entity

import "gorm.io/gorm"

type Favorite struct {
	ID      uint  `gorm:"primaryKey; not null; auto_increment" json:"id"`
	UserID  uint  `gorm:"not null" json:"user_id"`
	VideoID uint  `gorm:"not null" json:"video_id"`
	Video   Video `gorm:"foreignKey:VideoID"`
	Deleted gorm.DeletedAt
}
