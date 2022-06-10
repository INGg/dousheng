package service

import (
	"demo1/model"
	"demo1/model/entity"
	"demo1/repository"
	"errors"
)

//单例模式
var relationDao = repository.NewRelationDAO()
var userDao = repository.NewUserDAO()
var isSFollow bool

// RelationAction 关注操作
//func RelationAction(c *gin.Context) {
//	var req model.FollowActionRequest
//	if err := c.ShouldBind(&req); err != nil {
//		c.JSON(200, model.Response{
//			StatusCode: 1,
//			StatusMsg:  "参数解析失败",
//		})
//		// 成功获取参数
//	} else {
//		// 判断是否已经关注`
//		if isFollow, err := relationDao.IsFollow(req.UserID, req.ToUserID); err != nil {
//			c.JSON(200, model.Response{
//				StatusCode: 1,
//				StatusMsg:  "查询数据失败",
//			})
//			isSFollow = isFollow
//			return
//		}
//		// 关注操作
//		req.UserID = c.GetUint("user_id")
//		// 判断 关注还是取消关注
//		if req.ActionType == 1 {
//			if isSFollow {
//
//				c.JSON(200, model.Response{
//					StatusCode: 1,
//					StatusMsg:  "你已经关注他了",
//				})
//				return
//			}
//			if err := relationDao.AddRelation(req.UserID, req.ToUserID); err != nil {
//				c.JSON(200, model.Response{
//					StatusCode: 1,
//					StatusMsg:  "关注失败",
//				})
//				return
//			}
//			c.JSON(200, model.Response{
//				StatusCode: 0,
//				StatusMsg:  "关注成功",
//			})
//			//	取消关注操作
//		} else if req.ActionType == 2 {
//			fmt.Println("UserId: ", req.UserID)
//			fmt.Println("ToUSerId: ", req.ToUserID)
//			if !isSFollow {
//				c.JSON(200, model.Response{
//					StatusCode: 1,
//					StatusMsg:  "你还没关注他呢",
//				})
//				return
//			}
//			if err := relationDao.DeleteRelation(req.UserID, req.ToUserID); err != nil {
//
//				c.JSON(200, model.Response{
//					StatusCode: 1,
//					StatusMsg:  "取消关注失败",
//				})
//				return
//			}
//			c.JSON(200, model.Response{
//				StatusCode: 0,
//				StatusMsg:  "取消关注成功",
//			})
//		} else {
//			c.JSON(200, model.Response{
//				StatusCode: 1,
//				StatusMsg: "ActionType 解析失败",
//			}
//		}
//	}
//}
func AddRelation(req *model.FollowActionRequest) (*model.FollowActionResponse, error) {
	if req.UserID == req.ToUserID {
		return &model.FollowActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "you can't follow yourself",
			},
		}, errors.New("can't follow yourself")
	}

	//单例模式
	relationDAO := repository.NewRelationDAO()
	userDAO := repository.NewUserDAO()

	// 向relation表中写入关注
	if err := relationDAO.AddRelation(req.UserID, req.ToUserID); err != nil {
		return &model.FollowActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "AddRelation error",
			},
		}, err
	}

	// 用户发起者的关注人数加一
	if err := userDAO.UpdateUserFollowCount(req.UserID); err != nil {
		return &model.FollowActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "AddRelation error",
			},
		}, err
	}

	// 被关注人的粉丝数加一
	if err := userDAO.UpdateUserFollowerCount(req.ToUserID); err != nil {
		return &model.FollowActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "AddRelation error",
			},
		}, err
	}

	// 结果成功返回
	return &model.FollowActionResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
	}, nil
}

func CancelRelation(req *model.FollowActionRequest) (*model.FollowActionResponse, error) {
	//单例模式
	relationDAO := repository.NewRelationDAO()
	userDAO := repository.NewUserDAO()

	// 将关注关系删除
	if err := relationDAO.DeleteRelation(req.UserID, req.ToUserID); err != nil {
		return &model.FollowActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "cancel relation error",
			},
		}, err
	}

	// 用户发起者的关注人数减一
	if err := userDAO.ReduceFollowCount(req.UserID); err != nil {
		return &model.FollowActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "cancel relation error",
			},
		}, err
	}

	// 取关的人的粉丝数减一
	if err := userDAO.ReduceFollowerCount(req.ToUserID); err != nil {
		return &model.FollowActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "cancel relation error",
			},
		}, err
	}

	return &model.FollowActionResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
	}, nil
}

// FollowList 获取关注列表
func FollowList(req *model.UserFollowListRequest) (*model.UserFollowListResponse, error) {

	// 创建单例
	UserDAO := repository.NewUserDAO()
	relationDAO := repository.NewRelationDAO()

	//准备参数
	var userList []entity.User
	var userListR []entity.Relation

	// 找到他关注的用户的所有id
	if err := relationDAO.QueryUsersIDByFollowId(req.UserID, &userListR); err != nil { // 说明没有关注
		return &model.UserFollowListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "find user follow error",
			},
			UserList: nil,
		}, nil
	}

	if len(userListR) == 0 { // 他没关注任何人
		return &model.UserFollowListResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "ok, but list is nil",
			},
			UserList: nil,
		}, nil
	}

	// 获取author中的 authorId
	var userIdList = make([]uint, len(userListR))
	for i, author := range userListR {
		userIdList[i] = author.UserID
	}

	// 查找他关注的用户的信息
	if err := UserDAO.FindUsersByIdList(userIdList, &userList); err != nil {
		return &model.UserFollowListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get follower error",
			},
			UserList: nil,
		}, err
	}

	if req.FromUserID != 0 {
		for i := 0; i < len(userList); i++ {
			userList[i].IsFollow = relationDAO.QueryAFollowB(req.FromUserID, req.UserID)
		}
	}

	return &model.UserFollowListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		UserList: &userList,
	}, nil
}

// FollowerList 获取粉丝列表
func FollowerList(req *model.UserFollowerListRequest) (*model.UserFollowerListResponse, error) {
	// 创建单例
	relationDAO := repository.NewRelationDAO()
	userDAO := repository.NewUserDAO()

	var followerList []entity.User
	var relationList []entity.Relation

	if err := relationDAO.QueryFollowIdByUserID(req.UserID, &relationList); err != nil {
		return &model.UserFollowerListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get follower list error",
			},
			UserList: nil,
		}, err
	}

	if len(relationList) == 0 { // 说明该用户没有粉丝
		return &model.UserFollowerListResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "ok, but list is nil",
			},
			UserList: nil,
		}, nil

	}

	// 构造答案
	followerID := make([]uint, len(relationList))
	for i, relation := range relationList {
		followerID[i] = relation.FollowID
	}

	if err := userDAO.FindUsersByIdList(followerID, &followerList); err != nil {
		return &model.UserFollowerListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get follower list error",
			},
			UserList: nil,
		}, err
	}

	if req.FromUserID != 0 {
		for i := 0; i < len(followerList); i++ {
			followerList[i].IsFollow = relationDAO.QueryAFollowB(req.FromUserID, req.UserID)
		}
	}

	return &model.UserFollowerListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		UserList: &followerList,
	}, nil
}
