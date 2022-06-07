package entity

type Video struct {
	ID            uint   `gorm:"primaryKey; not null" json:"id"`
	PublishTime   int64  `gorm:"index; timestamp" json:"publish_time"`
	Author        User   `gorm:"foreignKey:AuthorID" json:"author"`
	AuthorID      uint   //`json:"author_id"`
	PlayUrl       string `gorm:"type:varchar(128)" json:"play_url"`  // 播放地址
	CoverUrl      string `gorm:"type:varchar(128)" json:"cover_url"` // 封面地址
	FavoriteCount int64  `gorm:"not_null; default:0" json:"favorite_count"`
	CommentCount  int64  `gorm:"not_null; default:0" json:"comment_count"`
	Title         string ` json:"title"`
}
