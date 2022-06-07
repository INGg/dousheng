package service

import (
	"demo1/repository"
	"fmt"
	"github.com/gin-gonic/gin"
)

// RelationAction 关注操作
func RelationAction(c *gin.Context) {
	var req FollowActionRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, Response{
			StatusCode: 1,
			StatusMsg:  "参数解析失败",
		})
		// 成功获取参数
	}

	//单例模式
	relationDAO := repository.NewRelationDAO()
	userDAO := repository.NewUserDAO()

	// 关注操作
	req.UserId = c.GetUint("user_id")
	// 判断 关注还是取消关注
	if req.ActionType == 1 {
		if err := relationDAO.AddRelation(req.UserId, req.ToUserId); err != nil {
			c.JSON(200, Response{
				StatusCode: 1,
				StatusMsg:  "关注失败",
			})
			return
		}
		if err := userDAO.UpdateUserFollowerCount(req.UserId); err != nil {
			c.JSON(200, Response{
				StatusCode: 1,
				StatusMsg:  "关注失败",
			})
			return
		}

		if err := userDAO.UpdateUserFollowerCount(req.ToUserId); err != nil {
			c.JSON(200, Response{
				StatusCode: 1,
				StatusMsg:  "关注失败",
			})
			return
		}

		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "关注成功",
		})
		//	取消关注操作
	} else if req.ActionType == 2 {
		fmt.Println("UserId: ", req.UserId)
		fmt.Println("ToUSerId: ", req.ToUserId)
		if err := relationDAO.DeleteRelation(req.UserId, req.ToUserId); err != nil {
			c.JSON(200, Response{
				StatusCode: 1,
				StatusMsg:  "取消关注失败",
			})
			return
		}
		if err := userDAO.ReduceFollowerCount(req.UserId); err != nil {
			c.JSON(200, Response{
				StatusCode: 1,
				StatusMsg:  "取消关注失败",
			})
			return
		}
		if err := userDAO.ReduceFollowerCount(req.ToUserId); err != nil {
			c.JSON(200, Response{
				StatusCode: 1,
				StatusMsg:  "取消关注失败",
			})
			return
		}
		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "取消关注成功",
		})
	} else {
		c.JSON(200, Response{
			StatusCode: 1,
			StatusMsg:  "ActionType解析失败",
		})
	}

}

// FollowList 获取关注列表
func FollowList(c *gin.Context) {
	var req UserFollowerListRequest
	var AuthorList []repository.UserRes
	UserD := repository.NewUserDAO()
	// 无法将参数赋值到req中
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, UserFollowListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "无法获取关注列表",
			},
			UserList: nil,
		})
		//	成功赋值给req,开始获取关注列表
	} else {
		// 创建单例
		relationDAO := repository.NewRelationDAO()

		//准备参数
		req.UserId = c.GetUint("user_id")
		var AuthorListR []repository.Relation

		relationDAO.QueryAuthorIdByFollowId(req.UserId, &AuthorListR)

		if AuthorListR == nil {
			c.JSON(200, UserFollowListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "你还没有关注呢",
				},
				UserList: nil,
			})
		} else {
			var AuthorIdList = make([]uint, len(AuthorList))
			// 获取author中的 authorId
			for i, author := range AuthorListR {
				AuthorIdList = append(AuthorIdList, author.AuthorId)
				fmt.Println("AuthorId: ", author.AuthorId)
				fmt.Println("AuThorIDList:", AuthorIdList[i])
			}

			err := UserD.FindMUserByIdList(AuthorIdList, &AuthorList)
			//fmt.Println("AuthorList:", AuthorList[0])
			// 无法从数据库中找到对应id列表

			if err != nil {
				c.JSON(200, UserFollowListResponse{
					Response: Response{
						StatusCode: 1,
						StatusMsg:  "用户查询出错",
					},
					UserList: nil,
				})
			} else {
				c.JSON(200, UserFollowListResponse{
					Response: Response{
						StatusCode: 0,
						StatusMsg:  "Success",
					},
					UserList: &AuthorList,
				})
			}
		}
	}
}

// FollowerList 获取粉丝列表
func FollowerList(c *gin.Context) {
	var req UserFollowerListRequest
	var FollowerList []repository.UserRes
	// 参数错误.无法赋值到req中
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, UserFollowerListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "参数错误",
			},
			UserList: nil,
		})
		//	已经成功将参数赋值到req中
	} else {
		// 创建单例
		relationDAO := repository.NewRelationDAO()
		userDAO := repository.NewUserDAO()

		// jwt 已经将 token中的user_id写入 gin.context中
		//req.UserId = c.GetUint("user_id")
		var RelationList []repository.Relation

		relationDAO.QueryFollowIdByAuthorId(req.UserId, &RelationList)

		if RelationList == nil {
			c.JSON(200, UserFollowerListResponse{
				Response: Response{
					StatusCode: 0,
					StatusMsg:  "你还没有粉丝呢",
				},
				UserList: nil,
			})

		} else {
			FollowerId := make([]uint, len(RelationList))
			for _, relation := range RelationList {
				FollowerId = append(FollowerId, relation.FollowerId)
			}
			if err := userDAO.FindMUserByIdList(FollowerId, &FollowerList); err != nil {
				c.JSON(200, UserFollowListResponse{
					Response: Response{
						StatusCode: 1,
						StatusMsg:  "用户查询出错",
					},
					UserList: nil,
				})
			} else {
				c.JSON(200, UserFollowerListResponse{
					Response: Response{
						StatusCode: 0,
						StatusMsg:  "success",
					},
					UserList: &FollowerList,
				})
			}

		}
	}
}
