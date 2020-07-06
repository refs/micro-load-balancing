package service

import (
	"micronaming/pkg/proto"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
)

// NewHello creates a new micro service.
func NewHello(name string, handler proto.SayHandler) micro.Service {
	s := micro.NewService(
		micro.Name(name),
		micro.Flags(&cli.StringFlag{
			Name: "msg",
		}),
	)

	s.Init()
	proto.RegisterSayHandler(s.Server(), handler)

	return s
}
