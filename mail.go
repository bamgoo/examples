package main

import (
	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"

	"github.com/bamgoo/bamgoo"
)

type (
	MailProvider = interface {
		Send(Map) (Map, Res)
	}

	SendCloudBuilder struct{}
	SendCloudMailer  struct {
		Setting Map
	}
)

func init() {
	bamgoo.Register("sendcloud", &SendCloudBuilder{})
}

func (sc *SendCloudBuilder) UseProvider(setting Map) (Any, error) {
	return &SendCloudMailer{Setting: setting}, nil
}

func (sc *SendCloudMailer) Send(args Map) (Map, Res) {
	return Map{"msg": "发送成功！"}, bamgoo.OK
}
