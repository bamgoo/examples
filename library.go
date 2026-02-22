package main

import (
	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"

	"github.com/bamgoo/bamgoo"
)

func init() {

	bamgoo.Register("sendcloud", bamgoo.Library{
		Name: "sendcloud", Desc: "SendCloud邮件",
		Methods: bamgoo.Methods{
			"send": bamgoo.Method{
				Name: "发送邮件", Desc: "发送邮件",
				Action: func(ctx *bamgoo.Context) (Map, Res) {
					return nil, bamgoo.OK
				},
			},
		},
	})
}
