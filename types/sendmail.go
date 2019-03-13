package types

import (
	"github.com/google/uuid"
)

// 发送邮件模块

// SrvSendMailUUID 发送邮件服务
var SrvSendMailUUID = uuid.UUID{0xa, 0x6f, 0x7d, 0xe5, 0x1, 0xf, 0x45, 0xf6, 0x8b, 0xa2, 0x40, 0xf2, 0x95, 0x67, 0x55, 0x80}

// MailAddr 邮件地址属性
type MailAddr struct {
	Addr string // 邮件地址
	Name string // 收件人姓名
}

// SendMailService 发送邮件接口
type SendMailService interface {
	SendTo(MailAddr, string, string) error // 发送给某人,主题,内容
}
