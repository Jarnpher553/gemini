package erro

type Err struct {
	Code int
	Msg  string
}

func (e *Err) Error() string {
	return e.Msg
}

var (
	successCode   int
	isSuccessCode bool
)

func SetSuccess(code int) {
	successCode = code
	isSuccessCode = true
}

func Success() int {
	if isSuccessCode {
		return successCode
	} else {
		return ErrSuccess
	}
}

// 错误码
const (
	ErrSuccess        = 200
	ErrDefault        = 500
	ErrReqContent     = 503
	ErrBreaker        = 3000
	ErrMaxRequest     = 3005
	ErrRateLimiter    = 3001
	ErrDelayLimiter   = 3002
	ErrReserveLimiter = 3004
	ErrAuthor         = 403
	ErrFileMime       = 501
	ErrNoFile         = 502
	ErrDbRead         = 504
	ErrDbModify       = 505
	ErrDbRemove       = 506
	ErrUserName       = 507
	ErrPassword       = 508
	ErrPermission     = 401
	ErrToken          = 510
	ErrDbInsert       = 511
	ErrTemplate       = 512
	ErrImport         = 513
	ErrExport         = 514
	ErrNotExist       = 517
	ErrDb             = 520
)

// 错误码对应错误信息
var ErrMsg = map[int]string{
	ErrSuccess:        "请求成功",
	ErrDefault:        "内部错误",
	ErrBreaker:        "服务熔断",
	ErrMaxRequest:     "熔断限流",
	ErrRateLimiter:    "服务繁忙",
	ErrDelayLimiter:   "服务等待",
	ErrReserveLimiter: "服务保持",
	ErrAuthor:         "未获取授权",
	ErrFileMime:       "文件类型错误",
	ErrNoFile:         "未上传文件",
	ErrReqContent:     "请求参数有误",
	ErrDbRead:         "查询失败",
	ErrDbModify:       "修改失败",
	ErrDbRemove:       "删除失败",
	ErrDbInsert:       "添加失败",
	ErrUserName:       "用户名错误",
	ErrPassword:       "密码错误",
	ErrPermission:     "无操作权限",
	ErrToken:          "令牌生成失败",
	ErrTemplate:       "下载模板失败",
	ErrImport:         "上传失败",
	ErrExport:         "导出失败",
	ErrNotExist:       "不存在记录",
	ErrDb:             "数据库异常",
}

func Register(code int, msg string) {
	ErrMsg[code] = msg
}
