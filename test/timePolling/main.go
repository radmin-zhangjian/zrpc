package main

import (
	"fmt"
	"time"
	"zrpc/utils"
)

func t2(args ...any) {
	fmt.Println(args...)
}

func main() {
	tp := utils.NewPolling()

	tp.Register(10, "t1", func(args ...any) {
		fmt.Println(args...)
	}, []any{1, 2, 3})
	tp.Register(10, "t2", t2, []any{5, 6, 7})
	tp.Register(12, "t2", t2, []any{8, 9, 0})

	//40秒后关闭
	time.AfterFunc(time.Second*20, func() {
		tp.Close()
	})

	tp.Run()

}
