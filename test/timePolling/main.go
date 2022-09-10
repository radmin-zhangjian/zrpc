package main

import (
	"fmt"
	"time"
	"zrpc/utils"
)

func t2(args ...any) {
	time.Sleep(8 * time.Second)
	fmt.Println(args...)
}

func main() {
	tp := utils.NewPolling(60)

	tp.Register(5, "t1", func(args ...any) {
		fmt.Println(args...)
	}, []any{1, 2, 3})
	tp.Register(6, "t2", t2, []any{5, 6, 7})
	tp.Register(8, "t2", t2, []any{8, 9, 0})
	tp.Register(61, "t2", t2, []any{8, 9, 0})

	//40秒后关闭
	time.AfterFunc(time.Second*120, func() {
		tp.Close()
	})

	tp.Run()

}
