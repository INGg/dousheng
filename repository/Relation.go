package repository

import "sync"

type Relation struct {
	ID       uint `gorm:"primaryKey; not null; auto_increment" json:"id"`
	UserId   uint `gorm:"not null" json:"user_id"`
	FollowId uint `gorm:"not null" json:"follow_id"`
	status   uint `gorm:"not null;default:0"`
}

type RelationDao struct {
}

var (
	relationDao  *RelationDao
	relationOnce sync.Once
)

func NewRelationDAO() *RelationDao {
	relationOnce.Do(func() {
		relationDao = &RelationDao{}
	})
	return relationDao
}
func (r *RelationDao) AddRelation(UserId uint, FollowId uint) error {
	//	将 用户id 和被其关注人的id 插入表中 relation表中,status 字段设置为1.
	return nil
}

func (r *RelationDao) QueryFollowIdByUserId(UserId uint) []uint {

	//	查询当前用户所关注的人
	return nil
}

func (r *RelationDao) QueryUserIdByFollowId(FollowId uint) []uint {
	// 查询当前用户的粉丝
	return nil
}
