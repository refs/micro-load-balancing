package main

import (
	"context"
	"flag"
	"fmt"
	"micronaming/pkg/proto"
	"micronaming/pkg/service"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/mdns"
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
	helloMsg   string
	instanceMD string
	qMode      bool
)

func allNodes(r registry.Registry) {
	registryServices, err := r.GetService("go.micro.api.hello")
	if err != nil {
		panic(err)
	}

	for _, serv := range registryServices {
		// list available nodes.
		if len(serv.Nodes) > 1 {
			for _, node := range serv.Nodes {
				// select nodes with a present metadata value.
				if _, ok := node.Metadata["ocis_instance"]; !ok {
					continue
				}
				fmt.Printf("addr: %v\nmetadata: %+v\n", node.Address, node.Metadata)
			}
		}
	}
}

func queriedNode(r registry.Registry) {
	registryServices, err := r.GetService("go.micro.api.hello")
	if err != nil {
		panic(err)
	}

	for _, serv := range registryServices {
		// list available nodes.
		if len(serv.Nodes) > 1 {
			for _, node := range serv.Nodes {
				// select nodes with a present metadata value.
				if _, ok := node.Metadata["ocis_instance"]; !ok {
					continue
				}

				// select the ones which label matches the queried
				if instanceMD != "" && node.Metadata["ocis_instance"] == instanceMD {
					fmt.Printf("addr: %v\n", node.Address)
					break
				}

			}
		}
	}
	return
}

func main() {
	flag.StringVar(&helloMsg, "msg", "", "message to be printed")
	flag.StringVar(&instanceMD, "ins", "", "set instance metadata")
	flag.BoolVar(&qMode, "q", false, "query micro registry for a list of nodes")

	flag.Parse()

	if qMode {
		r := mdns.NewRegistry()
		if instanceMD != "" {
			queriedNode(r)
			return
		}
		allNodes(r)
		return
	}

	var md map[string]string
	if instanceMD != "" {
		md = map[string]string{
			"ocis_instance": instanceMD,
		}
	}

	// create both services and run them on goroutines.
	hello1 := service.NewHello(
		"go.micro.api.hello",
		&Hello{msg: helloMsg},
		md,
	)
	if err := hello1.Run(); err != nil {
		panic(err)
	}
}
