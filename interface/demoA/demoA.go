package demoA

import "context"

// 入参数
type ArgsAdd struct {
	a int
	b int
}

//出参数
type ReplyAdd struct {
	sum int
}
type DemoAdd interface {
	Add(ctx context.Context, add ArgsAdd, replyAdd ReplyAdd) error
}