package entity

import "gorm.io/gorm"

type Comment struct {
	ID                 uint   `gorm:"primaryKey; not null" json:"comment_id"`
	CommentPublishTime int64  `gorm:"timestamp"`
	Content            string `gorm:"varchar(256)" json:"content"`
	AuthorID           uint   `gorm:"not null;"`
	Video              Video  `gorm:"foreignKey:VideoID"`
	VideoID            uint   `gorm:"index"`
	DeletedAt          gorm.DeletedAt
}
