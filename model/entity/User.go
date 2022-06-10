package entity

type User struct {
	ID            uint   `gorm:"primaryKey; not null" json:"id"`
	Name          string `gorm:"type:varchar(64); not null;" json:"name"`
	Token         string `gorm:"-" json:"token,omitempty"`
	Password      string `gorm:"char(24) ; not null;" json:"password,omitempty"`
	FollowCount   int64  `gorm:"not null; default:0" json:"follow_count"`   // 关注的人的数量
	FollowerCount int64  `gorm:"not null; default:0" json:"follower_count"` // 粉丝总数
	IsFollow      bool   `gorm:"-" json:"is_follow"`
}
