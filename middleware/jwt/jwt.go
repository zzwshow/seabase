package jwt

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"seabase/extend/e"
	"seabase/extend/redis"
	"seabase/extend/util"
	"seabase/service/base"
)

// JWT middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-Token")
		if token == "" {
			e.Resp(c,e.TokenNotFound,nil)
			c.Abort()
			return
		} else {
			claim,err := util.ParseToken(token)
			if err != nil{
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					e.Resp(c,e.TokenExpired,nil)
					c.Abort()
					return
				default:
					e.Resp(c,e.TokenInvalid,nil)
					c.Abort()
					return
				}
			}
			rdsVal,err := redis.Get(claim.RedisKey)
			if err != nil{
				e.Resp(c,err,nil)
				c.Abort()
				return
			}
			currentUser := &base.CurrentUser{}
			err = json.Unmarshal([]byte(rdsVal),currentUser)
			if err != nil {
				e.Resp(c,err,nil)
				c.Abort()
				return
			}
			if currentUser.UserName == "admin" {
				currentUser.IsAdmin = 1
			}
			c.Set("currentUser", currentUser)
			c.Set("X-Token", token)
		}
		c.Next()
	}
}