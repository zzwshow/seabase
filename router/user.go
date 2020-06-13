package router

import (
	"github.com/gin-gonic/gin"
	"seabase/api/v1/base"
)

type UserRouter struct {
}

func (*UserRouter) addGroupApi(rg *gin.RouterGroup) {
	userApi := new(base.UserManagerApi)
	rg.GET("/users", userApi.GetUsers)
}
