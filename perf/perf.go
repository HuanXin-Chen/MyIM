package perf

import (
	"fmt"
	"github.com/HuanXin-Chen/MyIM/common/sdk"
	"net"
)

var (
	TcpConnNum int32
)

func RunMain() {
	for i := 0; i < int(TcpConnNum); i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// 处理 panic
					fmt.Printf("第 %d 次循环发生 panic: %v，继续执行下一次循环\n", i+1, r)
				}
			}()
			sdk.NewChat(net.ParseIP("127.0.0.1"), 8900, "chx", "1223", "123")
		}()
		//sdk.NewChat(net.ParseIP("127.0.0.1"), 8900, "chx", "1223", "123")
	}
}
