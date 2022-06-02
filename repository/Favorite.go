package repository

type Favorite struct {
	ID      uint `gorm:"primaryKey; not null; auto_increment" json:"id"`
	UserId  uint `gorm:"not null" json:"user_id"`
	VideoId uint `gorm:"not null" json:"video_id"`
	status  uint `gorm:"not null;default:0" json:"status"`
}
