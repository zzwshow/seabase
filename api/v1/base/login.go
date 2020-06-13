package base

import (
	"github.com/gin-gonic/gin"
	"seabase/extend/e"
	"seabase/extend/util"
	"seabase/extend/vars"
	"seabase/service/base"
)

type UserLoginApi struct {
	UserName string `json:"username" validate:"min=4,max=11"`
	Password string `json:"password" validate:"min=6,max=30"`
}


func (uma *UserLoginApi) UserLogin(c *gin.Context) {
	userLoginInfo := new(UserLoginApi)
	if err := c.ShouldBindJSON(userLoginInfo);err != nil{
		e.Resp(c,err,nil)
		return
	}
	if err := util.AnalysisErrorByValidate(util.Vd.Struct(userLoginInfo));err != nil{
		e.Resp(c,err,nil)
		return
	}
	uS := new(base.UserLoginService)
	uS.UserName = userLoginInfo.UserName
	uS.Password = userLoginInfo.Password
	data,err := uS.UserLogin()
	if err != nil{
		e.Resp(c,err,nil)
		return
	}
	e.Resp(c,e.Success,vars.CommonMap{"data":data})
	return
}