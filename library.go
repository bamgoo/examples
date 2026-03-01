package main

import (
	. "github.com/infrago/base"
	_ "github.com/infrago/builtin"

	"github.com/infrago/infra"
)

func init() {

	infra.Register("sendcloud", infra.Library{
		Name: "sendcloud", Desc: "SendCloud邮件",
		Methods: infra.Methods{
			"send": infra.Method{
				Name: "发送邮件", Desc: "发送邮件",
				Action: func(ctx *infra.Context) (Map, Res) {
					return nil, infra.OK
				},
			},
		},
	})
}
