package addresscn

import (
	"github.com/hyacinthus/x/xobj"
)

// Service 地址解析服务
type Service struct {
	provider string // 数据来源 github(default)/http/cos
	url      string // http 模式时 文件的地址前缀 包含最后的斜线
	cos      xobj.Client
}
