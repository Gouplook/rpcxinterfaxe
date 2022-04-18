package demoA

import (
	"context"
	"github.com/Gouplook/rpcxinterfaxe/client"
	"github.com/Gouplook/rpcxinterfaxe/interface/demoA"
)

type DemoA struct {
	client.BaseClient
}

func (d *DemoA) Init() *DemoA {
	d.ServiceName = "rpcx_A"
	d.ServicePath = "DemoA"
	return d
}

func (d *DemoA) Add(ctx context.Context, ArgsAdd *demoA.ArgsAdd, replyAdd *demoA.ReplyAdd) error {
	return d.Call(ctx, "Add", ArgsAdd, replyAdd)
}
