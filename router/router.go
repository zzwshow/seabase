package router

import (
	"github.com/gin-gonic/gin"
	"seabase/api/v1/base"
	"seabase/middleware/jwt"
)

type RouterInterface interface {
	addGroupApi(rg *gin.RouterGroup)
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	userLoginHandler := new(base.UserLoginApi)
	r.POST("/api/v1/login",userLoginHandler.UserLogin)
	userManagerHandler := new(base.UserManagerApi)
	r.POST("/api/v1/users",userManagerHandler.CreateUser)
	
	apiV1 := r.Group("api/v1")
	apiV1.Use(jwt.JWT())
	new(UserRouter).addGroupApi(apiV1)

	return r
}
