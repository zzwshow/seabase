package base

import (
	"github.com/jinzhu/gorm"
	"seabase/model"
)

type UserModel struct {
	UserID   uint64 `gorm:"column:user_id;AUTO_INCREMENT;PRIMARY_KEY" json:"user_id"`
	Username string `gorm:"column:username;size:50;unique_index" json:"username"`
	Password string `gorm:"column:password;size:100" json:"password"`
	Name     string `gorm:"column:name;size:50;default:''" json:"name"`
	Number   string `gorm:"column:number;size:20;default:''" json:"number"`
	Email    string `gorm:"column:email;size:100;default:''" json:"email"`
	Mobile   string `gorm:"column:mobile;size:20;default:''" json:"mobile"`
	Avatar   string `gorm:"column:avatar;size:200;default:''" json:"avatar"`
	Status   uint   `gorm:"column:status;size:2;default:1" json:"status"`
	model.BaseModel
}

func (*UserModel) TableName() string {
	return "sea_users"
}

func (um *UserModel) SelectAll(pageNum int, pageSize int, condition map[string]interface{}) (umList []UserModel, err error) {
	var result *gorm.DB
	result = model.DB.Table(um.TableName()).Where(condition).Offset(pageNum).Limit(pageSize).Find(&umList)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return
}

func (um *UserModel) SelectOneByCondition(condition map[string]interface{})(userInfo *UserModel,err error){
	user := new(UserModel)
	err = model.DB.Table(um.TableName()).Where(condition).First(&user).Error
	return user,err
}


func (um *UserModel) Insert() (err error) {
	err = model.DB.Create(um).Error
	return
}
