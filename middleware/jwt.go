package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

var mySigningKey = []byte("dousheng")

const TokenExpireDuration = time.Hour * 24

// DouShengClaims 自定义Claims
type DouShengClaims struct {
	// 自行添加的字段
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

// GenToken 生成token
func GenToken(username string, uid uint) (string, error) {
	// 创建一个自己的声明
	claims := DouShengClaims{
		Username: username,
		UserID:   uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "lzd",
		},
	}

	// 使用指定签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySigningKey)
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*DouShengClaims, error) {
	// 解析token
	// 如果自定义claim结构体需要使用
	token, err := jwt.ParseWithClaims(tokenString, &DouShengClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 对token对象中的claim进行类型断言
	if claims, ok := token.Claims.(*DouShengClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//fmt.Printf("%+v\n", c.Request)
		token := c.DefaultQuery("token", "")
		fmt.Println("in jwt middle, token : ", token)
		if token == "" {
			token = c.PostForm("token")
			if token == "" {
				fmt.Println("token is empty")
				c.JSON(http.StatusUnauthorized, gin.H{
					"status_code": 1,
					"status_msg":  "token is empty",
				})
				c.Abort()
				return
			}
		}

		// dc:DouShengClaims
		dc, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("username", dc.Username)
		c.Set("user_id", dc.UserID)
		fmt.Println("jwt username : ", dc.Username)
		//c.Next()
	}
}
