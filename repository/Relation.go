package repository

import (
	"demo1/model/entity"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

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

//	AddRelation 将 用户id 和被其关注人的id 插入表中 relation
func (r *RelationDao) AddRelation(FollowerId uint, AuthorId uint) error {
	res := db.Create(&entity.Relation{
		UserID:   AuthorId,
		FollowID: FollowerId,
	})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		zap.L().Error("Add Relation error")
		return res.Error
	}
	zap.L().Info("insert relation success")
	return nil
}

// DeleteRelation 根据 userid followerId 删除对应记录
func (r *RelationDao) DeleteRelation(FollowerId uint, AuthorId uint) error {
	res := db.Where(&entity.Relation{
		UserID:   AuthorId,
		FollowID: FollowerId,
	}).Delete(&entity.Relation{})

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		zap.L().Error("Delete relation error")
		return res.Error
	}
	zap.L().Info("Delete relation success")
	return nil
}

// QueryFollowIdByUserID 查询当前用户的关注列表(id)
func (r *RelationDao) QueryFollowIdByUserID(uid uint, RelationList *[]entity.Relation) error {
	res := db.Model(&entity.Relation{}).Where("user_id = ?", uid).Find(RelationList)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// QueryUsersIDByFollowId 查询当前用户的粉丝(id)
func (r *RelationDao) QueryUsersIDByFollowId(FollowerId uint, relationList *[]entity.Relation) error {
	res := db.Model(&entity.Relation{}).Where("follow_id = ?", FollowerId).Find(relationList)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *RelationDao) QueryAFollowB(Auid uint, Buid uint) bool {
	res := db.Model(&entity.Relation{}).Where("user_id = ?", Auid).Where("follow_id = ?", Buid)
	if res.Error != nil {
		return false
	}
	return true
}
