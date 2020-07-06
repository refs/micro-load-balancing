package main

import (
	"context"
	"flag"
	"micronaming/pkg/proto"
	"micronaming/pkg/service"

	"github.com/golang/protobuf/ptypes/empty"
)

// Hello implements the SayHandler interface
type Hello struct {
	msg string
}

// Hello implements the SayHandler interface
func (h *Hello) Hello(ctx context.Context, a *empty.Empty, r *proto.Response) error {
	r.Msg = h.msg
	return nil
}

var (
	helloMsg      string
	queryRegistry bool
)

func main() {
	flag.StringVar(&helloMsg, "msg", "", "message to be printed")
	flag.BoolVar(&queryRegistry, "q", false, "query micro registry for a list of nodes")
	flag.Parse()

	// create both services and run them on goroutines.
	hello1 := service.NewHello("go.micro.api.hello", &Hello{msg: helloMsg})
	if err := hello1.Run(); err != nil {
		panic(err)
	}
}
