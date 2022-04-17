package demoA

import "context"

// 入参数
type ArgsAdd struct {
	Num_1 int
	Num_2 int
}

//出参数
type ReplyAdd struct {
	Sum int
}

// 定义接口
type DemoAdd interface {
	Add(ctx context.Context, add ArgsAdd, replyAdd ReplyAdd) error
}
