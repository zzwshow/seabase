package e

import "errors"

var (
	CountNonZeroError       = errors.New("数据已被使用")
	OnlyOneDefaultRoleError = errors.New("只能有一个默认角色")
	OnlyOneAdminRoleError   = errors.New("只能有一个管理角色")
	SameNameError           = errors.New("同名")
	NoQueryCondition        = errors.New("没有查询条件")
	InvalidPassword         = errors.New("密码错误")
	UserDoesNotExist        = errors.New("User does not exist or to many entries returned.")
	SameRecordExsit         = errors.New("存在相同记录")
	TaskStructToHandler     = errors.New("没有找到实现任务接口的结构体")
	NoFound                 = errors.New("查询数据没找到")
	CreateTaskFaild         = errors.New("创建任务执行函数失败")
	InvalidDeployId         = errors.New("Invalid deploy id")
	AssertionFailed         = errors.New("类型断言失败")
	GetMapKeyFailed         = errors.New("Failed to get map key")
	ModuleNoFound           = errors.New("模块不存在")
	NotEnoughData           = errors.New("没有足够的数据")
	UploadFileError         = errors.New("接收上传文件失败")
	RecordDoesNotExist      = errors.New("没有相关记录")
	HTTPResponseStausNot200 = errors.New("http响应状态非200")
	IpIsInvalid             = errors.New("无效的ip地址")
)