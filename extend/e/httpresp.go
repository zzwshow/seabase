package e

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpResponse struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	Success           = &HttpResponse{http.StatusOK, 20000, "请求成功"}
	RequestParamError = &HttpResponse{http.StatusBadRequest, 40001, "传入参数有误!!!"}
	TokenNotFound           = &HttpResponse{http.StatusUnauthorized, 4010002, "Token required."}
	TokenInvalid            = &HttpResponse{http.StatusUnauthorized, 4010003, "Invalid Token."}
	TokenExpired            = &HttpResponse{http.StatusUnauthorized, 4010004, "Token expired."}
	AlarmKeyInvalid         = &HttpResponse{http.StatusUnauthorized, 4010005, "accessKey或secretKey Invalid."}
	AccountPassUnmatch = &HttpResponse{http.StatusBadRequest, 4000002, "该账号原密码不匹配"}
	//SignupPassUnmatch = &Code{http.StatusBadRequest, 4000003, "注册两次输入密码不匹配"}
	//AccountNameExist = &Code{http.StatusBadRequest, 4000004, "账号昵称已被使用"}
	//UploadSuffixError = &Code{http.StatusBadRequest, 4000005, "该上传文件格式目前暂不支持"}
	//UploadSizeLimit = &Code{http.StatusBadRequest, 4000006, "目前上传仅支持小于5M的文件内容"}
	NotMySelf               = &HttpResponse{http.StatusUnauthorized, 4010001, "非本人"}

	ServiceInsideError      = &HttpResponse{http.StatusInternalServerError, 50000, "服务器内部错误"}
	UsernamePasswordError   = &HttpResponse{http.StatusInternalServerError, 50001, "您的用户名或密码输入错误"}
	Uniqueness              = &HttpResponse{http.StatusInternalServerError, 60000, "数据已经存在了，请核对一下哦"}
	RecordNotFound          = &HttpResponse{http.StatusInternalServerError, 60001, "没有查询到相关的记录"}
	ExistBindRelation       = &HttpResponse{http.StatusInternalServerError, 60002, "存在绑定关系"}
	AppIpInfo               = &HttpResponse{http.StatusInternalServerError, 60003, "请先解除应用主机绑定关系,在删除应用模块!"}
	BindHostRelation        = &HttpResponse{http.StatusInternalServerError, 60003, "应用未绑定主机,请联系运维人员"}
	BindIDCRelation         = &HttpResponse{http.StatusInternalServerError, 60004, "应用主机未绑定IDC,请联系运维人员"}
	BindAgentRelation       = &HttpResponse{http.StatusInternalServerError, 60005, "未绑定agent,请联系运维人员"}
	RequestNotAllow         = &HttpResponse{http.StatusForbidden, 403001, "拒绝访问."}
	DataHasBeenUsed         = &HttpResponse{http.StatusInternalServerError, 500011, "数据已被使用."}
	OnlyOneAdminRole        = &HttpResponse{http.StatusInternalServerError, 500011, "只能有一个管理角色"}
	OnlyOneDefaultRole      = &HttpResponse{http.StatusInternalServerError, 500011, "只能有一个默认角色"}
	OtherError              = &HttpResponse{http.StatusBadRequest, 500011, "内部错误!!!"}
	ArthasAgentStartError   = &HttpResponse{http.StatusInternalServerError, 40011, "Arthas Agent操作异常."}
	IpUsedAlready           = &HttpResponse{http.StatusInternalServerError, 40012, "ip已使用，网段不能删除"}
	ArthasAgentDoesNotExist = &HttpResponse{http.StatusInternalServerError, 40013, "Arthas Agent ID 不存在"}
)

func Resp(c *gin.Context, resp interface{}, data interface{}) {
	if resp == nil {
		resp = ServiceInsideError
	}
	switch resp := resp.(type) {
	case *HttpResponse:
		c.JSON(resp.Status, gin.H{
			"code": resp.Code,
			"msg":  resp.Message,
			"data": data,
		})
	case error:
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  resp.Error(),
			"data": data,
		})
	}

	return
}
