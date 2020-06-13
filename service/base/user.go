package base

import (
	"seabase/extend/conf"
	"seabase/extend/util"
	"seabase/model/base"
)

type UserService struct {
	Username string
	Password string
	Name     string
	Number   string
	Email    string
	Mobile   string
}

func (us *UserService) GetUsersAll(pageNum int, pageSize int, condition map[string]interface{}) (umList []base.UserModel, err error) {
	um := new(base.UserModel)
	umList, err = um.SelectAll(pageNum, pageSize, condition)
	return
}

func (us *UserService) CreateUser() (err error) {
	um := new(base.UserModel)
	um.Username = us.Username
	newP := util.AesEncrypt(us.Password,conf.ServerConf.DBSecret)
	um.Password = newP
	um.Name = us.Name
	um.Email = us.Email
	um.Mobile = us.Mobile
	err = um.Insert()
	return
}
