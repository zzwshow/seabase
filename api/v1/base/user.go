package base

import (
	"github.com/gin-gonic/gin"
	"seabase/extend/e"
	"seabase/extend/util"
	"seabase/extend/vars"
	"seabase/service/base"
)

type UserManagerApi struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Mobile   string `json:"mobile" validate:"min=6,max=11"`
	Name     string `json:"name" validate:"required"`
	Number   string `json:"number"`
	Email    string `json:"email"`
}

func (uma *UserManagerApi) GetUsers(c *gin.Context) {
	uS := new(base.UserService)
	condition := make(map[string]interface{})
	if param := c.Query("user_id"); param != "" {
		condition["user_id"] = param
	}
	data, err := uS.GetUsersAll(util.GetPage(c), util.GetPageSize(c), condition)
	if err != nil {
		e.Resp(c, e.RequestParamError, nil)
		return
	}
	e.Resp(c, e.Success, vars.CommonMap{"data": data})
	return
}

func (uma *UserManagerApi) CreateUser(c *gin.Context) {
	if err := c.ShouldBindJSON(uma); err != nil {
		e.Resp(c, err, nil)
		return
	}
	if err := util.AnalysisErrorByValidate(util.Vd.Struct(uma)); err != nil {
		e.Resp(c, err, nil)
		return
	}
	uS := new(base.UserService)
	uS.Username = uma.UserName
	uS.Password = uma.Password
	uS.Name = uma.Name
	uS.Mobile = uma.Mobile
	uS.Email = uma.Email
	uS.Number = uma.Number
	err := uS.CreateUser()
	if err != nil {
		e.Resp(c, err, nil)
		return
	}
	e.Resp(c, e.Success, nil)
	return
}
