package repository

import (
	"demo1/model/entity"
	"sync"
	"time"
)

type CommentDAO struct {
}

var (
	commentDAO  *CommentDAO
	commentOnce sync.Once
)

func NewCommentDAO() *CommentDAO {
	commentOnce.Do(func() {
		commentDAO = &CommentDAO{}
	})
	return commentDAO
}

// CreateComment 增加评论，返回评论的id和错误类型
func (c *CommentDAO) CreateComment(uid uint, videoID uint, content *string) (uint, error) {
	comment := &entity.Comment{
		CommentPublishTime: time.Now().Unix(),
		Content:            *content,
		AuthorID:           uid,
		VideoID:            videoID,
	}
	res := db.Create(&comment)
	if res.Error != nil {
		return 0, res.Error
	}
	return comment.ID, nil // 返回插入数据的主键
}

// DeleteCommentById 通过评论的id删除评论，成功删除返回true
func (c *CommentDAO) DeleteCommentById(cid uint) error {
	res := db.Delete(&entity.Comment{}, cid)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (c *CommentDAO) VideoCommentCount(vid uint) int64 {
	var count int64
	db.Model(&entity.Comment{}).Where("video_id = ?", vid).Count(&count)
	return count
}

func (c *CommentDAO) GetAllComment(list *[]entity.Comment, vid uint) error {
	res := db.Model(&entity.Comment{}).Order("comment_publish_time desc").Where("video_id = ?", vid).Find(list)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
