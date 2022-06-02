package repository

type Relation struct {
	ID       uint `gorm:"primaryKey; not null; auto_increment" json:"id"`
	UserId   uint `gorm:"not null" json:"user_id"`
	FollowId uint `gorm:"not null" json:"follow_id"`
	status   uint `gorm:"not null;default:0"`
}
