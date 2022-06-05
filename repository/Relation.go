package repository

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
)

type Relation struct {
	ID         uint `gorm:"primaryKey; not null; auto_increment" json:"id"`
	AuthorId   uint `gorm:"not null" json:"user_id"`
	FollowerId uint `gorm:"not null" json:"follow_id"`
	DeleteAt   gorm.DeletedAt
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
func (r *RelationDao) AddRelation(FollowerId uint, AuthorId uint) error {
	//	将 用户id 和被其关注人的id 插入表中 relation
	res := db.Create(&Relation{
		AuthorId:   AuthorId,
		FollowerId: FollowerId,
	})
	//resAuthor := db.Model(&User{}).Updates(map[string]interface{}{"follower_count": gorm.Expr("follower_count+1", AuthorId)})
	//resFollower := db.Model(&User{ID: FollowerId}).Updates(map[string]interface{}{"follow_count": gorm.Expr("follow_count+1", FollowerId)})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Add Relation error")
	}
	//if errors.Is(resAuthor.Error, gorm.ErrRecordNotFound) {
	//	log.Print("Update FollowerCount error")
	//}
	//if errors.Is(resFollower.Error, gorm.ErrRecordNotFound) {
	//	log.Print("Update FollowCount error")
	//}
	log.Print("insert relation success")
	return nil
}

func (r *RelationDao) DeleteRelation(AuthorId uint, FollowerId uint) error {
	//根据 userid followerId 删除对应记录
	res := db.Where(&Relation{
		AuthorId:   AuthorId,
		FollowerId: FollowerId,
	}).Delete(&Relation{})

	//resAuthor := db.Model(&User{ID: AuthorId}).Updates(map[string]interface{}{"follower_count": gorm.Expr("follower_count-", 1)})
	//resFollower := db.Model(&User{ID: FollowerId}).Updates(map[string]interface{}{"follow_count": gorm.Expr("follow_count-", 1)})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Delete relation error")
	}
	//if errors.Is(resAuthor.Error, gorm.ErrRecordNotFound) {
	//	log.Print("Update FollowerCount error")
	//}
	//if errors.Is(resFollower.Error, gorm.ErrRecordNotFound) {
	//	log.Print("Update FollowCount error")
	//}
	log.Print("Delete relation success")
	return nil
}

/* 查询当前用户的关注列表(id)
 */
func (r *RelationDao) QueryFollowIdByAuthorId(AuthorId uint, RelationList *[]Relation) error {
	res := db.Model(&Relation{}).Where("author_id = ?", AuthorId).Find(RelationList)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

/*查询当前用户的粉丝(id)
 */
func (r *RelationDao) QueryAuthorIdByFollowId(FollowerId uint, relationList *[]Relation) error {

	res := db.Model(&Relation{}).Where("follower_id = ?", FollowerId).Find(relationList)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
