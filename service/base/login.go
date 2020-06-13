package base

import (
	"encoding/json"
	"fmt"
	"seabase/extend/conf"
	"seabase/extend/e"
	"seabase/extend/redis"
	"seabase/extend/util"
	"seabase/model/base"
)

type UserLoginService struct {
	UserName string
	Password string
}

type CurrentUser struct {
	UserId          uint64              `json:"UserId"`
	UserName        string              `json:"UserName"`
	Number 			string				`json:"Number"`
	Name            string              `json:"Name"`
	Email           string              `json:"Email"`
	Mobile          string              `json:"Mobile"`
	Avatar          string              `json:"Avatar"`
	IsAdmin         uint                `json:"IsAdmin"`
}

func (uls *UserLoginService) UserLogin()(reqData map[string]interface{}, err error ){
	req := make(map[string]interface{})
	uM := new(base.UserModel)
	condition := map[string]interface{}{"username": uls.UserName}
	userInfo,err := uM.SelectOneByCondition(condition)
	if err != nil {
		return
	}
	if userInfo == nil{
		return nil,e.NoFound
	}
	passwod := util.AesDecrypt(userInfo.Password, conf.ServerConf.DBSecret)
	if uls.Password == passwod{
		redisKey := fmt.Sprintf("token_%s",userInfo.Username)
		token,err :=  util.GenerateToken(userInfo.Username,redisKey)
		if err != nil{
			return nil,err
		}
		// 保存到redis
		currentUser := CurrentUser{
			UserId: userInfo.UserID,
			UserName: userInfo.Username,
			Number: userInfo.Number,
			Name: userInfo.Name,
			Email: userInfo.Email,
			Mobile: userInfo.Mobile,
			IsAdmin: 0,
		}
		redisVal,err := json.Marshal(currentUser)
		err = redis.Set(redisKey,string(redisVal),conf.ServerConf.JWTExpire *60*60)
		if err != nil {
			return nil,err
		}
		// 返回客户端
		req["token"] = token
		req["username"] = userInfo.Username
		req["user_id"] = userInfo.UserID
		req["name"] = userInfo.Name
		return req,nil
	}
	return nil, e.InvalidPassword
}
