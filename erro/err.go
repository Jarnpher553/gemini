package erro

type Err struct {
	Code int
	Msg  string
}

func (e *Err) Error() string {
	return e.Msg
}

// 错误码
const (
	ErrSuccess        = 200
	ErrDefault        = 500
	ErrBreaker        = 3000
	ErrMaxRequest     = 3005
	ErrRateLimiter    = 3001
	ErrDelayLimiter   = 3002
	ErrReserveLimiter = 3004
	ErrAuthor         = 4003

	//add custom code below
	ErrFileMime   = 501
	ErrNoFile     = 502
	ErrReqContent = 503
	ErrDbRead     = 504
	ErrDbModify   = 505
	ErrDbRemove   = 506
	ErrUserName   = 507
	ErrPassword   = 508
	ErrPermission = 509
	ErrToken      = 510
	ErrDbInsert   = 511
	ErrTemplate   = 512
	ErrImport     = 513
	ErrExport     = 514
	ErrNotExist   = 517
	ErrClosed     = 519
	ErrDb         = 520
	ErrFileOpen   = 521
	ErrFileParse  = 522
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
	ErrDbModify:       "更新失败",
	ErrDbRemove:       "删除失败",
	ErrDbInsert:       "新增失败",
	ErrUserName:       "用户名错误",
	ErrPassword:       "密码错误",
	ErrPermission:     "无操作权限",
	ErrToken:          "令牌生成失败",
	ErrTemplate:       "下载模板失败",
	ErrImport:         "上传文件失败",
	ErrExport:         "导出数据失败",
	ErrNotExist:       "不存在记录",
	ErrDb:             "数据库错误",
	ErrFileOpen:       "文件读取失败",
	ErrFileParse:      "文件解析失败",
}

func Register(code int, msg string) {
	ErrMsg[code] = msg
}
