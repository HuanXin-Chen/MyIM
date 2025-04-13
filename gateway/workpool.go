package gateway

import (
	"fmt"

	"github.com/HuanXin-Chen/MyIM/common/config"
	"github.com/panjf2000/ants"
)

// 协程池初始化

var wPool *ants.Pool

func initWorkPoll() {
	var err error
	if wPool, err = ants.NewPool(config.GetGatewayWorkerPoolNum()); err != nil {
		fmt.Printf("InitWorkPoll.err :%s num:%d\n", err.Error(), config.GetGatewayWorkerPoolNum())
	}
}
